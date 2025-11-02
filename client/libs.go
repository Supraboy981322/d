/*These fns are for me to save time writing
 *  this, as I find writing `fmt.Println(...)`
 *  (or the `log`/`slog` equivalent) and 
 *  `if err != nil {...}` extremely time
 *  consuming and I'd rather spend that time
 *  writing functionality or doing homework */

package main

import (
	"os"
	"errors"
)


func wr(str string) {
	os.Stdout.WriteString(str)
}


func wrb(byt []byte) {
	os.Stdout.Write(byt)
}


func wrl(str string) {
	wr(str + "\n")
}


func werr(err error) {
	os.Stderr.WriteString("" + err.Error() + "\n")
}


func wserr(err string) {
	werr(errors.New(err))
}


func ferr(err error) {
	werr(err)
	os.Exit(1)
}


func fserr(err string) {
	ferr(errors.New(err))
}


func hanErr(err error) {
	if err != nil {
		werr(err)
	}
}


func hanFrr(err error) {
	if err != nil {
		ferr(err)
	}
}
