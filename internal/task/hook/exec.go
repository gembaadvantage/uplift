package hook

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/joho/godotenv"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

var cleanEnv = regexp.MustCompile(`\s*=\s*`)

// ExecOptions provides a way of customising the execution of commands
type ExecOptions struct {
	Debug  bool
	DryRun bool
	Env    []string
}

// Exec will execute a series of shell commands or scripts
func Exec(_ context.Context, cmds []string, opts ExecOptions) error {
	env := os.Environ()
	if len(opts.Env) > 0 {
		renv, err := resolveEnv(opts.Env)
		if err != nil {
			return err
		}

		env = append(env, renv...)
	}

	for _, c := range cmds {
		log.WithField("hook", c).Info("running")
		if opts.DryRun {
			continue
		}

		p, err := syntax.NewParser().Parse(strings.NewReader(c), "")
		if err != nil {
			return err
		}

		// Discard all output from commands and scripts unless in debug mode
		out := io.Discard
		if opts.Debug {
			// Stderr is used by apex for logging, stdout is reserved for capturing output
			out = os.Stderr
		}

		r, err := interp.New(
			interp.Params("-e"),
			interp.StdIO(os.Stdin, out, os.Stderr),
			interp.OpenHandler(openHandler),
			interp.Env(expand.ListEnviron(env...)),
		)
		if err != nil {
			return err
		}

		if err := r.Run(context.Background(), p); err != nil {
			return err
		}
	}

	return nil
}

func openHandler(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if path == "/dev/null" {
		return DevNull{}, nil
	}

	return interp.DefaultOpenHandler()(ctx, path, flag, perm)
}

func resolveEnv(env []string) ([]string, error) {
	renv := make([]string, 0, len(env))
	for _, e := range env {
		// Check for a .env extension
		if strings.HasSuffix(e, ".env") {
			_, err := os.Stat(e)
			if err != nil {
				return []string{}, fmt.Errorf("file %s does not exist", e)
			}

			log.WithField("path", e).Debug("Loading dotenv file")

			dotenv, err := godotenv.Read(e)
			if err != nil {
				return []string{}, err
			}

			for k, v := range dotenv {
				// A breaking change was introduced during the release of godotenv 1.5.0
				// that parsing an env file without a key e.g. '=VALUE' no longer returns
				// an error. Ensure this behaviour persists within uplift
				if k == "" {
					return []string{}, errors.New("can't separate key from value")
				}

				denv := fmt.Sprintf("%s=%s", k, v)

				log.WithField("var", denv).Debug("Injecting env var")
				renv = append(renv, denv)
			}
		} else {
			// Crude sanitisation of env var, trimming any whitespace around the assignment '='
			// otherwise it will be silently ignored
			e = cleanEnv.ReplaceAllString(e, "=")

			log.WithField("var", e).Debug("Injecting env var")
			renv = append(renv, e)
		}
	}

	return renv, nil
}
