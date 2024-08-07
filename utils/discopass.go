// func/discopass.go

package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func DiscoPass(filePath string) {
	// Open the input file
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening the input file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Create a new scanner for the input file
	scanner := bufio.NewScanner(f)

	// Loop through each line of the input file and search for passwords using regular expression
	for scanner.Scan() {
		line := scanner.Text()
		match, err := regexp.MatchString("^[a-zA-Z0-9]{12,}$", line)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if match {
			fmt.Println("Found Password:", line)
		}
	}
}
