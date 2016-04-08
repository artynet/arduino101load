package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/kardianos/osext"
	"runtime"
	"strings"
)

var verbose bool

func PrintlnVerbose(a ...interface{}) {
	if verbose {
		fmt.Println(a...)
	}
}

func main() {
	fmt.Println("Starting download script...")

	// ARG 1: ELF File to download
	// ARG 2: serialnumber of target device
	// path may contain \ need to change all to /

	args := os.Args[1:]

	bin_path, err := osext.ExecutableFolder()
	adb := bin_path + "/adb"
	adb = filepath.ToSlash(adb)
	if (len(args) != 3) {
		fmt.Println("Wrong parameter list")
		os.Exit(1)
	}

	bin_file_name := args[0]
	verbosity := args[1]
	serialnumber := args[2]

	if verbosity == "quiet" {
		verbose = false
	} else {
		verbose = true
	}

	PrintlnVerbose("Args to shell:", args)
	PrintlnVerbose("Serial Number: " + serialnumber)
	PrintlnVerbose("BIN FILE " + bin_file_name)

	if runtime.GOOS == "darwin" {
		library_path := os.Getenv("DYLD_LIBRARY_PATH")
		if !strings.Contains(library_path, bin_path) {
			os.Setenv("DYLD_LIBRARY_PATH", bin_path+":"+library_path)
		}
	}

	adb_search_command := []string{adb, "devices"}

	err, found := launchCommandAndWaitForOutput(adb_search_command, serialnumber, false)

	if (err == nil && found == false) {
		err, found = launchCommandAndWaitForOutput(adb_search_command, strings.ToUpper(serialnumber), false)
		if (found == true) {
			serialnumber = strings.ToUpper(serialnumber)
		}
	}

	if (err == nil && found == false) {
		err, found = launchCommandAndWaitForOutput(adb_search_command, strings.ToLower(serialnumber), false)
		if (found == true) {
			serialnumber = strings.ToUpper(serialnumber)
		}
	}

	if (err != nil) {
		fmt.Println("ERROR: Upload failed")
		os.Exit(1)
	}

	var serialnumberslice []string

	if (found == true) {
		serialnumberslice = []string{"-s",  serialnumber}
	}

	adb_push := []string{adb}
	adb_push = append(adb_push, serialnumberslice...)
	adb_push = append(adb_push, "push", bin_file_name, "/tmp/sketch.bin")
	err, _ = launchCommandAndWaitForOutput(adb_push, "", true)

	if err == nil {
		fmt.Println("SUCCESS!")
		os.Exit(0)
	} else {
		fmt.Println("ERROR: Upload failed")
		os.Exit(1)
	}
}

func launchCommandAndWaitForOutput(command []string, stringToSearch string, print_output bool) (error, bool) {
	oscmd := exec.Command(command[0], command[1:]...)
	tellCommandNotToSpawnShell(oscmd)
	stdout, _ := oscmd.StdoutPipe()
	stderr, _ := oscmd.StderrPipe()
	multi := io.MultiReader(stderr, stdout)
	err := oscmd.Start()
	in := bufio.NewScanner(multi)
	in.Split(bufio.ScanLines)
	found := false
	for in.Scan() {
		if print_output {
			PrintlnVerbose(in.Text())
		}
		if stringToSearch != "" {
			if strings.Contains(in.Text(), stringToSearch) {
				found = true
			}
		}
	}
	err = oscmd.Wait()
	return err, found
}
