package system

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// NewConsoleReader returns a pointer to a new bufio.Reader that takes in console input (os.Stdin)
func NewConsoleReader() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

// Query asks prints out the desired question, then formats and returns the provided console input
func Query(reader *bufio.Reader, question string) string {
	fmt.Println(question)
	fmt.Print("-> ")

	ipt, _ := reader.ReadString('\n')
	ipt = strings.Replace(ipt, "\r\n", "", -1) // necessary for Windows computers to format correctly
	return ipt
}

// QueryYN asks the desired question with an appended " (Y/N)", then interprets the input as true or false
func QueryYN(reader *bufio.Reader, question string) bool {
	question += " (Y/N)"

	ipt := Query(reader, question)

	if len(ipt) == 0 {
		return false
	}

	yn := ipt[0]
	if yn == 'Y' || yn == 'y' {
		return true
	}

	return false
}
