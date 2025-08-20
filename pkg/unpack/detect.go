package unpack

import (
	"os/exec"
)

// DetectTool finds the first available tool from the list; returns the executable name.
func DetectTool(preferred []Tool) (string, error) {
	candidates := preferred
	if len(candidates) == 0 {
		candidates = []Tool{Tool7zz, Tool7z, ToolUnrar}
	}
	for _, t := range candidates {
		if exe, err := exec.LookPath(string(t)); err == nil {
			return exe, nil
		}
	}
	return "", ErrNoToolFound
}

func is7z(exe string) bool {
	base := exe
	// exec.LookPath returns full path; keep simple check by suffix
	if len(base) >= 3 {
		if base[len(base)-3:] == "7zz" || base[len(base)-2:] == "7z" {
			return true
		}
	}
	return false
}
