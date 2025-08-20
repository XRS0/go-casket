# go-casket

Small Go helper to unpack archives by calling system tools 7-Zip (7zz/7z) or unrar.

## Features
- Auto-detects available tool (prefers 7zz, then 7z, then unrar)
- Password-protected archives support (-p for 7z/7zz and unrar)
- Overwrite control
- Optional timeout
- Optional working directory
- Simple API and example CLI
- Dockerfile and docker-compose for quick testing

## Requirements (external tools)
This library shells out to one of the system tools. Install at least one of:
- 7-Zip CLI: `7zz` (preferred, official) or `7z` (p7zip)
- unrar (read-only RAR extractor)

Recommended: use 7zz (official 7-Zip) when possible; it supports newer formats like RAR5.

### macOS (Homebrew)
- 7-Zip (official):
  - `brew install sevenzip`  # provides `7zz`
- Optional RAR:
  - `brew install unrar`     # non-free

### Debian/Ubuntu
- Prefer official 7-Zip if available:
  - `sudo apt-get update && sudo apt-get install 7zip-full`  # may provide `7zz` on newer distros
- Or p7zip + unrar:
  - `sudo apt-get update && sudo apt-get install p7zip-full unrar`

### Windows
- Winget: `winget install 7zip.7zip`
- Chocolatey: `choco install 7zip` and optionally `choco install unrar`

## Install (Go module)
```bash
go get github.com/XRS0/go-casket
```
Go 1.22+ is required (see `go.mod`).

## Supported tools and formats
- 7zz (official 7-Zip): wide format support including 7z, zip, rar/rar5, tar (gz/bz2/xz/zst), wim, iso, cab, dmg, vhd/vhdx, vmdk and more.
- 7z (p7zip): similar but often lacks RAR5 and some newer formats.
- unrar: extraction of RAR archives only.

Actual support depends on your installed tool. You can run `7zz i` or `7z i` to list supported formats.

## Library usage
```go
import (
    "github.com/XRS0/go-casket/pkg/unpack"
)

func example() error {
    opts := unpack.DefaultOptions()
    // Optional settings
    opts.Overwrite = true
    opts.Timeout = 0           // no timeout
    opts.WorkDir = ""          // default to dest dir
    // Set password if the archive is encrypted
    opts.Password = "mySecret" // or leave empty if not needed

    return unpack.Unpack("/path/to/archive.rar", "/tmp/out", opts, nil, nil)
}
```

Notes:
- Use `unpack.UnpackContext(ctx, ...)` to apply a context with deadline/cancellation.
- By default, stdout/stderr of the tool are wired to the current process. Pass writers to capture output.
- `opts.Preferred` can be set to control tool preference order.

## Example CLI
Build and run example:

```bash
# from repo root
go run ./cmd/example /path/to/archive.zip ./out
# with password via 3rd arg
go run ./cmd/example /path/to/archive.7z ./out mySecret
```

Usage output:
```
Usage: example <archive> <destDir> [password]
```