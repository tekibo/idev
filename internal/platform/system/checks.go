package system

// spec describes one detection path.
type spec struct {
	ID       string
	Packages map[string]string
	Command  string
	Pkg      string
	Library  string
	Socket   string
	Critical bool
}

// pkgs maps one dependency to the package names of each manager.
func pkgs(pacman, apt, dnf string) map[string]string {
	return map[string]string{"pacman": pacman, "apt": apt, "dnf": dnf}
}

// Checks returns the system checks for the selected surfaces.
func Checks(probe *Probe, manager Manager, includeOptional bool) []Check {
	specs := coreSpecs()
	if includeOptional {
		specs = append(specs, optionalSpecs()...)
	}
	checks := make([]Check, 0, len(specs))
	for _, item := range specs {
		checks = append(checks, probeSpec(probe, manager, item))
	}
	return checks
}

func probeSpec(probe *Probe, manager Manager, item spec) Check {
	managerName := ""
	if manager != nil {
		managerName = manager.Name()
	}
	pkg := item.Packages[managerName]
	check := Check{
		ID:        item.ID,
		Class:     ClassSystem,
		Status:    StatusMissing,
		FixableBy: FixUser,
		Package:   pkg,
		Critical:  item.Critical,
	}
	if pkg != "" {
		check.Hint = pkg
	}
	if item.Command != "" {
		if path, ok := probe.Command(item.Command); ok {
			check.Status = StatusOK
			check.Detail = "command: " + path
			return check
		}
	}
	if item.Pkg != "" && probe.PkgConfig(item.Pkg) {
		check.Status = StatusOK
		check.Detail = "pkg-config: " + item.Pkg
		return check
	}
	if item.Library != "" && probe.Library(item.Library) {
		check.Status = StatusOK
		check.Detail = "library: " + item.Library
		return check
	}
	if item.Socket != "" && probe.Socket(item.Socket) {
		check.Status = StatusOK
		check.Detail = "socket: " + item.Socket
		return check
	}
	return check
}
