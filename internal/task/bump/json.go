package bump

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/config"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func jsonBump(ctx *context.Context, path string, bumps []config.JSONBump) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	str := string(data)

	for _, bump := range bumps {
		log.WithFields(log.Fields{
			"file":   path,
			"path":   bump.Path,
			"semver": bump.SemVer,
		}).Debug("attempting file bump")

		// Strip any 'v' prefix if this must be a semantic version
		v := ctx.NextVersion.Raw
		if bump.SemVer {
			v = strictSemVer(v)
		}

		res := gjson.Get(str, bump.Path)
		if !res.Exists() {
			return false, errors.New("no version matched in file")
		}

		str, err = sjson.SetOptions(str, bump.Path, v, &sjson.Options{Optimistic: true})
		if err != nil {
			return false, err
		}
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
