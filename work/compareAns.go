package work

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

func (w Work) ansCompare(output *Output) bool {

	caseNum, _ := strconv.Atoi(output.caseNum)
	userAns := output.Result

	ans := []byte(w.OutputList[caseNum])

	ansHash := sha256.Sum256(ans)
	userAnsHash := sha256.Sum256(userAns)

	if bytes.Compare(ansHash[:], userAnsHash[:]) == 0 {
		return true
	}

	return false
}
