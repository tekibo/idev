package doctor

import (
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/system"
)

func flowCheck(cfg config.Config) Check {
	selection, err := config.LoadFlows(cfg.Home)
	detail := ""
	ok := err == nil && len(selection.Flows) > 0
	if err != nil {
		detail = err.Error()
	} else if !ok {
		detail = "no saved flows"
	}
	return Check{
		ID:        "private.flows",
		Class:     "private",
		Status:    status(ok),
		FixableBy: system.FixIdev,
		Detail:    detail,
		Hint:      "run idev setup --flow <name>",
	}
}

func privateTool(id string, cfg config.Config, name, hint string) Check {
	path := cfg.Path("bin", name)
	ok := config.IsExec(path)
	detail := ""
	if ok {
		detail = path
	}
	return Check{
		ID:        id,
		Class:     "private",
		Status:    status(ok),
		FixableBy: system.FixIdev,
		Detail:    detail,
		Hint:      hint,
		Critical:  true,
	}
}

func privateFile(id, path, hint string) Check {
	ok := path != "" && config.IsFile(path)
	detail := ""
	if ok {
		detail = path
	}
	return Check{
		ID:        id,
		Class:     "private",
		Status:    status(ok),
		FixableBy: system.FixIdev,
		Detail:    detail,
		Hint:      hint,
		Critical:  true,
	}
}

func privateOptionalFile(id, path, hint string) Check {
	check := privateFile(id, path, hint)
	check.Critical = false
	return check
}

func projectFile(id, path, hint string) Check {
	ok := path != "" && config.IsExec(path)
	detail := ""
	if ok {
		detail = path
	}
	return Check{
		ID:        id,
		Class:     "project",
		Status:    status(ok),
		FixableBy: system.FixUser,
		Detail:    detail,
		Hint:      hint,
		Critical:  true,
	}
}

func status(ok bool) string {
	if ok {
		return system.StatusOK
	}
	return system.StatusMissing
}
