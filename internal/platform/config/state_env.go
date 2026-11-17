package config

// StateEnv returns the persistent state paths kept under IDEV_HOME.
// The native manager keeps the same shape under IDEV_HOME so xtool, Expo and
// CocoaPods see one stable state root across runs.
func (c Config) StateEnv() map[string]string {
	return map[string]string{
		"XDG_CONFIG_HOME": c.Path(".config"),
		"XDG_DATA_HOME":   c.Path(".local", "share"),
		"XDG_CACHE_HOME":  c.Path(".cache"),
		"SDK_CACHE":       c.Path(".cache", "xtool"),
	}
}
