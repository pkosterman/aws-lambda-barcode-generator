package osutil

import (
	"bytes"
	"os/exec"
	"strings"
)

// Helper function to make simplify command line access
func Run(dir string, commandLine string) ([]byte, error) {

	// Split commandLine into an array separated by whitespace
	args := strings.Fields(commandLine)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := buf.Bytes()
	if err != nil {
		return out, err
	}
	return out, nil
}
