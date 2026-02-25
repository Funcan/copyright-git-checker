package main

// This is a git pre-commit hook that checks any stages file has an HPE
// copyright that needs updating.

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	fixFlag   bool
	quietFlag bool
)

var regex = "Copyright ([,0-9 -]+) Hewlett Packard"
var re = regexp.MustCompile(regex)

func fails(s string, year string) bool {
	if !re.MatchString(s) {
		return false
	}
	matches := re.FindStringSubmatch(s)
	if strings.HasSuffix(matches[1], year[len(year)-2:]) {
		return false
	}
	return true
}

// fixLine updates a copyright line to include the current year
func fixLine(s string, year string) string {
	if !re.MatchString(s) {
		return s
	}

	matches := re.FindStringSubmatch(s)
	yearsPart := matches[1]

	// If already ends with current year, no change needed
	if strings.HasSuffix(yearsPart, year[len(year)-2:]) {
		return s
	}

	// Build the new years string
	var newYears string

	// Check if it's a range (contains -)
	if strings.Contains(yearsPart, "-") {
		// Extract the start year from the range
		parts := strings.Split(yearsPart, "-")
		startYear := strings.TrimSpace(parts[0])
		// Handle comma-separated years before the range
		if strings.Contains(startYear, ",") {
			lastComma := strings.LastIndex(startYear, ",")
			startYear = strings.TrimSpace(startYear[lastComma+1:])
			prefix := strings.TrimSpace(yearsPart[:strings.LastIndex(yearsPart, "-")])
			if idx := strings.LastIndex(prefix, ","); idx != -1 {
				prefix = strings.TrimSpace(prefix[:idx+1])
				newYears = prefix + " " + startYear + "-" + year
			} else {
				newYears = startYear + "-" + year
			}
		} else {
			newYears = startYear + "-" + year
		}
	} else if strings.Contains(yearsPart, ",") {
		// It's a comma-separated list, extend the last year to a range
		newYears = yearsPart + "-" + year
	} else {
		// Single year, convert to range
		newYears = strings.TrimSpace(yearsPart) + "-" + year
	}

	// Replace the old years with new years in the original string
	return strings.Replace(s, "Copyright "+yearsPart+" Hewlett Packard", "Copyright "+newYears+" Hewlett Packard", 1)
}

func main() {
	flag.BoolVar(&fixFlag, "fix", false, "Automatically fix copyright headers to include the current year")
	flag.BoolVar(&quietFlag, "quiet", false, "Suppress informational messages (errors are still shown)")
	flag.Parse()

	// Get the list of staged files (excluding deleted files).
	cmd := exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=d")

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

		// Get file info to preserve permissions later
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting file info for %s: %v\n", file, err)
			continue
		}
		fileMode := fileInfo.Mode()

		// Read from the staged content for checking
		cmd := exec.Command("git", "show", ":"+file)
		var stagedContent bytes.Buffer
		cmd.Stdout = &stagedContent
		err = cmd.Run()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting staged content for %s: %v\n", file, err)
			continue
		}

		scanner := bufio.NewScanner(strings.NewReader(stagedContent.String()))

		currentYear := fmt.Sprintf("%d", time.Now().Year())
		lineNumber := 1
		needsFixing := false
		var lines []string

		// First pass: check all lines and collect them
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)
			if lineNumber <= 10 && fails(line, currentYear) {
				if fixFlag {
					needsFixing = true
					if !quietFlag {
						fmt.Fprintf(os.Stderr, "Fixing copyright message in %s:%d:'%s'\n", file, lineNumber, line)
					}
				} else {
					fmt.Fprintf(os.Stderr, "Bad copyright message in %s:%d:'%s'\n", file, lineNumber, line)
					nFails++
				}
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
			continue
		}

		// If --fix is enabled and file needs fixing, update it
		if fixFlag && needsFixing {
			// Read the working directory version to apply fixes to it
			// This preserves any unstaged changes in the file
			workingContent, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading working file %s: %v\n", file, err)
				nFails++
				continue
			}

			workingLines := strings.Split(string(workingContent), "\n")
			// Handle the case where Split creates an extra empty element for trailing newline
			hasTrailingNewline := strings.HasSuffix(string(workingContent), "\n")

			// Fix the lines in the first 10 lines of the working file
			for i := 0; i < len(workingLines) && i < 10; i++ {
				if fails(workingLines[i], currentYear) {
					workingLines[i] = fixLine(workingLines[i], currentYear)
				}
			}

			// Write the fixed content back to the file
			fixedContent := strings.Join(workingLines, "\n")
			// strings.Split on "a\n" gives ["a", ""], so Join restores it correctly
			// But if the file had a trailing newline, we need to ensure it's preserved
			if hasTrailingNewline && !strings.HasSuffix(fixedContent, "\n") {
				fixedContent += "\n"
			}

			err = os.WriteFile(file, []byte(fixedContent), fileMode)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing fixed file %s: %v\n", file, err)
				nFails++
				continue
			}

			// Re-stage the file so changes are included in the commit
			cmd := exec.Command("git", "add", file)
			err = cmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error re-staging file %s: %v\n", file, err)
				nFails++
			}
		}
	}

	if nFails > 0 {
		fmt.Fprintf(os.Stderr, "Found %d bad copyright messages. Aborting commit.\n", nFails)
		os.Exit(1)
	}
}
