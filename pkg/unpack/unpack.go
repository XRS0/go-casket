package unpack

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Unpack is a convenience wrapper around UnpackContext with context.Background().
func Unpack(archivePath, destDir string, opts Options, stdout, stderr io.Writer) error {
	return UnpackContext(context.Background(), archivePath, destDir, opts, stdout, stderr)
}

// UnpackContext extracts an archive file to destDir using 7z/7zz or unrar.
// It streams tool output to the provided writers if not nil, and respects the provided context.
func UnpackContext(ctx context.Context, archivePath, destDir string, opts Options, stdout, stderr io.Writer) error {
	if archivePath == "" || destDir == "" {
		return ErrInvalidArgs
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("create dest dir: %w", err)
	}

	exe, err := DetectTool(opts.Preferred)
	if err != nil {
		return err
	}

	use7z := is7z(exe)
	args := buildArgs(use7z, archivePath, destDir, opts.Overwrite, opts.Password)

	// Apply timeout if provided and no deadline is set yet.
	if opts.Timeout > 0 && ctx.Err() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, exe, args...)
	if opts.WorkDir != "" {
		cmd.Dir = opts.WorkDir
	} else {
		cmd.Dir = destDir
	}

	// Ensure environment doesn't localize prompts
	cmd.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")

	// Wire output
	if stdout != nil {
		cmd.Stdout = stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	} else {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %s: %w", exe, err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("extract failed with %s: %w", filepath.Base(exe), err)
	}
	return nil
}
