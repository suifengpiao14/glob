package glob

import (
	"io/fs"
	"strings"

	"github.com/gobwas/glob"
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
	g := glob.MustCompile(pattern)
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if g.Match(path) {
				matches = append(matches, path)
			}
		}
		return nil
	})
	return matches, err
}

func GlobURL(pattern string) ([]string, error) {

	starIndex := strings.Index(pattern, "*")
	if starIndex < 0 { // 没有*不枚举
		return []string{pattern}, nil
	}
	maxDepth := 0 // 默认不递归
	if strings.Contains(pattern, "**") {
		maxDepth = URL_enumerateFiles_maxDepth
	}
	u := pattern[:starIndex]
	lastSlash := strings.LastIndex(u, "/")
	if lastSlash > -1 {
		u = u[:lastSlash]
	}

	urls, err := EnumerateFiles(u, maxDepth, 0)
	if err != nil {
		return nil, err
	}

	var matches []string
	g := glob.MustCompile(pattern)
	for _, u := range urls {
		if g.Match(u) {
			matches = append(matches, u)
		}
	}
	return matches, err
}
