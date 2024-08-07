package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// MrLlamaReview performs a code review on the given directory
func MrLlamaReview(dir string) {
	// Change to the specified directory
	if err := os.Chdir(dir); err != nil {
		fmt.Println("Error changing directory:", err)
		os.Exit(1)
	}

	// Check if the directory is a git directory
	if _, err := exec.Command("git", "status").Output(); err != nil {
		fmt.Println("Error: directory is not a git repository")
		os.Exit(1)
	}

	// Get the current branch name
	currentBranch, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		fmt.Println("Error getting current branch:", err)
		os.Exit(1)
	}
	currentBranchName := strings.TrimSpace(string(currentBranch))

	// Print th ebranch name being compared against main
	fmt.Printf("Comparing branch '%s' against 'main'\n", currentBranchName)

	// Get the diff between the current branch and main
	diffOutput, err := exec.Command("git", "diff", "main", strings.TrimSpace(string(currentBranch))).Output()
	if err != nil {
		fmt.Println("Error getting diff:", err)
		os.Exit(1)
	}

	// Write the diff output to a file
	diffFilePath := "diff_output.txt"
	if err := os.WriteFile(diffFilePath, diffOutput, 0644); err != nil {
		fmt.Println("Error writing diff to file:", err)
		os.Exit(1)
	}

	// Read the diff file content
	diffContent, err := os.ReadFile(diffFilePath)
	if err != nil {
		fmt.Println("Error reading diff file:", err)
		os.Exit(1)
	}

	// Define the prompt
	prompt := fmt.Sprintf(`You are a Principal Security Engineer with expertise in code reviews and security issues. Please review the following code changes between the 'main' branch and the '%s' branch:
	%s

	Your review should cover the following aspects:
	1. Security vulnerabilities: Identify any potential security risks or vulnerabilities in the code.
	2. Code quality: Comment on the overall quality of the code, including readability, maintainability, and adherence to best practices.
	3. Functional correctness: Ensure that the changes function as intended and do not introduce any bugs.
	4. Performance implications: Assess whether the changes might negatively impact performance.
	5. Compliance with security policies: Verify that the code complies with the organization's security policies and standards.
	6. Suggestions for improvement: Provide any suggestions for improving the code, both in terms of security and general quality.

	Please provide detailed comments and suggestions for each aspect.`, currentBranchName, string(diffContent))

	// Create the payload for the API request
	payload := map[string]interface{}{
		"model":  "llama3.1",
		"prompt": prompt,
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		os.Exit(1)
	}

	// Make the API request
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error making API request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-200 response: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		os.Exit(1)
	}

	// Print the review output
	fmt.Println("Code Review Output:")
	if response, ok := result["response"].(string); ok {
		fmt.Println(response)
	} else {
		fmt.Println("Unexpected response format")
	}
}
