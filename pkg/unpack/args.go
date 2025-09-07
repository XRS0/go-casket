package unpack

import "path/filepath"

// buildArgs returns the command-line args for the chosen tool.
func buildArgs(use7z bool, archivePath, destDir string, overwrite bool, password string, sevenZType string, ignoreErrors bool) []string {
	if use7z {
		// 7z/7zz: x -y -o<dir> [-t<Type>] -p<PASS> <archive>
		args := []string{"x"}
		if overwrite {
			args = append(args, "-y")
		}
		args = append(args, "-o"+destDir)
		if sevenZType != "" {
			args = append(args, "-t"+sevenZType)
		}
		if password != "" {
			args = append(args, "-p"+password)
		}
		args = append(args, archivePath)
		return args
	}
	// unrar: x -o+/- -p<PASS>|-p- <archive> <destDir>
	args := []string{"x"}
	if overwrite {
		args = append(args, "-o+")
	} else {
		args = append(args, "-o-")
	}
	if password != "" {
		args = append(args, "-p"+password)
	} else {
		// Disable password prompt in non-interactive envs
		args = append(args, "-p-")
	}
	if ignoreErrors {
		args = append(args, "-kb")
	}
	// unrar requires normalized destination path
	args = append(args, archivePath, filepath.Clean(destDir))
	return args
}
