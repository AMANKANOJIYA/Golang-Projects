package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "child":
		child()
	default:
		parent()
	}
}

func parent() {

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {

	fmt.Println("running: ", os.Args[2:])

	exit, err := chroot("/home/n/Documents/talks/docker_without_docker/rootfs")
	defer exit()
	must(err)
	must(os.Chdir("/"))

	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("proc", 0)

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Env = append(cmd.Env, "PATH=/bin:/usr/bin")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

/// ----- Boilerplate from here ----- ///

func chroot(path string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(path); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

func must(err error) {
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return
		default:
			panic(err)
		}
	}
}
