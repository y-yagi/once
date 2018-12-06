package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		usage()
		os.Exit(0)
	}

	cmd := args[1]

	lockfile := "/tmp/once-" + cmd + ".lock"

	if _, err := os.Stat(lockfile); err == nil {
		os.Exit(1)
	}

	ioutil.WriteFile(lockfile, []byte(""), 0644)
	fd, _ := syscall.Open(lockfile, syscall.O_RDONLY, 0000)
	defer syscall.Close(fd)

	if err := syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		os.Exit(1)
	}

	err := runCmd(cmd, args[2:], os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf("Usage: once cmd\n")
}

func runCmd(cmd string, args []string, r io.Reader, out, err io.Writer) error {
	execCmd := exec.Command(cmd, args...)

	execCmd.Stderr = err
	execCmd.Stdout = out
	execCmd.Stdin = r

	return execCmd.Run()
}
