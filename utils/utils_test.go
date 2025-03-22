package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/runner"
	"github.com/biztos/greenhead/utils"
)

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
	exp := "error unmarshaling TOML: toml: line 2: expected '.' or '=', but got '?' instead"
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
