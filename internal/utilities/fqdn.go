package utilities

import "regexp"

func GetFQDN(url string) string {
	r := regexp.MustCompile(`http(s)?:\/\/`)

	return r.ReplaceAllString(url, "")
}
