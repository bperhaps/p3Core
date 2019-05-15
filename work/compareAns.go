package work

import (
	"bytes"
	"crypto/sha256"
)

func (w Work) ansCompare(output rOut) bool {

	caseNum := output.CaseNum()
	userAns := output.Result()

	ans := []byte(w.OutputList[caseNum])

	ansHash := sha256.Sum256(ans)
	userAnsHash := sha256.Sum256(userAns)

	if bytes.Compare(ansHash[:], userAnsHash[:]) == 0 {
		return true
	}

	return false
}
