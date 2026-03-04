package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Updated string `json:"updated"`
	} `json:"fields"`
}

type JiraSearchResponse struct {
	Issues []JiraIssue `json:"issues"`
}

func main() {
	// Early check for git repository
	if err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Run(); err != nil {
		log.Fatal("Current directory is not a git repository")
	}

	jiraURL := strings.TrimSuffix(os.Getenv("JIRA_URL"), "/")
	jiraToken := os.Getenv("JIRA_TOKEN")
	jiraUser := os.Getenv("JIRA_USER")

	if jiraURL == "" || jiraToken == "" || jiraUser == "" {
		log.Fatal("JIRA_URL, JIRA_TOKEN, and JIRA_USER must be set")
	}

	// JQL: assignee = 'user' AND statusCategory != 'Done' ORDER BY updated DESC
	jql := fmt.Sprintf("assignee = '%s' AND statusCategory != 'Done' ORDER BY updated DESC", jiraUser)
	searchURL := fmt.Sprintf("%s/rest/api/3/search/jql", jiraURL)

	payload := map[string]any{
		"jql":    jql,
		"fields": []string{"summary", "updated"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("failed to marshal search payload: %v", err)
	}

	req, err := http.NewRequest("POST", searchURL, strings.NewReader(string(body)))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.SetBasicAuth(jiraUser, jiraToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to fetch issues: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to fetch issues: status %s", resp.Status)
	}

	var searchResp JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		log.Fatalf("failed to decode response: %v", err)
	}

	issues := searchResp.Issues
	if len(issues) == 0 {
		fmt.Println("No active issues assigned to you.")
		return
	}

	var choice int
	if len(issues) == 1 {
		choice = 0
	} else {
		choice = selectIssue(issues) - 1
	}

	selected := issues[choice]
	branchName := formatBranchName(selected.Key, selected.Fields.Summary)
	branchName = getAvailableBranchName(branchName)
	fmt.Printf("Creating branch: %s\n", branchName)

	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to create git branch: %v", err)
	}
}

func selectIssue(issues []JiraIssue) int {
	fmt.Println("Select a Jira task:")
	for i, issue := range issues {
		updated, err := parseJiraTime(issue.Fields.Updated)
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

func parseJiraTime(timeStr string) (time.Time, error) {
	// Jira sometimes returns time as "2026-02-16T13:51:20.182+0000"
	// time.RFC3339Nano expects "Z07:00"
	const jiraTimeLayout = "2006-01-02T15:04:05.000-0700"
	updated, err := time.Parse(jiraTimeLayout, timeStr)
	if err != nil {
		// Fallback to RFC3339 if custom layout fails
		updated, err = time.Parse(time.RFC3339, timeStr)
	}
	return updated, err
}

func formatBranchName(key, summary string) string {
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

func getAvailableBranchName(baseName string) string {
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
		if !branchExists(candidate) {
			return candidate
		}
		counter++
		candidate = fmt.Sprintf("%s-%d", baseName, counter)
	}
}

func branchExists(name string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+name).Run()
	return err == nil
}
