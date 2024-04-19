package schemas

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Wrapper[T any] struct {
	When time.Time
	Body T
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#User
type User struct {
	Login bigquery.NullString `json:"login,omitempty" bigquery:"login"`
	Type  bigquery.NullString `json:"type,omitempty" bigquery:"type"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#Organization
type Organization struct {
	Login bigquery.NullString `json:"login,omitempty" bigquery:"login"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#Repository
type Repository struct {
	Owner    User
	Name     bigquery.NullString `json:"name,omitempty" bigquery:"name"`
	URL      bigquery.NullString `json:"url,omitempty" bigquery:"url"`
	FullName bigquery.NullString `json:"full_name,omitempty" bigquery:"full_name"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#PullRequestBranch
type PullRequestBranch struct {
	Ref  bigquery.NullString `json:"ref,omitempty" bigquery:"ref"`
	SHA  bigquery.NullString `json:"sha,omitempty" bigquery:"sha"`
	Repo Repository          `json:"repo,omitempty" bigquery:"repo"`
	User User                `json:"user,omitempty" bigquery:"user"`
}

type Label struct {
	Name bigquery.NullString `json:"name,omitempty" bigquery:"name"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#PullRequest
type PullRequest struct {
	Number bigquery.NullInt64  `json:"number,omitempty" bigquery:"number"`
	State  bigquery.NullString `json:"state,omitempty" bigquery:"state"`
	Title  bigquery.NullString `json:"title,omitempty" bigquery:"title"`

	Base PullRequestBranch `json:"base,omitempty" bigquery:"base"`
	Head PullRequestBranch `json:"head,omitempty" bigquery:"head"`

	Labels []Label `json:"labels" bigquery:"labels"`

	CreatedAt bigquery.NullTimestamp `json:"created_at,omitempty" bigquery:"created_at"`
	UpdatedAt bigquery.NullTimestamp `json:"updated_at,omitempty" bigquery:"updated_at"`
	ClosedAt  bigquery.NullTimestamp `json:"closed_at,omitempty" bigquery:"closed_at"`
	MergedAt  bigquery.NullTimestamp `json:"merged_at,omitempty" bigquery:"merged_at"`

	Mergeable      bigquery.NullBool   `json:"mergeable,omitempty" bigquery:"mergeable"`
	MergeableState bigquery.NullString `json:"mergeable_state,omitempty" bigquery:"mergeable_state"`
	MergedBy       User                `json:"merged_by,omitempty" bigquery:"merged_by"`
	MergeCommitSHA bigquery.NullString `json:"merge_commit_sha,omitempty" bigquery:"merge_commit_sha"`

	Additions    bigquery.NullInt64 `json:"additions,omitempty" bigquery:"additions"`
	Deletions    bigquery.NullInt64 `json:"deletions,omitempty" bigquery:"deletions"`
	ChangedFiles bigquery.NullInt64 `json:"changed_files,omitempty" bigquery:"changed_files"`
}

// https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request
// https://pkg.go.dev/github.com/google/go-github/v60/github#PullRequestEvent
type PullRequestEvent struct {
	// assigned,opened  etc.
	Action     bigquery.NullString `json:"action,omitempty" bigquery:"action"`
	Sender     User                `json:"sender,omitempty" bigquery:"sender"`
	Assignee   User                `json:"assignee,omitempty" bigquery:"assignee"`
	Repository Repository          `json:"repository,omitempty" bigquery:"repository"`

	PullRequest PullRequest `json:"pull_request,omitempty" bigquery:"pull_request"`

	// Populated when action is synchronize
	Before bigquery.NullString `json:"before,omitempty" bigquery:"before"`
	After  bigquery.NullString `json:"after,omitempty" bigquery:"after"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#Workflow
type Workflow struct {
	ID    bigquery.NullInt64  `json:"id,omitempty" bigquery:"id"`
	Name  bigquery.NullString `json:"name,omitempty" bigquery:"name"`
	Path  bigquery.NullString `json:"path,omitempty" bigquery:"path"`
	State bigquery.NullString `json:"state,omitempty" bigquery:"state"`

	CreatedAt bigquery.NullTimestamp `json:"created_at,omitempty" bigquery:"created_at"`
	UpdatedAt bigquery.NullTimestamp `json:"updated_at,omitempty" bigquery:"updated_at"`
}

// https://pkg.go.dev/github.com/google/go-github/v60/github#WorkflowRun
type WorkflowRun struct {
	ID           bigquery.NullInt64     `json:"id,omitempty" bigquery:"id"`
	Name         bigquery.NullString    `json:"name,omitempty" bigquery:"name"`
	NodeID       bigquery.NullString    `json:"node_id,omitempty" bigquery:"node_id"`
	HeadBranch   bigquery.NullString    `json:"head_branch,omitempty" bigquery:"head_branch"`
	HeadSHA      bigquery.NullString    `json:"head_sha,omitempty" bigquery:"head_sha"`
	RunNumber    bigquery.NullInt64     `json:"run_number,omitempty" bigquery:"run_number"`
	RunAttempt   bigquery.NullInt64     `json:"run_attempt,omitempty" bigquery:"run_attempt"`
	Event        bigquery.NullString    `json:"event,omitempty" bigquery:"event"`
	DisplayTitle bigquery.NullString    `json:"display_title,omitempty" bigquery:"display_title"`
	Status       bigquery.NullString    `json:"status,omitempty" bigquery:"status"`
	Conclusion   bigquery.NullString    `json:"conclusion,omitempty" bigquery:"conclusion"`
	WorkflowID   bigquery.NullInt64     `json:"workflow_id,omitempty" bigquery:"workflow_id"`
	CheckSuiteID bigquery.NullInt64     `json:"check_suite_id,omitempty" bigquery:"check_suite_id"`
	URL          bigquery.NullString    `json:"url,omitempty" bigquery:"url"`
	HTMLURL      bigquery.NullString    `json:"html_url,omitempty" bigquery:"html_url"`
	CreatedAt    bigquery.NullTimestamp `json:"created_at,omitempty" bigquery:"created_at"`
	UpdatedAt    bigquery.NullTimestamp `json:"updated_at,omitempty" bigquery:"updated_at"`
	RunStartedAt bigquery.NullTimestamp `json:"run_started_at,omitempty" bigquery:"run_started_at"`
	JobsURL      bigquery.NullString    `json:"jobs_url,omitempty" bigquery:"jobs_url"`
	LogsURL      bigquery.NullString    `json:"logs_url,omitempty" bigquery:"logs_url"`
	ArtifactsURL bigquery.NullString    `json:"artifacts_url,omitempty" bigquery:"artifacts_url"`
	CancelURL    bigquery.NullString    `json:"cancel_url,omitempty" bigquery:"cancel_url"`
	RerunURL     bigquery.NullString    `json:"rerun_url,omitempty" bigquery:"rerun_url"`
	WorkflowURL  bigquery.NullString    `json:"workflow_url,omitempty" bigquery:"workflow_url"`
}

//type WorkflowRun struct {
//	ID           bigquery.NullInt64     `json:"id,omitempty" bigquery:"id"`
//	RunNumber    bigquery.NullInt64     `json:"run_number,omitempty" bigquery:"run_number"`
//	RunAttempt   bigquery.NullInt64     `json:"run_attempt,omitempty" bigquery:"run_attempt"`
//	HeadBranch   bigquery.NullString    `json:"head_branch,omitempty" bigquery:"head_branch"`
//	HeadSHA      bigquery.NullString    `json:"head_sha,omitempty" bigquery:"head_sha"`
//	Name         bigquery.NullString    `json:"name,omitempty" bigquery:"name"`
//	Event        bigquery.NullString    `json:"event,omitempty" bigquery:"event"`
//	Status       bigquery.NullString    `json:"status,omitempty" bigquery:"status"`
//	RunStartedAt bigquery.NullTimestamp `json:"run_started_at,omitempty" bigquery:"run_started_at"`
//
//	// success, failure, cancelled, etc.
//	Conclusion bigquery.NullString `json:"conclusion,omitempty" bigquery:"conclusion"`
//}

// https://docs.github.com/developers/webhooks-and-events/webhook-events-and-payloads#workflow_run
// subset of https://pkg.go.dev/github.com/google/go-github/v60/github#WorkflowRunEvent
type WorkflowRunEvent struct {
	// completed, etc.
	Action       bigquery.NullString `json:"action,omitempty" bigquery:"action"`
	Workflow     Workflow            `json:"workflow,omitempty" bigquery:"workflow"`
	WorkflowRun  WorkflowRun         `json:"workflow_run,omitempty" bigquery:"workflow_run"`
	Organization Organization        `json:"organization,omitempty" bigquery:"organization"`
	Repository   Repository          `json:"repository,omitempty" bigquery:"repository"`
	Sender       User                `json:"sender,omitempty" bigquery:"sender"`
}
