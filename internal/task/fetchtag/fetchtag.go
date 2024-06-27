package fetchtag

import (
	"github.com/gembaadvantage/uplift/internal/context"
	git "github.com/purpleclay/gitz"
)

// Task for fetching tags from a remote repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "fetching all tags"
}

// Skip running the task if no flag has been provided
func (t Task) Skip(ctx *context.Context) bool {
	return !ctx.FetchTags
}

// Run the task fetching all tags from the remote repository
func (t Task) Run(ctx *context.Context) error {
	_, err := ctx.GitClient.Fetch(git.WithAll(), git.WithTags())
	return err
}
