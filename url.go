package glob

import (
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var (
	URL_enumerateFiles_maxDepth = 15 // 递归目录最大深度
	URL_File_exts               = []string{".txt", ".md", ".sql", ".tpl", ".csv", ".json", ".html"}
)

func EnumerateFiles(url string, maxDepth int, depth int) (files []string, err error) {
	files = make([]string, 0)
	// 防止无限递归
	if depth > maxDepth {
		return
	}

	// 发送 HTTP 请求获取页面内容
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 解析页面内容
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return
		case html.SelfClosingTagToken, html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if isRelativeURL(link) {
							link = resolveRelativeURL(url, link)
						}
						if isFile(link) {
							files = append(files, link)
						} else {
							sub, err := EnumerateFiles(link, maxDepth, depth+1)
							if err != nil {
								return nil, err
							}
							files = append(files, sub...)

						}
					}
				}
			}
		}
	}
}

func isRelativeURL(url string) bool {
	return !strings.HasPrefix(url, "http")
}

func resolveRelativeURL(baseURL, url string) string {
	if strings.HasSuffix(baseURL, "/") {
		return baseURL + url
	}
	return baseURL + "/" + url
}

func isFile(url string) bool {
	for _, ext := range URL_File_exts {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}
