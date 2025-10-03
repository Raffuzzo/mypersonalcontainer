package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: ./mypersonalcontainer run <command> [args...]")
	os.Exit(1)
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func child() {
	// Controllo che rootfs ci sia
	if _, err := os.Stat("rootfs/bin/sh"); os.IsNotExist(err) {
		panic("rootfs non valido: manca /bin/sh")
	}

	// Imposta il root filesystem
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	must(os.MkdirAll("rootfs/oldrootfs", 0700))
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	must(os.Chdir("/"))

	// Mi stampo le info 
	fmt.Println("✅ Container started!")
	fmt.Println("PID 1 inside container:", os.Getpid())
	fmt.Println("Hostname:", getHostname())

	// Esegue il comando
	if len(os.Args) < 3 {
		usage()
	}
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// Funzione get per retrievarmi subito l hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
