package main

import "os/exec"

const (
	shel = "powershell.exe"
	comd = []string{}
)

func main() {
	exec.Command(shel, comd...).Wait()
}
