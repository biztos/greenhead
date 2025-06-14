package rgxp_test

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/rgxp"
)

func TestParseBasicOK(t *testing.T) {

	require := require.New(t)

	src := "/foo/ism"
	r, err := rgxp.Parse(src)
	require.NoError(err, "parse")
	require.Equal(src, r.String(), "string")
	require.True(r.MatchString("FOO"), "matches")

}

func TestParseFailNoSlashes(t *testing.T) {

	require := require.New(t)

	src := "not.*delimited"
	_, err := rgxp.Parse(src)
	require.ErrorIs(err, rgxp.ErrNotRegexp)

}

func TestParseFailBadFlags(t *testing.T) {

	require := require.New(t)

	src := "/weird/flags"
	_, err := rgxp.Parse(src)
	require.ErrorIs(err, rgxp.ErrNotRegexp)
}

func TestParseFailDupeFlags(t *testing.T) {

	require := require.New(t)

	src := "/dupes/isi"
	_, err := rgxp.Parse(src)
	require.ErrorIs(err, rgxp.ErrDupeFlag)

}

func TestParseFailBadPattern(t *testing.T) {

	require := require.New(t)

	src := "/almost[.*/"
	_, err := rgxp.Parse(src)
	require.ErrorIs(err, rgxp.ErrInvalid)

}

func TestMustParseOK(t *testing.T) {

	require := require.New(t)

	src := "/foo/ism"
	r := rgxp.MustParse(src)
	require.Equal(src, r.String(), "string")
	require.True(r.MatchString("FOO"), "matches")

}

func TestMustParsePanics(t *testing.T) {

	require := require.New(t)

	require.PanicsWithError(rgxp.ErrNotRegexp.Error()+`: "x"`,
		func() { rgxp.MustParse("x") })

}

type RgxpConfig struct {
	Re *rgxp.Rgxp `toml:"some_regexp"`
}

func TestMarshalTomlOK(t *testing.T) {

	require := require.New(t)

	c := &RgxpConfig{rgxp.MustParse("/foo/ism")}
	b, err := toml.Marshal(c)
	require.NoError(err)
	require.Equal("some_regexp = \"/foo/ism\"\n", string(b))

}

func TestUnmarshalTomlOK(t *testing.T) {

	require := require.New(t)

	var c RgxpConfig
	err := toml.Unmarshal([]byte(`some_regexp = "/foo/is"`), &c)
	require.NoError(err)
	require.Equal("/foo/is", c.Re.String())

}

func TestUnmarshalTomlError(t *testing.T) {

	require := require.New(t)

	var c RgxpConfig
	err := toml.Unmarshal([]byte(`some_regexp = "nope not regexp"`), &c)
	require.ErrorContains(err, rgxp.ErrNotRegexp.Error())

}

func TestUnmarshalTomlRoundTrips(t *testing.T) {

	require := require.New(t)

	res := []string{
		"/foo/",
		"/foo/i",
		"/foo/bar/baz/ism",
		"/foo.*[\"]+bar/s",
		"/foo\n\nbar/",
	}
	for _, src := range res {
		c := &RgxpConfig{rgxp.MustParse(src)}
		b, err := toml.Marshal(c)
		require.NoError(err, "marshal")
		require.Equal(fmt.Sprintf("some_regexp = %q\n", src), string(b), "toml")
		var c2 RgxpConfig
		require.NoError(toml.Unmarshal(b, &c2), "unmarshal")
		require.Equal(src, c2.Re.String())
	}

}

func TestParseOptionalStringOK(t *testing.T) {

	require := require.New(t)

	r, err := rgxp.ParseOptional("foo")
	require.NoError(err)
	require.False(r.IsRegexp())
	require.Equal("foo", r.String())

}

func TestParseOptionalRegexpOK(t *testing.T) {

	require := require.New(t)

	r, err := rgxp.ParseOptional("/foo/")
	require.NoError(err)
	require.True(r.IsRegexp())
	require.Equal("/foo/", r.String())

}

func TestParseOptionalFails(t *testing.T) {

	require := require.New(t)

	_, err := rgxp.ParseOptional("/foo[/")
	require.ErrorIs(err, rgxp.ErrInvalid)

}

func TestMustParseOptionalOK(t *testing.T) {

	require := require.New(t)

	src := "/foo/ism"
	r := rgxp.MustParseOptional(src)
	require.Equal(src, r.String(), "string")
	require.True(r.IsRegexp(), "is")
	require.True(r.MatchString("FOO"), "matches")

}

func TestMustParseOptinalPanics(t *testing.T) {

	require := require.New(t)

	require.PanicsWithError(rgxp.ErrInvalid.Error()+`: "/foo[/"`,
		func() { rgxp.MustParseOptional("/foo[/") })

}

type OptConfig struct {
	Re *rgxp.OptionalRgxp `toml:"maybe_regexp"`
}

func TestMarshalTomlOptionalRegexpOK(t *testing.T) {

	require := require.New(t)

	c := &OptConfig{rgxp.MustParseOptional("/foo/ism")}
	b, err := toml.Marshal(c)
	require.NoError(err)
	require.Equal("maybe_regexp = \"/foo/ism\"\n", string(b))

}

func TestMarshalTomlOptionalStringOK(t *testing.T) {

	require := require.New(t)

	c := &OptConfig{rgxp.MustParseOptional("plain")}
	b, err := toml.Marshal(c)
	require.NoError(err)
	require.Equal("maybe_regexp = \"plain\"\n", string(b))

}

func TestUnmarshalTomlOptionalRegexpOK(t *testing.T) {

	require := require.New(t)

	var c OptConfig
	err := toml.Unmarshal([]byte(`maybe_regexp = "/foo/is"`), &c)
	require.NoError(err)
	require.True(c.Re.IsRegexp())
	require.Equal("/foo/is", c.Re.String())

}

func TestUnmarshalTomlOptionalStringOK(t *testing.T) {

	require := require.New(t)

	var c OptConfig
	err := toml.Unmarshal([]byte(`maybe_regexp = "plain"`), &c)
	require.NoError(err)
	require.False(c.Re.IsRegexp())
	require.Equal("plain", c.Re.String())

}

func TestUnmarshalTomlOptionalError(t *testing.T) {

	require := require.New(t)

	var c OptConfig
	err := toml.Unmarshal([]byte(`maybe_regexp = "/foo[/"`), &c)
	require.ErrorContains(err, rgxp.ErrInvalid.Error())

}

func TestMatchOrEqualString(t *testing.T) {

	require := require.New(t)

	type check struct {
		name string
		exp  bool
		src  string
	}
	checks := []check{
		{"regexp hits", true, "/f.*r/"},
		{"regexp misses", false, "/not/"},
		{"string hits", true, "foo bar"},
		{"string misses", false, "not"},
	}

	for _, c := range checks {
		r := rgxp.MustParseOptional(c.src)
		require.Equal(c.exp, r.MatchOrEqualString("foo bar"), c.name)
	}

}

// The case of "string-match a regexp" is in the docs so let's make sure it
// works.  If you really have to do this, you should probably not be here! :-)
func TestMatchOrEqualStringRegexpStyle(t *testing.T) {

	require := require.New(t)

	r := rgxp.MustParseOptional("/^[/]foo bar[/]$/")

	require.True(r.MatchOrEqualString("/foo bar/"))

}
