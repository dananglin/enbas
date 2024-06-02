// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

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

func ParseStatusContentType(value string) (StatusContentType, error) {
	switch value {
	case statusContentTypePlainValue, statusContentTypeTextPlainValue:
		return StatusContentTypePlainText, nil
	case statusContentTypeMarkdownValue, statusContentTypeTextMarkdownValue:
		return StatusContentTypeMarkdown, nil
	}

	return StatusContentTypeUnknown, InvalidStatusContentTypeError{Value: value}
}

func (s StatusContentType) MarshalJSON() ([]byte, error) {
	value := s.String()
	if value == unknownValue {
		return nil, InvalidStatusContentTypeError{Value: value}
	}

	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to encode %s to JSON: %w", value, err)
	}

	return data, nil
}

type InvalidStatusContentTypeError struct {
	Value string
}

func (e InvalidStatusContentTypeError) Error() string {
	return "'" + e.Value + "' is an invalid status content type: valid values are " +
		statusContentTypePlainValue + " or " + statusContentTypeTextPlainValue + " for plain text, or " +
		statusContentTypeMarkdownValue + " or " + statusContentTypeTextMarkdownValue + " for Markdown"
}
