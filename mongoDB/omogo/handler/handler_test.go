package handler

import (
	"fmt"
	"testing"
)

type printerror struct {
}

func (p printerror) Error() string {
	return "Test PrintError function succeed!"
}

func Test_PrintError(t *testing.T) {
	var err error
	err = new(printerror)
	message := err.Error()
	if message == "Test PrintError function succeed!" {
		fmt.Printf("Succeed!")
	} else {
		fmt.Printf("Failed!")
	}
}

func Test_SetFlags(t *testing.T) {

}

func Test_PrintWithTable(t *testing.T) {

}

func Test_Clinet(t *testing.T) {

}

func Test_insertDocuments(t *testing.T) {

}

func Test_Run(t *testing.T) {

}
