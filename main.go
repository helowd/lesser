package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

func getTerminalSize() (int, int, error) {
	var dimensions struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&dimensions)),
	)
	if errno != 0 {
		return 0, 0, fmt.Errorf("Failed to get terminal size: %v", errno)
	}
	return int(dimensions.rows), int(dimensions.cols), nil
}

func runCommandWithLess(command string) {
	cmd := exec.Command("bash", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating StdoutPipe: %v\n", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating StderrPipe: %v\n", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting command: %v\n", err)
		return
	}

	rows, _, err := getTerminalSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		return
	}

	buf, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdout: %v\n", err)
		return
	}

	if len(buf) > 0 {
		output := string(buf)
		lines := strings.Split(output, "\n")

		if len(lines) > rows {
			lessCmd := exec.Command("less")
			lessCmd.Stdin = strings.NewReader(output)
			lessCmd.Stdout = os.Stdout
			lessCmd.Stderr = os.Stderr
			if err := lessCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running less: %v\n", err)
                return 
			}
            fmt.Print(output)
		} else {
			fmt.Print(output)
		}
	}

	bufErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stderr: %v\n", err)
		return
	}

	if len(bufErr) > 0 {
		fmt.Fprint(os.Stderr, string(bufErr))
	}

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error waiting for command: %v\n", err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: le <command>")
		os.Exit(1)
	}

	commandToRun := strings.Join(os.Args[1:], " ")
	runCommandWithLess(commandToRun)
}
