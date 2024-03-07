package path

import "strings"

func ExtractPrefix(url string) string {
	if url == "" || url == "/" {
		return "/"
	}

	pathPrefix := url
	parts := strings.Split(pathPrefix, "/")
	if len(parts) > 1 {
		pathPrefix = "/" + parts[1]
	}
	return pathPrefix
}
