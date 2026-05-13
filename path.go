package httpfs

import (
	"path/filepath"
	"strings"

	libpath "github.com/reiver/go-path"
)

// PathJoin returns the path resulting from joining the root path and the sub-path.
//
// Typcally the root-path is a directory on the file-system,
// and the sub-path comes from the HTTP request.
func PathJoin(root string, subpath string) (path string, ok bool) {
	if "" == root {
		return "", false
	}

	path = filepath.Join(root, libpath.Canonical(subpath))
	if "" == path {
		return "", false
	}

	cleanedRoot := filepath.Clean(root)
	if path != cleanedRoot && !strings.HasPrefix(path, cleanedRoot+string(filepath.Separator)) {
		return "", false
	}

	return path, true
}

// pathResolve resolves symlinks in path and verifies the result is still under root.
// Returns the resolved path, or "" if resolution fails or the path escapes root.
//
// Typcally the root-path is a directory on the file-system,
// and the sub-path comes from the HTTP request.
//
// Use this for files that must already exist (GET, HEAD, DELETE).
func pathResolve(root string, path string) string {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if nil != err {
		return ""
	}

	resolved, err := filepath.EvalSymlinks(path)
	if nil != err {
		return ""
	}

	if resolved != resolvedRoot && !strings.HasPrefix(resolved, resolvedRoot+string(filepath.Separator)) {
		return ""
	}

	return resolved
}

// pathResolveParent resolves symlinks in the nearest existing ancestor of path
// and verifies the result is still under root. Returns the path with the resolved
// ancestor and the remaining (non-existent) segments appended, or "" if resolution
// fails or the path escapes root.
//
// Typcally the root-path is a directory on the file-system,
// and the sub-path comes from the HTTP request.
//
// Use this for files that may not exist yet (PUT).
func pathResolveParent(root string, path string) string {
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if nil != err {
		return ""
	}

	// Walk up from path until we find an existing ancestor.
	existing := path
	var remaining []string
	for {
		_, err := filepath.EvalSymlinks(existing)
		if nil == err {
			break
		}
		base := filepath.Base(existing)
		if "." == base || ".." == base {
			return ""
		}
		remaining = append(remaining, base)
		parent := filepath.Dir(existing)
		if parent == existing {
			return ""
		}
		existing = parent
	}

	resolvedExisting, err := filepath.EvalSymlinks(existing)
	if nil != err {
		return ""
	}

	if resolvedExisting != resolvedRoot && !strings.HasPrefix(resolvedExisting, resolvedRoot+string(filepath.Separator)) {
		return ""
	}

	// Re-append the non-existent segments in the correct order.
	result := resolvedExisting
	for i := len(remaining) - 1; i >= 0; i-- {
		result = filepath.Join(result, remaining[i])
	}

	return result
}
