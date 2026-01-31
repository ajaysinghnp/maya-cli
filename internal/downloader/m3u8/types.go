package m3u8

import "github.com/ajaysinghnp/maya-cli/internal/logger/iface"

type Options struct {
	URL        string
	Output     string // final video file path
	TempDir    string // temporary folder for partial downloads
	Resume     bool
	Concurrent int
	Log        iface.Logger // or *logger.Logger depending on your interface design
}
