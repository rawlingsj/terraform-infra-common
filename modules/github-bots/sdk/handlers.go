package sdk

import (
	"context"

	"github.com/google/go-github/v61/github"
)

type EventHandlerFunc interface {
	EventType() EventType
}

type PullRequestHandler func(ctx context.Context, pre github.PullRequestEvent, pr *github.PullRequest) error

func (r PullRequestHandler) EventType() EventType {
	return PullRequestEvent
}

type WorkflowRunHandler func(ctx context.Context, wre github.WorkflowRunEvent, wr *github.WorkflowRun) error

func (r WorkflowRunHandler) EventType() EventType {
	return WorkflowRunEvent
}

type CheckRunHandler func(ctx context.Context, wre github.CheckRunEvent, wr *github.CheckRun) error

func (r CheckRunHandler) EventType() EventType {
	return CheckRunEvent
}

type EventType string

const (
	PullRequestEvent EventType = "dev.chainguard.github.pull_request"
	WorkflowRunEvent EventType = "dev.chainguard.github.workflow_run"
	CheckRunEvent    EventType = "dev.chainguard.github.check_run"
)
