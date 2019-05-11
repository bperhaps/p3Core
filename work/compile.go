package work

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func compile(w *Work) (bool, error) {

	os.MkdirAll(w.binaryPath, os.ModePerm)
	outputPath := w.binaryPath + "main"

	src := "test.c"

	var compileCmd *exec.Cmd
	var compileErr io.ReadCloser

	switch w.Language {
	case "c":
		compileCmd = exec.Command("gcc", src, "-w", "-o", outputPath)
	}

	compileErr, _ = compileCmd.StderrPipe()

	compileCmd.Start()

	errorByte, _ := ioutil.ReadAll(compileErr)
	if len(errorByte) != 0 {
		return false, errors.New(string(errorByte))
	}

	compileCmd.Wait()

	if _, err := os.Stat(outputPath); err != nil {
		return false, errors.New("binary file not exist.")
	}

	return true, nil

}
