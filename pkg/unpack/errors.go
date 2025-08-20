package unpack

import "errors"

var (
	ErrInvalidArgs = errors.New("archivePath and destDir are required")
	ErrNoToolFound = errors.New("no supported archiver found on PATH (tried 7zz, 7z, unrar)")
)
