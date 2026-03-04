package main

import (
	"fmt"
	"log"
	"os"

	"github.com/maclav3/jira-newbranch/internal/branch"
	"github.com/maclav3/jira-newbranch/internal/jira"
	"github.com/maclav3/jira-newbranch/internal/ui"
)

func main() {
	if !branch.IsGitRepo() {
		log.Fatal("Current directory is not a git repository")
	}

	jiraURL := os.Getenv("JIRA_URL")
	jiraToken := os.Getenv("JIRA_TOKEN")
	jiraUser := os.Getenv("JIRA_USER")

	if jiraURL == "" || jiraToken == "" || jiraUser == "" {
		log.Fatal("JIRA_URL, JIRA_TOKEN, and JIRA_USER must be set")
	}

	client := jira.NewClient(jiraURL, jiraUser, jiraToken)
	issues, err := client.GetActiveIssues()
	if err != nil {
		log.Fatalf("failed to fetch issues: %v", err)
	}

	if len(issues) == 0 {
		fmt.Println("No active issues assigned to you.")
		return
	}

	var choice int
	if len(issues) == 1 {
		choice = 0
	} else {
		choice = ui.SelectIssue(issues) - 1
	}

	selected := issues[choice]
	branchName := branch.FormatBranchName(selected.Key, selected.Fields.Summary)
	branchName, accepted := branch.EnsureAvailableBranchName(branchName)
	if !accepted {
		fmt.Println("Cancelled.")
		os.Exit(0)
	}

	fmt.Printf("Creating branch: %s\n", branchName)

	if err := branch.CheckoutNewBranch(branchName); err != nil {
		log.Fatalf("failed to create git branch: %v", err)
	}
}
