package nextsemver

import (
	"fmt"
	"strings"

	semv "github.com/Masterminds/semver"
	"github.com/apex/log"
	git "github.com/purpleclay/gitz"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
)

// Task that determines the next semantic version of a repository
// based on the conventional commit of the last commit message
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "next semantic version"
}

// Skip is disabled for this task
func (t Task) Skip(_ *context.Context) bool {
	return false
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	var tagSuffix string
	if ctx.FilterOnPrerelease {
		tagSuffix = buildTagSuffix(ctx)
	}
	tag, err := latestTag(ctx.GitClient, tagSuffix)
	if err != nil {
		return err
	}

	if tag == "" {
		log.Debug("repository not tagged with version")
	} else {
		log.WithField("version", tag).Debug("identified latest version within repository")
	}
	ctx.CurrentVersion, _ = semver.Parse(tag)

	glog, err := ctx.GitClient.Log(git.WithRefRange(git.HeadRef, tag))
	if err != nil {
		return err
	}

	// Identify any commit that will trigger the largest semantic version bump
	inc := semver.ParseLogWithOptions(glog.Commits, semver.ParseOptions{TrimHeader: ctx.Changelog.TrimHeader})
	if inc == semver.NoIncrement {
		ctx.NoVersionChanged = true

		log.Warn("no commits trigger a change in semantic version")
		return nil
	}
	log.WithField("increment", string(inc)).Info("largest increment detected from commits")

	if tag == "" {
		tag = "v0.0.0"
	}

	// Remove the prefix if needed
	if ctx.NoPrefix {
		tag = strings.TrimPrefix(tag, "v")
	}

	pver, _ := semv.NewVersion(tag)

	// Handle the fact semver returns a pointer when it initialises a new
	// semv.Version, but all of its methods work on a copy
	nxt := *pver

	if ctx.IgnoreExistingPrerelease {
		log.Info("stripped existing prerelease metadata from version")
		nxt, _ = nxt.SetPrerelease("")
		nxt, _ = nxt.SetMetadata("")
	}

	switch inc {
	case semver.MajorIncrement:
		nxt = nxt.IncMajor()
	case semver.MinorIncrement:
		nxt = nxt.IncMinor()
	case semver.PatchIncrement:
		nxt = nxt.IncPatch()
	}

	// Append any prerelease suffixes
	if ctx.Prerelease != "" {
		nxt, _ = nxt.SetPrerelease(ctx.Prerelease)
		nxt, _ = nxt.SetMetadata(ctx.Metadata)

		log.WithFields(log.Fields{
			"prerelease": ctx.Prerelease,
			"metadata":   ctx.Metadata,
		}).Debug("appending prerelease version")
	}

	ctx.NextVersion = semver.Version{
		Prefix:     ctx.CurrentVersion.Prefix,
		Patch:      nxt.Patch(),
		Minor:      nxt.Minor(),
		Major:      nxt.Major(),
		Prerelease: nxt.Prerelease(),
		Metadata:   nxt.Metadata(),
		Raw:        nxt.Original(),
	}

	log.WithField("version", ctx.NextVersion.Raw).Info("identified next semantic version")
	return nil
}

func latestTag(gitc *git.Client, suffix string) (string, error) {
	tags, err := gitc.Tags(git.WithShellGlob("*.*.*"),
		git.WithSortBy(git.CreatorDateDesc, git.VersionDesc))
	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", nil
	}

	if suffix == "" {
		return tags[0], nil
	}

	for _, t := range tags {
		if strings.HasSuffix(t, suffix) {
			return t, nil
		}
	}

	return "", nil
}

func buildTagSuffix(ctx *context.Context) string {
	var suffix string
	if ctx.Prerelease != "" {
		suffix = fmt.Sprintf("-%s", ctx.Prerelease)
		if ctx.Metadata != "" {
			suffix = fmt.Sprintf("%s+%s", suffix, ctx.Metadata)
		}
	}
	return suffix
}
