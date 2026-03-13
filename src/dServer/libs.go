package main

import (
	"os"
	"strconv"
	"errors"
)

//write str with no newline
func wr(str string) {
	os.Stdout.WriteString(str)
}

//write byte to stdout
func wrb(byt []byte) {
	os.Stdout.Write(byt)
}

//write ln
func wrl(str string) {
	wr(str + "\n")
}

//just print err
func werr(err error) {
	os.Stderr.WriteString("" + err.Error() + "\n")
}

//just print err from str
func wserr(err string) {
	werr(errors.New(err))
}

//fatal err 
func ferr(err error) {
	werr(err)
	os.Exit(1)
}

//fatal err from str
func fserr(err string) {
	ferr(errors.New(err))
}

//handle errs as non-fatal
func hanErr(err error) {
	if err != nil {
		werr(err)
	}
}

//handle errs as fatal
func hanFrr(err error) {
	if err != nil {
		ferr(err)
	}
}

//mk err
func merr(str string, err error) {
	erorStr := str + err.Error()
	eror := errors.New(erorStr)
	werr(eror)
}

func idx_of_str(haystack []string, needle string) int {
	for i, thing :=  range haystack {
		if thing == needle {
			wrl("\033[32m" + needle + "\033[0m (" + strconv.Itoa(i) + ")")
			return i
		}
	}
	return -1
}

