# Jira NewBranch CLI

A simple CLI application that creates a git branch with a name based on the Jira task you are currently working on.

This project was **vibe-coded** with the help of **Junie**, an autonomous programmer by JetBrains.

## Features

- **Jira Integration**: Fetches tasks assigned to you that are not "Done".
- **Interactive Selection**: Choose from your most recently updated tasks.
- **Smart Branch Naming**: Converts task ID and summary into a clean git branch name.
  - Limits description to 4 keywords.
  - Skips common function words (to, and, or, be, for, with, the, a, an).
  - Example: `PRJ-2134 Adjust foo and bar to be ready for prod` → `PRJ-2134-adjust-foo-bar-ready`.
- **Git Safety**: Fails early if you are not inside a git repository.

## Prerequisites

- **Go 1.25** or later.
- **Git** installed and initialized in your project.

## Installation

```bash
go build -o jira-newbranch main.go
# Optionally move it to your PATH
mv jira-newbranch /usr/local/bin/
```

## Configuration

Set the following environment variables:

- `JIRA_URL`: The base URL of your Jira instance (e.g., `https://your-domain.atlassian.net`).
- `JIRA_TOKEN`: Your Jira API token (Cloud) or Personal Access Token (Data Center).
- `JIRA_USER`: Your Jira email (Cloud) or username (Data Center).

## Usage

Run the tool from within any git repository:

```bash
jira-newbranch
```

1. The tool will verify you are in a git repository.
2. It will fetch your active Jira tasks.
3. Select a task by entering its number.
4. A new git branch will be created and checked out automatically.

## Testing

To run the unit tests for branch naming and time parsing logic:

```bash
go test -v ./...
```
