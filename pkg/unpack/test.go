package unpack

import (
	"context"
	"io"
	"os/exec"
	"path/filepath"
)

// TestArchive verifies archive integrity without extracting files.
// It auto-detects format by signature and chooses the best tool:
// - RAR/RAR5: prefer unrar if available, else 7z/7zz with -tRar
// - Others: 7z/7zz with -t<Type>
// It streams tool output to the provided writers if not nil, and respects ctx/opts.Timeout.
func TestArchive(ctx context.Context, archivePath string, opts Options, stdout, stderr io.Writer) error {
	if archivePath == "" {
		return ErrInvalidArgs
	}

	// Apply timeout if provided and no deadline is set yet.
	if opts.Timeout > 0 {
		if _, has := ctx.Deadline(); !has {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
			defer cancel()
		}
	}

	atype, _ := detectArchiveType(archivePath)
	sevenZType := mapSevenZType(atype)

	// Prefer unrar for RAR
	if sevenZType == "Rar" {
		if err := TestWithUnrar(ctx, archivePath, opts.Password, stdout, stderr); err == nil {
			return nil
		}
		// If unrar not available or failed, try 7z/7zz
		return TestWith7z(ctx, archivePath, opts.Password, sevenZType, stdout, stderr)
	}

	// Default: test with 7z/7zz
	return TestWith7z(ctx, archivePath, opts.Password, sevenZType, stdout, stderr)
}

// TestWith7z verifies archive using 7z/7zz: `7z t [-t<Type>] [-p<PASS>] <archive>`
func TestWith7z(ctx context.Context, archivePath string, password string, sevenZType string, stdout, stderr io.Writer) error {
	exe := ""
	if p, err := exec.LookPath(string(Tool7zz)); err == nil {
		exe = p
	} else if p, err := exec.LookPath(string(Tool7z)); err == nil {
		exe = p
	} else {
		return ErrNoToolFound
	}
	args := buildTestArgs(true, archivePath, password, sevenZType)
	return runTool(ctx, exe, args, filepath.Dir(archivePath), "", stdout, stderr)
}

// TestWithUnrar verifies archive using unrar: `unrar t -p<PASS>|-p- <archive>`
func TestWithUnrar(ctx context.Context, archivePath string, password string, stdout, stderr io.Writer) error {
	exe, err := exec.LookPath(string(ToolUnrar))
	if err != nil {
		return ErrNoToolFound
	}
	args := buildTestArgs(false, archivePath, password, "")
	return runTool(ctx, exe, args, filepath.Dir(archivePath), "", stdout, stderr)
}

// buildTestArgs builds args for tool test (integrity check) mode.
func buildTestArgs(use7z bool, archivePath string, password string, sevenZType string) []string {
	if use7z {
		// 7z/7zz: t [-t<Type>] [-p<PASS>] <archive>
		args := []string{"t"}
		if sevenZType != "" {
			args = append(args, "-t"+sevenZType)
		}
		if password != "" {
			args = append(args, "-p"+password)
		}
		args = append(args, archivePath)
		return args
	}
	// unrar: t -p<PASS>|-p- <archive>
	args := []string{"t", "-idq"}
	if password != "" {
		args = append(args, "-p"+password)
	} else {
		args = append(args, "-p-") // disable password prompt
	}
	args = append(args, archivePath)
	return args
}

// ValidateArchive is a convenience wrapper that uses background context.
func ValidateArchive(archivePath string, opts Options, stdout, stderr io.Writer) error {
	return TestArchive(context.Background(), archivePath, opts, stdout, stderr)
}
