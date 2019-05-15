package work

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"time"
)

const (
	CompileErr = 1
	AnswerErr  = 2
	MemoryErr  = 3
	TimeOutErr = 4
	RuntimeErr = 5
	UnknownErr = 6
)

type rType struct {
	CaseNum int
	Type    string
	Result  string
	Reason  string
	Time    string
	Memory  int
}

func NewrType(data interface{}) *rType {
	switch data.(type) {
	case *rOut:
		d := data.(*rOut)
		return &rType{
			CaseNum: d.CaseNum(),
			Type:    "output",
			Result:  string(d.result),
			Time:    d.Time().String(),
			Memory:  d.Memory(),
		}
	case *rErr:
		d := data.(*rErr)
		return &rType{
			CaseNum: d.CaseNum(),
			Type:    "error",
			Reason:  d.errType,
			Result:  d.Error(),
		}
	}

	return nil
}

func (w Work) getCmd() *exec.Cmd {
	var cmd *exec.Cmd

	switch w.Language {
	case "c":
		cmd = exec.Command(w.binaryPath + "main")
	}

	return cmd
}

func (w *Work) killProcess(cmdList []*exec.Cmd) {
	for _, cmd := range cmdList {
		cmd.Process.Kill()
	}
}

func (w *Work) runProcess(cmd *exec.Cmd, caseNum int, inputData []byte) {

	stdIn, _ := cmd.StdinPipe()
	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	startTime := time.Now()
	cmd.Start()

	var mem_r int
	go runtimeMemoryCheck(cmd.Process.Pid, w.MemLimit, &mem_r, func() {
		//Memory Err occured
		w.execResult <- Error(caseNum, MemoryErr)
	})

	//input data to input stream
	stdIn.Write(inputData)
	stdIn.Close()

	//reading stdErr and stdOut
	errput, err := ioutil.ReadAll(stdErr)
	if err != nil {
		fmt.Println(err.Error())
	}
	output, err := ioutil.ReadAll(stdOut)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		errput = []byte(err.Error())
	}
	endTime := time.Since(startTime)

	//Has the process make a error?
	if len(errput) != 0 {
		w.execResult <- Error(caseNum, RuntimeErr)
		return
	}

	//has the process success to run?
	if len(output) != 0 {
		w.execResult <- Output(caseNum, output, endTime, mem_r)
		return
	}

	//Unknown error
	w.execResult <- Error(caseNum, UnknownErr)
}

func (w *Work) execution(caseNum int, inputData string) {
	cmd := w.getCmd()
	w.cmdList = append(w.cmdList, cmd)

	go w.runProcess(cmd, caseNum, []byte(inputData))

	//Timeout
	<-time.After(time.Duration(w.TimeLimit) * time.Millisecond)
	w.execResult <- Error(caseNum, TimeOutErr)
}

func (w *Work) marking() (interface{}, *rErr) {
	var resultArr []interface{}

	for i := 0; i < len(w.InputList); i++ {
		result := <-w.execResult

		switch result.(type) {
		case *rOut:
			if w.Mode == Practice {
				resultArr = append(resultArr, result)
				continue
			}

			if w.Mode == Actual {
				if w.ansCompare(result.(rOut)) {
					resultArr = append(resultArr, result)
					continue
				}
				//Answer is not correct
				return nil, Error(result.(rOut).CaseNum(), AnswerErr)
			}

		case *rErr:
			if w.Mode == Practice {
				resultArr = append(resultArr, result)
				continue
			}

			if w.Mode == Actual {
				err := result.(rErr)
				return nil, &err
			}
		}
	}

	return resultArr, nil

}

func (w *Work) getResult() (interface{}, *rErr) {
	cmdList := []*exec.Cmd{}
	defer w.killProcess(cmdList)

	for caseNum, caseData := range w.InputList {
		go w.execution(caseNum, caseData)
	}

	return w.marking()
}

func (w *Work) Run() string {
	defer os.RemoveAll(w.binaryPath)

	var resultArr []*rType

	err := compile(w)
	//compile Error
	if err != nil {
		resultArr = append(resultArr, NewrType(err))
		jsonData, _ := json.MarshalIndent(resultArr, "", "    ")
		return string(jsonData)
	}

	out, err := w.getResult()

	for _, r := range out.([]interface{}) {
		resultArr = append(resultArr, NewrType(r))
	}

	sort.Slice(resultArr, func(i, j int) bool {
		return resultArr[i].CaseNum < resultArr[j].CaseNum
	})

	r, _ := json.MarshalIndent(resultArr, "", "    ")
	return string(r)
}

//check memory usage in runtime,
func runtimeMemoryCheck(pid int, memLimit int64, mem_r *int, memErr func()) {
	for {
		mem, err := calculateMemory(pid)
		if err != nil {
			break
		}

		*mem_r = int(mem)

		//Does the process use more memory than the limit?
		if mem > memLimit {
			memErr()
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

//referenced by https://stackoverflow.com/questions/31879817/golang-os-exec-realtime-memory-usage
//this function calculate PSS and return PSS and err
func calculateMemory(pid int) (int64, error) {

	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := int64(0)
	pfx := []byte("Pss:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			var size int64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			res += size
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}

	return res, nil
}
