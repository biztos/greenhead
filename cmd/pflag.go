package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexpArrayValue facilitates the use of Regexps in Cobra pflags.
type RegexpArrayValue struct {
	regexps *[]*regexp.Regexp
}

// String implements pflag.Value.
func (r *RegexpArrayValue) String() string {
	var patterns []string
	for _, re := range *r.regexps {
		patterns = append(patterns, re.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(patterns, ", "))
}

// Set implements pflag.Value.
func (r *RegexpArrayValue) Set(s string) error {
	re, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	*r.regexps = append(*r.regexps, re)
	return nil
}

// Type implements pflag.Value.
func (r *RegexpArrayValue) Type() string {
	return "regexpArray"
}
