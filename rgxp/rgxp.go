// Package rgxp provides convenience wrappers for regexp.Regexp.
package rgxp

import (
	"errors"
	"fmt"
	"regexp"
)

var ErrNotRegexp = errors.New("not a regular expression in /pattern/[ism] format")

var ErrDupeFlag = errors.New("duplicate flag in regular expression")

var ErrInvalid = errors.New("invalid regular expression")

// Rgxp wraps regexp.Regexp
type Rgxp struct {
	regexp.Regexp
	src string
}

var matchRegexp = regexp.MustCompile("(?s)^/(.*)/([smi]{0,3})$")

// Parse parses src, which must be a regular expression of the format:
//
//	"/pattern/[ism]"
//
// The pattern part must be a valid Go regular expression.
func Parse(src string) (*Rgxp, error) {
	m := matchRegexp.FindStringSubmatch(src)
	if m == nil {
		return nil, fmt.Errorf("%w: %q", ErrNotRegexp, src)
	}
	gore := m[1]
	flags := m[2]
	have_flag := map[rune]bool{}
	for _, f := range flags {
		if have_flag[f] {
			return nil, fmt.Errorf("%w: %q in %q", ErrDupeFlag, f, src)
		}
		have_flag[f] = true
	}
	if flags != "" {
		gore = "(?" + flags + ")" + gore
	}
	re, err := regexp.Compile(gore)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrInvalid, src)
	}

	return &Rgxp{*re, src}, nil
}

// MustParse calls Parse and panics on error.
func MustParse(src string) *Rgxp {
	r, err := Parse(src)
	if err != nil {
		panic(err)
	}
	return r
}

// String returns the source from which r was built.
func (r *Rgxp) String() string {
	return r.src
}

// MarshalTOML returns the text from which r was built.
func (r *Rgxp) MarshalTOML() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", r.src)), nil
}

// UnmarshalText parses text as the source regular expression.
func (r *Rgxp) UnmarshalText(text []byte) error {
	r2, err := Parse(string(text))
	if err != nil {
		return err
	}
	r.src = r2.src
	r.Regexp = r2.Regexp
	return nil
}

// OptionalRgxp wraps Rgxp for values that might or might not be a Regexp.
// (It might just be a plain string.)  Anything that matches the standard
// format is considered a regular expression:
//
//	"/pattern/[ism]"
//
// Attempting to match on an OptionalRgxp for which IsRegexp returns false is
// likely to panic.
//
// Always check IsRegexp before assuming you have a valid Regexp!
type OptionalRgxp struct {
	Rgxp
	isRegexp bool
}

// ParseOptional parses src into an OptionalRegexp.
func ParseOptional(src string) (*OptionalRgxp, error) {
	r, err := Parse(src)
	if errors.Is(err, ErrNotRegexp) {
		// Fine, we'll take the string then!
		return &OptionalRgxp{Rgxp{regexp.Regexp{}, src}, false}, nil
	} else if err != nil {
		return nil, err
	}
	return &OptionalRgxp{*r, true}, nil
}

// MustParseOptional calls ParseOptional and panics on error.
func MustParseOptional(src string) *OptionalRgxp {
	r, err := ParseOptional(src)
	if err != nil {
		panic(err)
	}
	return r
}

// IsRegexp tells whether r holds a usable Regexp or just plain text.
func (r *OptionalRgxp) IsRegexp() bool {
	return r.isRegexp
}

// MarshalTOML returns the text from which r was built.
func (r *OptionalRgxp) MarshalTOML() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", r.src)), nil
}

// UnmarshalText parses text as the source regular expression or plain string.
func (r *OptionalRgxp) UnmarshalText(text []byte) error {
	r2, err := ParseOptional(string(text))
	if err != nil {
		return err
	}
	r.Rgxp = r2.Rgxp
	r.isRegexp = r2.isRegexp
	return nil
}
