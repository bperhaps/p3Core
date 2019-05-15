package main

import (
	"fmt"
	"p3Core/work"
)

//var workingPool map[string]*work.Work

func main() {

	w := work.NewWork([]byte(`{
		"userId" : "sms2831",
		"language" : "c",
		"problemNumber" :1,
		"memLimit" : 3000,
		"timeLimit" : 2000,
		"inputList" : ["1", "2", "3"],
		"outputList" : ["7", "7", "7"],
		"mode" : 1
	}`))
	r := w.Run()

	fmt.Println(r)
}
