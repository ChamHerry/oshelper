package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("adb", "shell", "uiautomator", "runtest", "your_test_file.jar", "-c", "your.package.name.YourTestClass")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Output: %s\n", output)
}
