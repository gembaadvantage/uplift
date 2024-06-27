package beforetag

import (
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/task/hook"
)

// Task for executing any custom shell commands or scripts
// before tagging within the release workflow
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "before tagging repository"
}

// Skip running the task
func (t Task) Skip(ctx *context.Context) bool {
	return ctx.Config.Hooks == nil || len(ctx.Config.Hooks.BeforeTag) == 0 || ctx.NoVersionChanged
}

// Run the task
func (t Task) Run(ctx *context.Context) error {
	return hook.Exec(ctx.Context, ctx.Config.Hooks.BeforeTag, hook.ExecOptions{
		DryRun: ctx.DryRun,
		Debug:  ctx.Debug,
		Env:    ctx.Config.Env,
	})
}
