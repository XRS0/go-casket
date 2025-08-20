package unpack

import "path/filepath"

// buildArgs returns the command-line args for the chosen tool.
func buildArgs(use7z bool, archivePath, destDir string, overwrite bool, password string) []string {
	if use7z {
		// 7z x -y -o<dir> -p<PASS> -- <archive>
		args := []string{"x"}
		if overwrite {
			args = append(args, "-y")
		}
		args = append(args, "-o"+destDir)
		if password != "" {
			args = append(args, "-p"+password)
		}
		args = append(args, "--", archivePath)
		return args
	}
	// unrar x -o+ -p<PASS> <archive> <destDir>
	args := []string{"x"}
	if overwrite {
		args = append(args, "-o+")
	} else {
		args = append(args, "-o-")
	}
	if password != "" {
		args = append(args, "-p"+password)
	}
	// unrar requires normalized destination path
	args = append(args, archivePath, filepath.Clean(destDir))
	return args
}
