package main

// This is a git pre-commit hook that checks any stages file has an HPE
// copyright that needs updating.

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var regex = "Copyright ([,0-9 -]+) Hewlett Packard"

func fails(s string, year string) bool {
	re := regexp.MustCompile(regex)

	if !re.MatchString(s) {
		return false
	}
	matches := re.FindStringSubmatch(s)
	if strings.HasSuffix(matches[1], year[len(year)-2:]) {
		return false
	}
	return true
}

func main() {
	// Get the list of staged files.
	cmd := exec.Command("git", "diff", "--cached", "--name-only")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting staged files: %v\n", err)
		os.Exit(1)
	}

	stagedFiles := strings.Split(out.String(), "\n")

	nFails := 0
	for _, file := range stagedFiles {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		cmd := exec.Command("git", "show", ":"+file)
		var fileContent bytes.Buffer
		cmd.Stdout = &fileContent
		err = cmd.Run()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting file content for %s: %v\n", file, err)
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(fileContent.String()))

		currentYear := fmt.Sprintf("%d", time.Now().Year())
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if fails(line, currentYear) {
				fmt.Fprintf(os.Stderr, "Bad copyright message in %s:%d:'%s'\n", file, lineNumber, line)
				nFails++
			}
			lineNumber++
			if lineNumber > 10 {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
		}
	}

	if nFails > 0 {
		fmt.Fprintf(os.Stderr, "Found %d bad copyright messages. Aborting commit.\n", nFails)
		os.Exit(1)
	}
}
