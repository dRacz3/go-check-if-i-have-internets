package main

import (
	"io"
	"net/http"
	"strings"
)

type LinkChecker struct {
	target_url   string
	check_string string
	client       http.Client
}

func (l LinkChecker) checkLink() bool {
	resp, err := l.client.Get(l.target_url)
	if err != nil {
		return false
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), l.check_string)

}
