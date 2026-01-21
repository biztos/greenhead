package utils_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/ghd/runner"
	"github.com/biztos/greenhead/ghd/utils"
)

type CanNotMarshalJson struct {
	Bar func()
}

func TestMustJsonStringOK(t *testing.T) {

	require := require.New(t)
	v := map[string]any{"foo": "bar", "bar": []int{1, 2, 3}}
	var s string
	require.NotPanics(func() { s = utils.MustJsonString(v) }, "not panic")
	exp := `{"bar":[1,2,3],"foo":"bar"}` // uh-oh, any order guarantee?
	require.Equal(exp, s, "json as expected")
}

func TestMustJsonStringPanics(t *testing.T) {

	require := require.New(t)

	v := map[string]any{"foo": &CanNotMarshalJson{}}
	require.PanicsWithValue("json: unsupported type: func()",
		func() { utils.MustJsonString(v) }, "panic")

}

func TestMustJsonStringPrettyOK(t *testing.T) {

	require := require.New(t)
	v := map[string]any{"foo": "bar", "bar": []int{1, 2, 3}}
	var s string
	require.NotPanics(func() { s = utils.MustJsonStringPretty(v) }, "not panic")
	exp := `{
  "bar": [
    1,
    2,
    3
  ],
  "foo": "bar"
}`
	require.Equal(exp, s, "json as expected")
}

func TestMustJsonStringPrettyPanics(t *testing.T) {

	require := require.New(t)

	v := map[string]any{"foo": &CanNotMarshalJson{}}
	require.PanicsWithValue("json: unsupported type: func()",
		func() { utils.MustJsonStringPretty(v) }, "panic")

}

func TestDurOk(t *testing.T) {

	require := require.New(t)

	start := time.Now()
	time.Sleep(time.Nanosecond)
	s := utils.Dur(start)
	d, err := time.ParseDuration(s)
	require.NoError(err)

	// A bit dumb here, but let's just say it should not take a millisec!
	require.NotZero(d)
	require.True(d < time.Millisecond, "less than one millisecond")

}

func TestDurLogOk(t *testing.T) {

	require := require.New(t)

	start := time.Now()
	time.Sleep(time.Nanosecond)
	a := utils.DurLog(start)
	require.Equal(2, len(a), "has two elements")

	label, ok := a[0].(string)
	require.True(ok, "first is string")
	require.Equal(label, "duration")

	s, ok := a[1].(string)
	require.True(ok, "second is string")
	// As above...
	d, err := time.ParseDuration(s)
	require.NoError(err)
	require.NotZero(d)
	require.True(d < time.Millisecond, "less than one millisecond")

}

func TestUnmarshalFileErrBadExt(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	err := utils.UnmarshalFile("will-not-read.xlsx", cfg)
	require.EqualError(err, `unsupported extension: ".xlsx"`, "error")

}

func TestUnmarshalFileErrBadJson(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "runner_malformed.json")
	err := utils.UnmarshalFile(path, cfg)
	exp := "error parsing JSON: invalid character 'X' looking for beginning of value"
	require.EqualError(err, exp, "error")

}

func TestUnmarshalFileErrBadToml(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "runner_malformed.toml")
	err := utils.UnmarshalFile(path, cfg)
	exp := "error unmarshaling TOML: toml: line 2: expected '.' or '=', but got '?' instead"
	require.EqualError(err, exp, "error")

}

func TestUnmarshalFileErrJsonCanNotToml(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "can_not_toml.json")
	err := utils.UnmarshalFile(path, cfg)
	exp := "error marshaling TOML: toml: cannot encode array with nil element"
	require.EqualError(err, exp, "error")

}

func TestUnmarshalFileErrBadFile(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	err := utils.UnmarshalFile("no-such-file.toml", cfg)
	require.EqualError(err, `error reading file: open no-such-file.toml: no such file or directory`, "error")

}

func TestUnmarshalFileTomlOK(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "runner_simple.toml")
	err := utils.UnmarshalFile(path, cfg)
	require.NoError(err)
	require.True(cfg.Debug, "debug")
	require.True(cfg.Silent, "silent")
	require.True(cfg.Stream, "stream")
	require.Equal("file.log", cfg.LogFile, "log_file")
	require.Equal("dump.dir", cfg.DumpDir, "dump_dir")

}

func TestUnmarshalFileJsonOK(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "runner_simple.json")
	err := utils.UnmarshalFile(path, cfg)
	require.NoError(err)
	require.True(cfg.Debug, "debug")
	require.True(cfg.Silent, "silent")
	require.True(cfg.Stream, "stream")
	require.Equal("file.log", cfg.LogFile, "log_file")
	require.Equal("dump.dir", cfg.DumpDir, "dump_dir")

}

func TestUnmarshalFileJson5OK(t *testing.T) {

	require := require.New(t)

	cfg := &runner.Config{}
	path := filepath.Join("testdata", "runner_simple.json5")
	err := utils.UnmarshalFile(path, cfg)
	require.NoError(err)
	require.True(cfg.Debug, "debug")
	require.True(cfg.Silent, "silent")
	require.True(cfg.Stream, "stream")
	require.Equal("file.log", cfg.LogFile, "log_file")
	require.Equal("dump.dir", cfg.DumpDir, "dump_dir")

}
