package work

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func compile(w *Work) *rErr {

	os.MkdirAll(w.binaryPath, os.ModePerm)
	outputPath := w.binaryPath + "main"

	src := "source/test.c"

	var compileCmd *exec.Cmd

	switch w.Language {
	case "c":
		compileCmd = exec.Command("gcc", src, "-w", "-o", outputPath)
	}

	compileErr, err := compileCmd.StderrPipe()
	if err != nil {
		log.Panic(err)
	}

	compileCmd.Start()

	errMsg, _ := ioutil.ReadAll(compileErr)
	if len(errMsg) != 0 {
		return Error_R(-1, CompileErr, errMsg)
	}

	compileCmd.Wait()

	return nil

}
