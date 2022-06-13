package dockerregistry

import "regexp"

var (
	alphaNumeric = `[a-z0-9]+`
	separator = `(?:[._]|__|[-]*)`

	nameComponent = expression(
		alphaNumeric,
		optional(repeated(separator, alphaNumeric)))

	tag = `[\w][\w.-]{0,127}`
	TagRegexp = regexp.MustCompile(tag)
)


func expression(res ...string) string {
	var s string
	for _, re := range res {
		s += re
	}

	return s
}


// optional wraps the expression in a non-capturing group and makes the
// production optional.
func optional(res ...string) string {
	return group(expression(res...)) + `?`
}

// repeated wraps the regexp in a non-capturing group to get one or more
// matches.
func repeated(res ...string) string {
	return group(expression(res...)) + `+`
}

// group wraps the regexp in a non-capturing group.
func group(res ...string) string {
	return `(?:` + expression(res...) + `)`
}

// capture wraps the expression in a capturing group.
func capture(res ...string) string {
	return `(` + expression(res...) + `)`
}

// anchored anchors the regular expression by adding start and end delimiters.
func anchored(res ...string) string {
	return `^` + expression(res...) + `$`
}