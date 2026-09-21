package buildinfo

// 构建时通过 ldflags 注入。
var (
	Mode      = "debug"
	Version   = "v1.0.0"
	GitInfo   = ""
	BuildTime = ""
)

const (
	ModeDebug   = "debug"
	ModeRelease = "release"
)

func IsRelease() bool {
	return Mode == ModeRelease
}
