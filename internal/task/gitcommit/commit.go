package gitcommit

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	git "github.com/purpleclay/gitz"
)

// Task for committing all staged changes and optionally pushing
// them to a git remote
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "committing changes"
}

// Skip running the task
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.DryRun || ctx.NoVersionChanged
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	stg, err := ctx.GitClient.Staged()
	if err != nil {
		return err
	}

	if len(stg) == 0 {
		log.Info("no files staged skipping push")
		return nil
	}

	log.Debug("attempting to commit changes")
	if _, err := ctx.GitClient.Commit(ctx.CommitDetails.Message,
		git.WithCommitConfig("user.name", ctx.CommitDetails.Author.Name,
			"user.email", ctx.CommitDetails.Author.Email)); err != nil {
		return err
	}
	log.Info("staged changes committed")

	if ctx.NoPush {
		log.Warn("skipping push of commit to remote")
		return nil
	}

	log.Info("pushing commit to remote")
	var pushOpts []string
	if ctx.Config.Git != nil {
		pushOpts = filterPushOptions(ctx.Config.Git.PushOptions)
	}
	_, err = ctx.GitClient.Push(git.WithPushOptions(pushOpts...))
	return err
}

func filterPushOptions(options []config.GitPushOption) []string {
	filtered := []string{}
	for _, opt := range options {
		if !opt.SkipBranch {
			log.WithField("option", opt.Option).Debug("with push option")
			filtered = append(filtered, opt.Option)
		}
	}

	return filtered
}
