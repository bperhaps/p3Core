package work

import (
	"errors"
	"time"
)

type rOut struct {
	caseNum int
	result  []byte
	time    time.Duration
	memory  int
}

func (r rOut) Memory() int {
	return r.memory
}

func (r rOut) CaseNum() int {
	return r.caseNum
}

func (r rOut) Time() time.Duration {
	return r.time
}

func (r rOut) Result() []byte {
	return r.result
}

type rErr struct {
	errType string
	caseNum int
	error
}

func (r rErr) CaseNum() int {
	return r.caseNum
}

func Output(caseNum int, result []byte, time time.Duration, mem_r int) *rOut {
	return &rOut{
		caseNum: caseNum,
		result:  result,
		time:    time,
		memory:  mem_r,
	}
}

func Error(caseNum int, errType int) *rErr {

	errTypeStr, err := getErrorTypeByString(errType)

	return &rErr{
		errType: errTypeStr,
		caseNum: caseNum,
		error:   err,
	}
}

func Error_R(caseNum int, errType int, r []byte) *rErr {

	errTypeStr, _ := getErrorTypeByString(errType)

	return &rErr{
		errType: errTypeStr,
		caseNum: caseNum,
		error:   errors.New(string(r)),
	}
}

func getErrorTypeByString(errType int) (string, error) {
	var r string
	var e error
	switch errType {
	case CompileErr:
		r = "CompileErr"
		e = errors.New("Compile error occured.")
	case AnswerErr:
		r = "AnswerErr"
		e = errors.New("Answer is not correct.")
	case MemoryErr:
		r = "MemoryErr"
		e = errors.New("Exceeded memory.")
	case TimeOutErr:
		r = "TimeOutErr"
		e = errors.New("Time out.")
	case RuntimeErr:
		r = "RuntimeErr"
		e = errors.New("Runtime error occured")
	default:
		r = "UnknownErr"
		e = errors.New("Unknown error occured")
	}

	return r, e
}
