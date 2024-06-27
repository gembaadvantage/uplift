package gitcheck

import (
	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
)

// Task for detecting if uplift is being run within a recognised
// git repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "checking git"
}

// Skip running the task
func (t Task) Skip(_ *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	repo, err := ctx.GitClient.Repository()
	if err != nil {
		return ErrNoRepository
	}

	log.Debug("checking if repository is dirty")
	status, err := ctx.GitClient.PorcelainStatus()
	if err != nil {
		return err
	}

	if len(status) > 0 {

		if len(ctx.IncludeArtifacts) == 0 {
			return ErrDirty{status}
		}

		for _, sts := range status {
			if !stringInSlice(sts.Path, ctx.IncludeArtifacts) {
				return ErrDirty{status}
			}
		}
	}

	log.Debug("checking for detached head")
	if repo.DetachedHead {
		if ctx.IgnoreDetached {
			log.Warn("detached HEAD detected. This may impact certain operations within uplift")
		} else {
			return ErrDetachedHead{}
		}
	}

	log.Debug("checking if shallow clone used")
	if repo.ShallowClone {
		if ctx.IgnoreShallow {
			log.Warn("shallow clone detected. This may impact certain operations within uplift")
		} else {
			return ErrShallowClone{}
		}
	}

	return nil
}

func stringInSlice(stringToFind string, listOfStrings []string) bool {
	for _, value := range listOfStrings {
		if value == stringToFind {
			return true
		}
	}
	return false
}
