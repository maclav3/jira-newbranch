package branch

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/maclav3/jira-newbranch/internal/ui"
)

func FormatBranchName(key, summary string) string {
	// Skip function words: to, and, or, be, for, with, etc.
	skipWords := map[string]struct{}{
		"to":   {},
		"and":  {},
		"or":   {},
		"of":   {},
		"be":   {},
		"for":  {},
		"with": {},
		"the":  {},
		"a":    {},
		"an":   {},
	}

	// Clean summary: lowercase, remove non-alphanumeric, split into words
	summary = strings.ToLower(summary)
	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	summary = reg.ReplaceAllString(summary, " ")
	words := strings.Fields(summary)

	var filtered []string
	for _, w := range words {
		if _, skip := skipWords[w]; !skip {
			filtered = append(filtered, w)
			if len(filtered) == 4 {
				break
			}
		}
	}

	suffix := strings.Join(filtered, "-")
	if suffix == "" {
		return key
	}
	return fmt.Sprintf("%s-%s", key, suffix)
}

func GetAvailableBranchName(baseName string) string {
	candidate := baseName
	counter := 1

	// Check if baseName already has a numeric suffix (e.g., "branch-name-2").
	// If it does, extract the base name and the counter to correctly increment it.
	re := regexp.MustCompile(`^(.*)-(\d+)$`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) == 3 {
		baseName = matches[1]
		fmt.Sscanf(matches[2], "%d", &counter)
	}

	for {
		if !BranchExists(candidate) {
			return candidate
		}
		counter++
		candidate = fmt.Sprintf("%s-%d", baseName, counter)
	}
}

func BranchExists(name string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name).Run()
	return err == nil
}

func EnsureAvailableBranchName(branchName string) (string, bool) {
	// If the initial branch name is already taken, ask for confirmation.
	if !BranchExists(branchName) {
		return branchName, true
	}

	nextName := GetAvailableBranchName(branchName)
	msg := fmt.Sprintf("Branch '%s' already exists. Create '%s' instead?", branchName, nextName)
	if !ui.Confirm(msg) {
		return "", false
	}
	return nextName, true
}

func IsGitRepo() bool {
	return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil
}

func CheckoutNewBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
