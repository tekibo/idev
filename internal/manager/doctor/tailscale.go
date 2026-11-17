package doctor

import (
	"os"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/system"
)

func tailscaleCheck(cfg config.Config) Check {
	check := Check{
		ID:        "private.tailscale",
		Class:     "private",
		FixableBy: system.FixUser,
		Hint:      "run sudo systemctl enable --now tailscaled, tailscale up, and sudo tailscale set --operator=$USER",
	}
	if _, ok := cfg.FindBinary("tailscale"); !ok {
		check.Status = system.StatusMissing
		check.Detail = "tailscale is not on PATH"
		return check
	}
	if !tailscaleSocket() {
		check.Status = system.StatusWarn
		check.Detail = "tailscaled is not running"
		return check
	}
	check.Status = system.StatusOK
	return check
}

func tailscaleSocket() bool {
	for _, path := range []string{"/var/run/tailscale/tailscaled.sock", "/run/tailscale/tailscaled.sock"} {
		info, err := os.Stat(path)
		if err == nil && info.Mode()&os.ModeSocket != 0 {
			return true
		}
	}
	return false
}
