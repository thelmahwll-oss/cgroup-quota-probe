package platform

import "runtime"

// Info holds platform detection results used for fixture selection.
type Info struct {
	OS   string
	Arch string
}

// Current returns the current platform info.
func Current() Info {
	return Info{OS: runtime.GOOS, Arch: runtime.GOARCH}
}

// FixtureSuffix returns the test fixture suffix for the current platform.
func (i Info) FixtureSuffix() string {
	switch i.OS {
	case "linux":
		return "linux_amd64"
	case "darwin":
		if i.Arch == "arm64" {
			return "darwin_arm64"
		}
		return "darwin_amd64"
	case "windows":
		return "windows_amd64"
	default:
		return "unknown"
	}
}
