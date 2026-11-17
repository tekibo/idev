package help

func ipadScreenCommand() *Command {
	return &Command{
		Name:    "ipad-screen",
		Group:   "Devices",
		Summary: "Mirror and control a USB iPad over Tailscale.",
		Usage:   "idev ipad-screen <start|stop|status|run|install|menu> [options]",
		Description: "Starts a password-protected HTTPS viewer on loopback and publishes it through Tailscale Serve. " +
			"The iPad must be connected by USB with Developer Mode and trust enabled. Config and state stay under IDEV_HOME.",
		Children: []*Command{
			{
				Name:    "start",
				Summary: "Start the viewer and publish it through Tailscale Serve.",
				Usage:   "idev ipad-screen start --udid UDID [--port PORT] [--upgrade]",
				Options: []Option{
					req("--udid UDID", "USB device UDID."),
					def("--port PORT", "Loopback HTTPS port published by Serve.", "8080"),
					opt("--upgrade", "Upgrade the private pymobiledevice3 venv if it is too old."),
				},
			},
			{
				Name:    "stop",
				Summary: "Stop the viewer and remove its exact Serve mapping.",
				Usage:   "idev ipad-screen stop",
			},
			{
				Name:    "status",
				Summary: "Show the viewer state and published URL.",
				Usage:   "idev ipad-screen status",
			},
			{
				Name:    "run",
				Summary: "Run the viewer in the foreground for the user service.",
				Usage:   "idev ipad-screen run --udid UDID [--port PORT]",
				Options: []Option{
					req("--udid UDID", "USB device UDID."),
					def("--port PORT", "Loopback HTTPS port published by Serve.", "8080"),
				},
			},
			{
				Name:    "install",
				Summary: "Install the command link, menu entry, and disabled user unit.",
				Usage:   "idev ipad-screen install",
			},
			{
				Name:    "menu",
				Summary: "Run the interactive start, stop, or status prompt.",
				Usage:   "idev ipad-screen menu",
			},
		},
		Examples: []Example{
			example("Start the viewer for one iPad", "idev ipad-screen start --udid 00008110-001234567890001E"),
			example("Check the published URL", "idev ipad-screen status"),
		},
	}
}
