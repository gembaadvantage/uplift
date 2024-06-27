package bump

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/goreleaser/fileglob"
	git "github.com/purpleclay/gitz"
)

// Task for bumping versions within files
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "bumping files"
}

// Skip running the task if no version has changed
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.NoVersionChanged || ctx.SkipBumps
}

// Run the task bumping the semantic version of any file identified within
// the uplift configuration file
func (t Task) Run(ctx *context.Context) error {
	if len(ctx.Config.Bumps) == 0 {
		log.Info("no files to bump")
		return nil
	}

	n := 0
	for _, bump := range ctx.Config.Bumps {
		// For simplicity, treat all file paths as Globs. If no glob characters are detected, a []string
		// will still be returned for the individual file, providing a consistent approach to bumping.
		resolvedBumps, err := resolveGlob(bump.File)
		if err != nil {
			return err
		}

		for _, resolvedBump := range resolvedBumps {
			var ok bool
			var bumpErr error

			// A unified path based approach to file bumping is definitely needed here.
			if len(bump.Regex) > 0 {
				ok, bumpErr = regexBump(ctx, resolvedBump, bump.Regex)
			}

			if len(bump.JSON) > 0 {
				ok, bumpErr = jsonBump(ctx, resolvedBump, bump.JSON)
			}

			if bumpErr != nil {
				return bumpErr
			}

			if ok {
				n++
			}

			if ctx.NoStage {
				log.Info("skip staging of file")
				continue
			}

			if _, err := ctx.GitClient.Stage(git.WithPathSpecs(resolvedBump)); err != nil {
				return err
			}
			log.WithField("file", resolvedBump).Info("successfully staged file")
		}
	}

	if n > 0 {
		log.WithField("count", n).Debug("bumped files staged")
	}

	return nil
}

func resolveGlob(pattern string) ([]string, error) {
	if !fileglob.ContainsMatchers(pattern) {
		return []string{pattern}, nil
	}

	return fileglob.Glob(pattern)
}
