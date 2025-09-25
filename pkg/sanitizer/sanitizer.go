package sanitizer

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	policy *bluemonday.Policy
}

func NewSanitizer() *Sanitizer {
	policy := bluemonday.StrictPolicy()
	return &Sanitizer{policy: policy}
}

func (s *Sanitizer) SanitizeSQL(input string) string {
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")
	input = strings.ReplaceAll(input, "/*", "")
	input = strings.ReplaceAll(input, "*/", "")

	return s.SanitizeXSS(input)
}

func (s *Sanitizer) SanitizeXSS(input string) string {
	return s.policy.Sanitize(input)
}
