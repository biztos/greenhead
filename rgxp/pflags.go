package rgxp

import (
	"fmt"
	"strings"
)

// RgxpArrayValue facilitates the use of Regexps in Cobra pflags.
type RgxpArrayValue struct {
	Rgxps *[]*Rgxp
}

// String implements pflag.Value.
func (r *RgxpArrayValue) String() string {
	var patterns []string
	for _, re := range *r.Rgxps {
		patterns = append(patterns, fmt.Sprintf("%q", re.String()))
	}
	return fmt.Sprintf("[%s]", strings.Join(patterns, ", "))
}

// Set implements pflag.Value.
func (r *RgxpArrayValue) Set(s string) error {
	re, err := Parse(s)
	if err != nil {
		return err
	}
	*r.Rgxps = append(*r.Rgxps, re)
	return nil
}

// Type implements pflag.Value.
func (r *RgxpArrayValue) Type() string {
	return "rgxpArray"
}

// OptionalRgxpArrayValue facilitates the use of Regexps in Cobra pflags.
type OptionalRgxpArrayValue struct {
	OptionalRgxps *[]*OptionalRgxp
}

// String implements pflag.Value.
func (r *OptionalRgxpArrayValue) String() string {
	var patterns []string
	for _, re := range *r.OptionalRgxps {
		patterns = append(patterns, fmt.Sprintf("%q", re.String()))
	}
	return fmt.Sprintf("[%s]", strings.Join(patterns, ", "))
}

// Set implements pflag.Value.
func (r *OptionalRgxpArrayValue) Set(s string) error {
	re, err := ParseOptional(s)
	if err != nil {
		return err
	}
	*r.OptionalRgxps = append(*r.OptionalRgxps, re)
	return nil
}

// Type implements pflag.Value.
func (r *OptionalRgxpArrayValue) Type() string {
	return "optRgxpArray"
}
