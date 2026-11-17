package privileged

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/terminal"
)

const sudoPrompt = "idev-sudo-prompt: "

// authenticate runs sudo -v and the real command in one PTY session. sudo
// caches credentials per terminal by default, so the command must use the
// terminal that authenticated.
func (r *Runner) authenticate(ctx context.Context, plan Plan) error {
	password, err := r.Provider.Password(ctx)
	if err != nil {
		return err
	}
	defer zero(password)
	if len(password) == 0 {
		return ErrPasswordRequired
	}
	pair, err := terminal.Open()
	if err != nil {
		return err
	}
	defer func() { _ = pair.Close() }()
	if err := pair.DisableEcho(); err != nil {
		return err
	}
	authCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd := execx.NewCommandContext(authCtx, execx.Command{
		Name: "/bin/sh",
		Args: []string{"-c", r.script(plan)},
		Dir:  plan.Command.Dir,
		Env:  plan.Command.Env,
	})
	if err := pair.Start(cmd); err != nil {
		return err
	}
	wait := make(chan error, 1)
	go func() { wait <- cmd.Wait() }()
	out := plan.Command.Stdout
	if out == nil {
		out = plan.Command.Stderr
	}
	if out == nil {
		out = io.Discard
	}
	prompts, done := watchPrompts(pair.Master, []byte(sudoPrompt), out)
	defer close(done)

	seen, err := waitForPrompt(ctx, r.promptTimeout(), prompts, wait)
	if err != nil {
		return err
	}
	if !seen {
		return nil
	}
	line := make([]byte, len(password)+1)
	copy(line, password)
	line[len(password)] = '\n'
	_, writeErr := pair.Master.Write(line)
	zero(line)
	if writeErr != nil {
		return writeErr
	}
	return waitForResult(ctx, prompts, wait)
}

// script builds the one shell session that authenticates and runs the command.
func (r *Runner) script(plan Plan) string {
	auth := shellJoin([]string{r.sudo(), "-p", sudoPrompt, "-v"})
	parts := append([]string{r.sudo(), "-n", plan.Command.Name}, plan.Command.Args...)
	return auth + " || exit 125; " + shellJoin(parts) + "; exit $?"
}

// shellQuote quotes one argument for /bin/sh. Every argument is quoted, so
// shell metacharacters never become syntax.
func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func shellJoin(parts []string) string {
	quoted := make([]string, 0, len(parts))
	for _, part := range parts {
		quoted = append(quoted, shellQuote(part))
	}
	return strings.Join(quoted, " ")
}

func (r *Runner) promptTimeout() time.Duration {
	if r.PromptTimeout > 0 {
		return r.PromptTimeout
	}
	return 30 * time.Second
}

func zero(data []byte) {
	for index := range data {
		data[index] = 0
	}
}
