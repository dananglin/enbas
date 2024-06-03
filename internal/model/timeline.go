// SPDX-FileCopyrightText: 2024 Dan Anglin <d.n.i.anglin@gmail.com>
//
// SPDX-License-Identifier: GPL-3.0-or-later

package model

const (
	TimelineCategoryHome   = "home"
	TimelineCategoryPublic = "public"
	TimelineCategoryTag    = "tag"
	TimelineCategoryList   = "list"
)

type InvalidTimelineCategoryError struct {
	Value string
}

func (e InvalidTimelineCategoryError) Error() string {
	return "'" +
		e.Value +
		"' is not a valid timeline category (valid values are " +
		TimelineCategoryHome + ", " +
		TimelineCategoryPublic + ", " +
		TimelineCategoryTag + ", " +
		TimelineCategoryList + ")"
}
