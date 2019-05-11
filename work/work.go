package work

import (
	"encoding/json"
	"os"
)

const (
	Practice = 1
	Actual   = 2
)

type Work struct {
	UserId        string   `json:"userId"`
	Language      string   `json:"language"`
	ProblemNumber int      `json:"problemNumber"`
	MemLimit      int64    `json:"memLimit"`
	TimeLimit     int      `json:"timeLimit"`
	InputList     []string `json:"inputList"`
	OutputList    []string `json:"outputList"`

	binaryPath string

	execCmd []string
}

func NewWork(jsonData []byte) *Work {
	//test

	os.Setenv("P3_INPUTPATH", "input/")
	os.Setenv("P3_OUTPUTPATH", "output/")
	os.Setenv("P3_BINARYPATH", "binary/")

	//////////////////////////

	binaryPath := os.Getenv("P3_BINARYPATH")

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		os.MkdirAll(binaryPath, os.ModePerm)
	}

	work := &Work{}
	json.Unmarshal(jsonData, &work)

	work.binaryPath = binaryPath

	return work
}
