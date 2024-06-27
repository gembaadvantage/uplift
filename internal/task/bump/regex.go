package bump

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/semver"
)

func regexBump(ctx *context.Context, path string, bumps []config.RegexBump) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	str := string(data)

	for _, bump := range bumps {
		log.WithFields(log.Fields{
			"file":   path,
			"regex":  bump.Pattern,
			"count":  bump.Count,
			"semver": bump.SemVer,
		}).Debug("attempting file bump")

		mstr, err := match(bump.Pattern, str)
		if err != nil {
			return false, err
		}

		if strings.Contains(mstr, ctx.NextVersion.Raw) {
			log.WithFields(log.Fields{
				"file":    path,
				"current": ctx.CurrentVersion.Raw,
				"next":    ctx.NextVersion.Raw,
			}).Info("skipping bump")
			return false, nil
		}

		// Use strings replace to ensure the replacement count is honoured
		n := -1
		if bump.Count > 0 {
			n = bump.Count
		}

		// Strip any 'v' prefix if this must be a semantic version
		v := ctx.NextVersion.Raw
		if bump.SemVer {
			v = strictSemVer(v)
		}

		verRpl := semver.Regex.ReplaceAllString(mstr, v)
		str = strings.Replace(str, mstr, verRpl, n)
	}

	log.WithFields(log.Fields{
		"file":    path,
		"current": ctx.CurrentVersion.Raw,
		"next":    ctx.NextVersion.Raw,
	}).Info("file bumped")

	// Don't make any file changes if part of a dry-run
	if ctx.DryRun {
		log.Info("file not modified in dry run mode")
		return false, nil
	}

	return true, os.WriteFile(path, []byte(str), 0o644)
}

func match(pattern string, data string) (string, error) {
	verRgx := strings.Replace(pattern, semver.Token, semver.Pattern, 1)

	rgx, err := regexp.Compile(verRgx)
	if err != nil {
		return "", err
	}

	m := rgx.FindString(data)
	if m == "" {
		return "", errors.New("no version matched in file")
	}

	return string(m), nil
}

func strictSemVer(v string) string {
	if v[0] == 'v' {
		log.Debug("trimming 'v' prefix from version")
		v = v[1:]
	}
	return v
}
