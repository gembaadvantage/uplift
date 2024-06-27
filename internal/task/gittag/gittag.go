package gittag

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	git "github.com/purpleclay/gitz"
)

// Task for tagging a repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "tagging repository"
}

// Skip running the task if no version has changed
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged
}

// Run the task tagging a repository with the next semantic version. Supports both
// standard and annotated git tags
func (t Task) Run(ctx *context.Context) error {
	if ctx.CurrentVersion.Raw == ctx.NextVersion.Raw {
		log.WithFields(log.Fields{
			"current": ctx.CurrentVersion.Raw,
			"next":    ctx.NextVersion.Raw,
		}).Warn("no version change detected")
		return nil
	}

	log.WithField("tag", ctx.NextVersion.Raw).Info("identified next tag")
	if ctx.DryRun {
		log.Info("skipping tag in dry run mode")
		return nil
	}

	if ctx.PrintCurrentTag || ctx.PrintNextTag {
		printRepositoryTag(ctx)
		return nil
	}

	log.Debug("attempting to tag repository")
	if ctx.Config.AnnotatedTags {
		if _, err := ctx.GitClient.Tag(ctx.NextVersion.Raw,
			git.WithTagConfig("user.name", ctx.CommitDetails.Author.Name, "user.email", ctx.CommitDetails.Author.Email),
			git.WithAnnotation(ctx.CommitDetails.Message)); err != nil {
			return err
		}
		log.Info("tagged repository with annotated tag")
	} else {
		if _, err := ctx.GitClient.Tag(ctx.NextVersion.Raw); err != nil {
			return err
		}
		log.Info("tagged repository with lightweight tag")
	}

	if ctx.NoPush {
		log.Warn("skipping push of tag to remote")
		return nil
	}

	log.Info("pushing tag to remote")
	var pushOpts []string
	if ctx.Config.Git != nil {
		pushOpts = filterPushOptions(ctx.Config.Git.PushOptions)
	}

	_, err := ctx.GitClient.Push(git.WithRefSpecs(ctx.NextVersion.Raw),
		git.WithPushOptions(pushOpts...))
	return err
}

func printRepositoryTag(ctx *context.Context) {
	tags := make([]string, 0, 2)

	if ctx.PrintCurrentTag {
		tags = append(tags, ctx.CurrentVersion.Raw)
	}

	if ctx.PrintNextTag {
		tags = append(tags, ctx.NextVersion.Raw)
	}

	fmt.Fprint(ctx.Out, strings.Join(tags, " "))
}

func filterPushOptions(options []config.GitPushOption) []string {
	filtered := []string{}
	for _, opt := range options {
		if !opt.SkipTag {
			log.WithField("option", opt.Option).Debug("with push option")
			filtered = append(filtered, opt.Option)
		}
	}

	return filtered
}
