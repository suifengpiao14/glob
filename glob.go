package glob

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func GlobDirectory(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}
	var matches []string
	reg := getPattern(pattern)
	index := strings.Index(pattern, "**")
	root := pattern[:index]
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			err := errors.Errorf("dir:%s filepath.Walk info is nil", ".")
			return err
		}
		if !info.IsDir() {
			path = strings.ReplaceAll(path, "\\", "/")
			if reg.MatchString(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
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
	regStr = strings.ReplaceAll(regStr, ".", "\\.")
	regStr = strings.ReplaceAll(regStr, "**", ".*")
	reg := regexp.MustCompile(regStr)
	return reg
}
