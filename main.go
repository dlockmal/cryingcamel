//go:generate go run main.go -version v0.0.1
package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/cryingcamel/utils"
)

func main() {
	// Define the command-line-flags
	discopassCmd := flag.NewFlagSet("discopass", flag.ExitOnError)
	discopassFile := discopassCmd.String("file", "", "Input file to search for passwords")

	mrllamareviewCmd := flag.NewFlagSet("mrllamareview", flag.ExitOnError)
	mrllamareviewDir := mrllamareviewCmd.String("dir", "", "Git-initialized directory to review")

	// Check the first argument provided (the command)
	if len(os.Args) < 2 {
		fmt.Println("Expected 'discopass' or 'mrllamareview' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "discopass":
		discopassCmd.Parse(os.Args[2:])
		if *discopassFile == "" {
			fmt.Println("Usage: discopass --file <input_file>")
			os.Exit(1)
		}
		utils.DiscoPass(*discopassFile)
	case "mrllamareview":
		mrllamareviewCmd.Parse(os.Args[2:])
		if *mrllamareviewDir == "" {
			fmt.Println("Usage: mrllamareview --dir <directory>")
			os.Exit(1)
		}
		utils.MrLlamaReview(*mrllamareviewDir)
	default:
		fmt.Println("Expected 'discopass' or 'mrllamareview' subcommands")
		os.Exit(1)
	}
}
