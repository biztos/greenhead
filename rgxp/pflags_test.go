package rgxp_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/rgxp"
)

func TestRgxpArrayString(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.RgxpArrayValue{&[]*rgxp.Rgxp{
		rgxp.MustParse("/foo/"),
		rgxp.MustParse("/bar/"),
	}}

	require.Equal(`["/foo/", "/bar/"]`, ary.String())

}

func TestRgxpArraySetOK(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.RgxpArrayValue{&[]*rgxp.Rgxp{
		rgxp.MustParse("/foo/"),
	}}
	require.NoError(ary.Set("/bar/")) // appends

	require.Equal(`["/foo/", "/bar/"]`, ary.String())

}

func TestRgxpArraySetError(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.RgxpArrayValue{&[]*rgxp.Rgxp{}}
	require.ErrorIs(ary.Set("plain"), rgxp.ErrNotRegexp)

}

func TestRgxpArrayType(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.RgxpArrayValue{&[]*rgxp.Rgxp{}}
	require.Equal("rgxpArray", ary.Type())

}

func TestOptionalRgxpArrayString(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.OptionalRgxpArrayValue{&[]*rgxp.OptionalRgxp{
		rgxp.MustParseOptional("/foo/"),
		rgxp.MustParseOptional("bar"),
	}}

	require.Equal(`["/foo/", "bar"]`, ary.String())

}

func TestOptionalRgxpArraySetOK(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.OptionalRgxpArrayValue{&[]*rgxp.OptionalRgxp{
		rgxp.MustParseOptional("/foo/"),
	}}
	require.NoError(ary.Set("bar")) // appends

	require.Equal(`["/foo/", "bar"]`, ary.String())

}

func TestOptionalRgxpArraySetError(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.OptionalRgxpArrayValue{&[]*rgxp.OptionalRgxp{}}
	require.ErrorIs(ary.Set("/foo[/"), rgxp.ErrInvalid)

}

func TestOptionalRgxpArrayType(t *testing.T) {

	require := require.New(t)

	ary := &rgxp.OptionalRgxpArrayValue{&[]*rgxp.OptionalRgxp{}}
	require.Equal("optRgxpArray", ary.Type())

}
