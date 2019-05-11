package work

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"syscall"
	"time"
)

var AnswerErr = errors.New("inCorrect")
var MemoryErr = errors.New("memoryErr")
var TimeOutErr = errors.New("timeoutErr")
var RuntimeErr = errors.New("runtimeErr")
var UnknownErr = errors.New("unknownErr")

type Output struct {
	caseNum string
	Result  []byte
	time    time.Duration
	Errors  error
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

func (w *Work) execution(mode int) []*Output {
	cmdList := []*exec.Cmd{}
	defer w.killProcess(cmdList)

	outputChan := make(chan *Output)

	var caseLength int
	caseLength = len(w.InputList)
	for i, inputData := range w.InputList {
		go func(i int, inputData string) {
			resultChan := make(chan *Output)
			cmd := w.getCmd()
			cmdList = append(cmdList, cmd)
			go runProcess(cmd, strconv.Itoa(i), []byte(inputData), w, resultChan, mode)

			select {
			case <-time.After(time.Duration(w.TimeLimit) * time.Millisecond):
				outputChan <- &Output{
					strconv.Itoa(i),
					nil,
					time.Duration(0),
					TimeOutErr,
				}
			case output := <-resultChan:
				outputChan <- output
			}
		}(i, inputData)
	}

	var result []*Output

	for i := 0; i < caseLength; i++ {
		output := <-outputChan

		switch output.Errors {
		case nil:
			switch mode {
			case Practice:
				result = append(result, output)
				continue
			case Actual:
				if w.ansCompare(output) {
					result = append(result, output)
					continue
				}

				return []*Output{{
					output.caseNum,
					nil,
					time.Duration(0),
					AnswerErr,
				}}
			}
		default:
			switch mode {
			case Practice:
				result = append(result, output)
				continue
			case Actual:
				return []*Output{output}
			}
		}
	}

	//practice
	return result
}

func (w *Work) Run(mode int) (string, error) {
	defer os.RemoveAll(w.binaryPath)

	_, CompileErr := compile(w)
	//compile Error
	if CompileErr != nil {
		return "", CompileErr
	}

	outputs := w.execution(mode)

	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].caseNum < outputs[j].caseNum
	})

	var result string

	switch mode {
	case Practice:

		for _, output := range outputs {
			if output.Errors != nil {
				result += fmt.Sprintf("case %s : %s\n", output.caseNum, output.Errors.Error())
			} else {
				result += fmt.Sprintf("case %s : %s\n", output.caseNum, string(output.Result))
			}
		}
		return result, nil
	case Actual:
		for _, output := range outputs {
			if output.Errors != nil {
				result += fmt.Sprintf("%s\n", output.Errors.Error())
			} else {
				result += fmt.Sprintf("%s %s\n", string(output.Result), output.time)
			}
		}
	}

	return result, nil

}

func runProcess(cmd *exec.Cmd, caseNum string, inputData []byte, w *Work, resultChan chan *Output, mode int) {

	stdIn, _ := cmd.StdinPipe()
	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	startTime := time.Now()
	cmd.Start()
	go memoryCheck(cmd.Process.Pid, caseNum, w.MemLimit, resultChan)

	stdIn.Write(inputData)
	stdIn.Close()

	errput, err := ioutil.ReadAll(stdErr)
	output, err := ioutil.ReadAll(stdOut)

	if err != nil {
		fmt.Println(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		errput = append(errput, []byte(err.Error())...)
	}
	endTime := time.Since(startTime)

	maxrss := cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss

	if maxrss > w.MemLimit {
		resultChan <- &Output{
			caseNum,
			nil,
			time.Duration(0),
			MemoryErr,
		}
	} else if len(output) != 0 {
		resultChan <- &Output{
			caseNum,
			output,
			endTime,
			nil,
		}
	} else if len(errput) != 0 {
		resultChan <- &Output{
			caseNum,
			errput,
			time.Duration(0),
			RuntimeErr,
		}
	} else {
		resultChan <- &Output{
			caseNum,
			nil,
			time.Duration(0),
			UnknownErr,
		}
	}
}

func memoryCheck(pid int, caseNum string, memLimit int64, memoryChan chan *Output) {
	for {
		mem, err := calculateMemory(pid)
		if err != nil {
			break
		}

		if mem > memLimit {
			memoryChan <- &Output{
				caseNum,
				[]byte("memory error"),
				time.Duration(0),
				MemoryErr,
			}
			break
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
