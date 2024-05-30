package model

import (
	"encoding/json"
	"fmt"
)

type StatusContentType int

const (
	StatusContentTypePlainText StatusContentType = iota
	StatusContentTypeMarkdown
	StatusContentTypeUnknown
)

const (
	statusContentTypeTextPlainValue    = "text/plain"
	statusContentTypePlainValue        = "plain"
	statusContentTypeTextMarkdownValue = "text/markdown"
	statusContentTypeMarkdownValue     = "markdown"
)

func (s StatusContentType) String() string {
	mapped := map[StatusContentType]string{
		StatusContentTypeMarkdown:  statusContentTypeTextMarkdownValue,
		StatusContentTypePlainText: statusContentTypeTextPlainValue,
	}

	output, ok := mapped[s]
	if !ok {
		return unknownValue
	}

	return output
}

func ParseStatusContentType(value string) StatusContentType {
	switch value {
	case statusContentTypePlainValue, statusContentTypeTextPlainValue:
		return StatusContentTypePlainText
	case statusContentTypeMarkdownValue, statusContentTypeTextMarkdownValue:
		return StatusContentTypeMarkdown
	}

	return StatusContentTypeUnknown
}

func (s StatusContentType) MarshalJSON() ([]byte, error) {
	value := s.String()
	if value == unknownValue {
		return nil, fmt.Errorf("%q is not a valid status content type", value)
	}

	return json.Marshal(value)
}
