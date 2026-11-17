package xtool

import (
	"os"

	"github.com/tekibo/idev/internal/platform/config"
)

func rollback(cfg config.Config) {
	_ = os.RemoveAll(cfg.XtoolRoot())
	_ = os.Remove(cfg.XtoolLink())
	_ = os.RemoveAll(cfg.XtoolSource())
}
