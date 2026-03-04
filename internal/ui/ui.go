package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maclav3/jira-newbranch/internal/jira"
)

func SelectIssue(issues []jira.JiraIssue) int {
	fmt.Println("Select a Jira task:")
	for i, issue := range issues {
		updated, err := jira.ParseJiraTime(issue.Fields.Updated)
		if err != nil {
			log.Fatalf("failed to parse updated time %q: %v", issue.Fields.Updated, err)
		}
		fmt.Printf("[%d] %s: %s (Updated: %s)\n", i+1, issue.Key, issue.Fields.Summary, updated.Format(time.RFC3339))
	}

	fmt.Print("Enter number: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	var choice int
	_, err = fmt.Sscanf(input, "%d", &choice)
	if err != nil || choice < 1 || choice > len(issues) {
		log.Fatal("Invalid selection")
	}

	return choice
}
