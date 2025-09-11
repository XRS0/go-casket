package unpack

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	// Archive type detection via Kaitai-generated parser
	"github.com/XRS0/go-casket/pkg/unpack/tools"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
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

	// Detect archive type from signature
	atype, _ := detectArchiveType(archivePath)
	sevenZType := mapSevenZType(atype)

	exe, err := DetectTool(opts.Preferred)
	if err != nil {
		return err
	}

	use7z := is7z(exe)
	args := buildArgs(use7z, archivePath, destDir, opts.Overwrite, opts.Password, sevenZType, opts.IgnoreErrors)

	// Apply timeout if provided and no deadline is set yet.
	if opts.Timeout > 0 && ctx.Err() == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// If detected RAR, prefer unrar when available, else use 7z with -tRar. If 7z fails, try unrar as fallback.
	if atype == tools.ArchiveHeader_ArchiveType__Rar4 || atype == tools.ArchiveHeader_ArchiveType__Rar5 {
		if unrarExe, lookErr := exec.LookPath(string(ToolUnrar)); lookErr == nil {
			urArgs := buildArgs(false, archivePath, destDir, opts.Overwrite, opts.Password, "", opts.IgnoreErrors)
			if err := runTool(ctx, unrarExe, urArgs, destDir, opts.WorkDir, stdout, stderr); err != nil {
				return fmt.Errorf("extract failed with unrar: %w", err)
			}
			return nil
		}
		// no unrar, proceed with 7z using -tRar
	}

	// Run the chosen tool (7z/7zz or whatever DetectTool returned)
	if err := runTool(ctx, exe, args, destDir, opts.WorkDir, stdout, stderr); err != nil {
		// Fallback: if 7z/7zz failed and archive is RAR, try unrar when available
		if use7z && (atype == tools.ArchiveHeader_ArchiveType__Rar4 || atype == tools.ArchiveHeader_ArchiveType__Rar5) {
			if unrarExe, lookErr := exec.LookPath(string(ToolUnrar)); lookErr == nil {
				fallbackArgs := buildArgs(false, archivePath, destDir, opts.Overwrite, opts.Password, "", opts.IgnoreErrors)
				if fbErr := runTool(ctx, unrarExe, fallbackArgs, destDir, opts.WorkDir, stdout, stderr); fbErr == nil {
					return nil
				}
			}
		}
		return fmt.Errorf("extract failed with %s: %w", filepath.Base(exe), err)
	}
	return nil
}

// runTool executes the external extraction tool with common wiring.
func runTool(ctx context.Context, exe string, args []string, destDir, workDir string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, exe, args...)
	if workDir != "" {
		cmd.Dir = workDir
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
		return err
	}
	return nil
}

// detectArchiveType reads magic and returns the tools.ArchiveHeader_ArchiveType
func detectArchiveType(path string) (tools.ArchiveHeader_ArchiveType, error) {
	f, err := os.Open(path)
	if err != nil {
		return tools.ArchiveHeader_ArchiveType__Unknown, err
	}
	defer f.Close()
	ks := kaitai.NewStream(f)
	ah := tools.NewArchiveHeader()
	if err := ah.Read(ks, ah, ah); err != nil {
		return tools.ArchiveHeader_ArchiveType__Unknown, err
	}
	return ah.Archive()
}

// mapSevenZType maps detected type to 7z -t<Type> values
func mapSevenZType(t tools.ArchiveHeader_ArchiveType) string {
	switch t {
	case tools.ArchiveHeader_ArchiveType__Zip:
		return "Zip"
	case tools.ArchiveHeader_ArchiveType__Rar4:
		return "Rar"
	case tools.ArchiveHeader_ArchiveType__Rar5:
		return ""
	case tools.ArchiveHeader_ArchiveType__SevenZ:
		return "7z"
	case tools.ArchiveHeader_ArchiveType__Gzip:
		return "GZip"
	case tools.ArchiveHeader_ArchiveType__Bzip2:
		return "BZip2"
	case tools.ArchiveHeader_ArchiveType__Xz:
		return "XZ"
	case tools.ArchiveHeader_ArchiveType__Cab:
		return "Cab"
	case tools.ArchiveHeader_ArchiveType__Arj:
		return "Arj"
	case tools.ArchiveHeader_ArchiveType__Tar:
		return "Tar"
	case tools.ArchiveHeader_ArchiveType__Unknown:
		fallthrough
	default:
		return ""
	}
}
