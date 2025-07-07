package logx

import "runtime/debug"

type buildInfo struct {
	GoVersion string `json:"go_version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
}

func newBuildInfo() buildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return buildInfo{}
	}

	var (
		gitCommit string
		buildDate string
		modified  bool
	)

	// Extract build settings
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			gitCommit = setting.Value
		case "vcs.time":
			buildDate = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	// Append +CHANGES to commit hash if repository was modified
	if modified {
		gitCommit += "+CHANGES"
	}

	return buildInfo{
		GoVersion: info.GoVersion,
		GitCommit: gitCommit,
		BuildDate: buildDate,
	}
}
