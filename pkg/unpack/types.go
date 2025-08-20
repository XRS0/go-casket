package unpack

import "time"

// Tool represents an external archive tool we can call.
type Tool string

const (
	Tool7z    Tool = "7z"
	Tool7zz   Tool = "7zz"
	ToolUnrar Tool = "unrar"
)

// Options control how extraction is performed.
type Options struct {
	// Tool preference order; first found on PATH will be used.
	Preferred []Tool
	// If true, overwrite existing files without prompt.
	Overwrite bool
	// Extraction timeout. Zero means no timeout.
	Timeout time.Duration
	// Working directory for command execution. If empty, uses output dir.
	WorkDir string
	// Password for encrypted archives (7z/zip/rar). Empty means no password is passed.
	Password string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Preferred: []Tool{Tool7zz, Tool7z, ToolUnrar},
		Overwrite: true,
		Timeout:   0,
	}
}
