package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
)

func TestCommitMessages(t *testing.T) {
	// Mock the output of the Git command
	// 	output := []byte(`123455 feat: tst
	// 123456 fix: tst
	// 123457 docs: tst`)

	cmd := exec.Command("git", "log", "--pretty=format:%h %s", "--abbrev-commit", "--no-merges", "main..")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing Git command:", err)
		return
	}

	// Execute the test
	// Create a scanner to read the output
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// Create a list to store the commits
	commits := []Commit{}

	// Loop through the output
	for scanner.Scan() {
		line := scanner.Text()
		commit := Commit{
			Hash:    line[:7],
			Message: line[8:],
		}
		commits = append(commits, commit)
	}

	for _, commit := range commits {
		t.Run(commit.Hash, func(t *testing.T) {
			_, err := parser.NewMachine(conventionalcommits.WithTypes(conventionalcommits.TypesConventional)).Parse([]byte(commit.Message))
			if err != nil {
				t.Errorf("\nCommit message does not conform to Conventional Commits v1.0 specifications\n\nCommit:\n%s %s\n\nError:\n%s\n\nReference:\nhttps://www.conventionalcommits.org/en/v1.0.0/#specification", commit.Hash, commit.Message, err)
			}
		})
	}

}
