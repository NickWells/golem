package xdg

import (
	"os"
	"strings"
)

// The values given here represent the environment variables and default
// values from the XDG Base Directory Specification. See that spec for more
// details about how to use these values.
const (
	DataHomeEnvVar     = "XDG_DATA_HOME"
	DataHomeEnvVarDflt = "$HOME/.local/share"

	ConfigHomeEnvVar     = "XDG_CONFIG_HOME"
	ConfigHomeEnvVarDflt = "$HOME/.config"

	DataDirsEnvVar     = "XDG_DATA_DIRS"
	DataDirsEnvVarDflt = "/usr/local/share/:/usr/share/"

	ConfigDirsEnvVar     = "XDG_CONFIG_DIRS"
	ConfigDirsEnvVarDflt = "/etc/xdg"

	CacheHomeEnvVar     = "XDG_CACHE_HOME"
	CacheHomeEnvVarDflt = "$HOME/.cache"

	ListSep = ":"
)

// DataHome returns the value of the XDG_DATA_HOME environment variable or
// the default value
func DataHome() string {
	rval := os.ExpandEnv("$" + DataHomeEnvVar)
	if rval == "" {
		rval = os.ExpandEnv(DataHomeEnvVarDflt)
	}

	return rval
}

// ConfigHome returns the value of the XDG_CONFIG_HOME environment variable or
// the default value
func ConfigHome() string {
	rval := os.ExpandEnv("$" + ConfigHomeEnvVar)
	if rval == "" {
		rval = os.ExpandEnv(ConfigHomeEnvVarDflt)
	}

	return rval
}

// DataDirs returns the values in the XDG_DATA_DIRS environment variable or
// the default values
func DataDirs() []string {
	val := os.ExpandEnv("$" + DataDirsEnvVar)
	if val == "" {
		val = DataDirsEnvVarDflt
	}

	return strings.Split(val, ListSep)
}

// ConfigDirs returns the values in the XDG_CONFIG_DIRS environment variable or
// the default values
func ConfigDirs() []string {
	val := os.ExpandEnv("$" + ConfigDirsEnvVar)
	if val == "" {
		val = ConfigDirsEnvVarDflt
	}

	return strings.Split(val, ListSep)
}

// CacheHome returns the value of the XDG_CACHE_HOME environment variable or
// the default value
func CacheHome() string {
	rval := os.ExpandEnv("$" + CacheHomeEnvVar)
	if rval == "" {
		rval = os.ExpandEnv(CacheHomeEnvVarDflt)
	}

	return rval
}
