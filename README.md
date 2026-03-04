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

### Using Go

The simplest way is to use `go install`:

```bash
go install github.com/maclav3/jira-newbranch@latest
```

### Manual build
Or build it manually:

```bash
go build -o jira-newbranch main.go
# Optionally move it to your PATH
mv jira-newbranch /usr/local/bin/
```

### Installation with Task

Or use the `install` task to build and install the executable to a specific path (defaults to `~/bin/`):

```bash
task install
# Or specify a custom destination
task install DEST=/usr/local/bin/
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

To run the unit tests for branch naming and time parsing logic, you can use `go test`:

```bash
go test -v ./...
```

Or use the provided `Taskfile`:

```bash
task test
```

## Formatting

To format all Go files in the project:

```bash
task fmt
```

## CI/CD and Versioning

This project uses GitHub Actions for continuous integration and automated versioning.

### Automated Tagging
On every push to the `main` branch, the CI pipeline:
1. Runs all tests.
2. If tests pass, it automatically creates a new SemVer-compatible git tag (e.g., `v1.0.1`).

### Bumping Versions
By default, the pipeline performs a **patch** bump. You can trigger a **minor** or **major** bump using the following methods:

1. **Commit Messages**: Include `#major`, `#minor`, or `#patch` in your commit message.
2. **Manual Dispatch**: Go to the **Actions** tab in GitHub, select the **CI/CD** workflow, and use the **Run workflow** button to choose the bump level.
