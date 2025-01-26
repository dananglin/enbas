package model

type Thread struct {
	Context     Status
	Ancestors   StatusList
	Descendants StatusList
}
