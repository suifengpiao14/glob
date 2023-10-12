package glob

import (
	"io/fs"
	"regexp"
	"strings"

	"github.com/yargevad/filepathx"
)

func GlobDirectory(pattern string) ([]string, error) {
	matches, err := filepathx.Glob(pattern)
	return matches, err
}

func GlobFS(fsys fs.FS, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return fs.Glob(fsys, pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func getPattern(pattern string) *regexp.Regexp {
	regStr := strings.TrimLeft(pattern, ".")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	return reg
}
