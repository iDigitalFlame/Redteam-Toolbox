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

type server struct {
	dir http.Handler
	log io.WriteCloser
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s <dir> <log>\n", os.Args[0])
		os.Exit(1)
	}

	h := &server{dir: http.FileServer(http.Dir(os.Args[1]))}
	h.init(os.Args[2])
	http.Handle("/", h)

	defer h.log.Close()
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}
func (s *server) init(p string) {
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	s.log = f
}
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		b := &bytes.Buffer{}
		io.Copy(b, r.Body)
		fmt.Fprintf(s.log, "%s\n", string(b.Bytes()))
		r.Body.Close()
	} else {
		s.dir.ServeHTTP(w, r)
	}
}
