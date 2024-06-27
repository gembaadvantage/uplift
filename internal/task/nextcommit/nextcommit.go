package nextcommit

import (
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
	git "github.com/purpleclay/gitz"
)

// Task for generating the next commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "building next commit"
}

// Skip is disabled for this task
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged
}

// Run the task and generate the next commit by either impersonating the author
// from the last commit or by generating a user defined commit
func (t Task) Run(ctx *context.Context) error {
	c := git.CommitDetails{
		Author: git.Person{
			Name:  "uplift-bot",
			Email: "uplift@gembaadvantage.com",
		},
		Message: fmt.Sprintf("ci(uplift): uplifted for version %s", ctx.NextVersion.Raw),
	}

	if ctx.Config.CommitAuthor != nil {
		if ctx.Config.CommitAuthor.Name != "" {
			log.Debug("overwriting commit author name from uplift config")
			c.Author.Name = ctx.Config.CommitAuthor.Name
		}

		if ctx.Config.CommitAuthor.Email != "" {
			log.Debug("overwriting commit author email from uplift config")
			c.Author.Email = ctx.Config.CommitAuthor.Email
		}
	}

	if ctx.Config.CommitMessage != "" {
		log.Debug("overwriting commit message from uplift config")
		c.Message = strings.ReplaceAll(ctx.Config.CommitMessage, semver.Token, ctx.NextVersion.Raw)
	}

	cfg, err := ctx.GitClient.Config()
	if err != nil {
		return err
	}

	name := cfg["user.name"]
	email := cfg["user.email"]
	if name != "" && email != "" {
		log.Debug("overwriting commit author from git config")
		c.Author.Name = name
		c.Author.Email = email
	}

	ctx.CommitDetails = c
	log.WithFields(log.Fields{
		"name":    ctx.CommitDetails.Author.Name,
		"email":   ctx.CommitDetails.Author.Email,
		"message": ctx.CommitDetails.Message,
	}).Info("changes will be committed with")
	return nil
}
