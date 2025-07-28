// blade.go
package blade

import (
	"regexp"
)

// The compileContent function orchestrates the transpilation from Blade to Go template syntax.
func compileContent(value string) string {
	value = compileComments(value)
	value = compileEchos(value)
	value = compileIf(value)
	// Additional lexers for loops, includes, etc., will be added here.
	return value
}

// compileComments removes Blade comments from the template.
func compileComments(value string) string {
	reg := regexp.MustCompile(`{{--(.+?)--}}`)
	return reg.ReplaceAllString(value, "")
}

// compileEchos translates Blade's echo syntax to Go's.
func compileEchos(value string) string {
	// This pattern now correctly captures the complex variables (for example, `$Page.title``) and replaces the `$`` with the Go template dot notation.
	reg := regexp.MustCompile(`{{\s*\$(.+?)\s*}}`)
	value = reg.ReplaceAllString(value, `{{ .${1} }}`)

	// This handles unescaped data, for example, `{!! $Content !!}`
	reg = regexp.MustCompile(`{!!\s*\$(.+?)\s*!!}`)
	value = reg.ReplaceAllString(value, `{{ .${1} }}`)

	return value
}

// compileIf translates Blade's conditional statements.
func compileIf(value string) string {
	// This pattern now correctly captures the complex variables within the directives.
	reg := regexp.MustCompile(`@if\s*\(\s*\$(.+?)\s*\)`)
	value = reg.ReplaceAllString(value, `{{if .${1}}}`)

	reg = regexp.MustCompile(`@elseif\s*\((.+?)\)`)
	value = reg.ReplaceAllString(value, `{{else if ${1}}}`)

	reg = regexp.MustCompile(`@else`)
	value = reg.ReplaceAllString(value, `{{else}}`)

	reg = regexp.MustCompile(`@endif`)
	value = reg.ReplaceAllString(value, `{{end}}`)

	return value
}
