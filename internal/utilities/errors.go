package utilities

type UnspecifiedProgramError struct{}

func (e UnspecifiedProgramError) Error() string {
	return "the program to view these files is unspecified"
}

type UnspecifiedBrowserError struct{}

func (e UnspecifiedBrowserError) Error() string {
	return "the browser to view this link is not specified"
}
