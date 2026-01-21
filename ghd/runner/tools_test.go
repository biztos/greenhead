package runner_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/ghd/registry"
	"github.com/biztos/greenhead/ghd/runner"
	"github.com/biztos/greenhead/ghd/tools"
)

type TestInput struct {
	Val string `json:"val"`
}

func testTool(name string) tools.Tooler {
	return tools.NewTool[TestInput, string](name, name+" ok\nyes!",
		func(ctx context.Context, in TestInput) (string, error) {
			return name + " " + in.Val, nil
		})
}

func blankRunner() *runner.Runner {
	r, err := runner.NewRunner(&runner.Config{})
	if err != nil {
		panic(err) // sort of the most basic thing you could do...
	}
	return r
}

func TestListToolsNoTools(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	r, err := runner.NewRunner(&runner.Config{})
	require.NoError(err, "NewRunner errs not")

	buf := new(bytes.Buffer)
	require.NoError(r.ListTools(buf, true), "ListTools errs not")

	exp := "<no tools>\n"
	require.Equal(exp, buf.String(), "list output")

}

func TestListToolsNameOnly(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(testTool("foo")), "reg foo")
	require.NoError(registry.Register(testTool("bar")), "reg bar")

	r, err := runner.NewRunner(&runner.Config{})
	require.NoError(err, "NewRunner errs not")

	buf := new(bytes.Buffer)
	require.NoError(r.ListTools(buf, true), "ListTools errs not")

	exp := "foo\nbar\n"
	require.Equal(exp, buf.String(), "list output")

}

func TestListToolsLong(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(testTool("foo")), "reg foo")
	require.NoError(registry.Register(testTool("barzoo")), "reg bar")

	r, err := runner.NewRunner(&runner.Config{})
	require.NoError(err, "NewRunner errs not")

	buf := new(bytes.Buffer)
	require.NoError(r.ListTools(buf, false), "ListTools errs not")

	exp := "foo    - foo ok\nbarzoo - barzoo ok\n"
	require.Equal(exp, buf.String(), "list output")

}

func TestShowToolError(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	buf := new(bytes.Buffer)
	err := blankRunner().ShowTool(buf, "foo")
	require.EqualError(err, `tool is not registered: "foo"`)
	require.Equal("", buf.String(), "no output")

}

func TestShowToolOK(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(testTool("foo")), "reg foo")

	buf := new(bytes.Buffer)
	err := blankRunner().ShowTool(buf, "foo")
	require.NoError(err, "show tool")
	exp := `foo

foo ok
yes!

Input Schema:

{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "val": {
      "type": "string"
    }
  },
  "required": [
    "val"
  ]
}

Return Type: string, error

`

	require.Equal(exp, buf.String(), "show output")

}

func TestRunToolErrorNoTool(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	out, err := runner.RunTool("foo", "[")
	require.EqualError(err, `tool is not registered: "foo"`)
	require.Nil(out, "output")

}

func TestRunToolErrorArgs(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(testTool("foo")), "reg foo")

	out, err := runner.RunTool("foo", "[")
	require.Error(err, "got error")
	require.Regexp("error parsing json", err, "error as expected")
	require.Nil(out, "output")

}

func TestRunToolOK(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(testTool("foo")), "reg foo")

	out, err := runner.RunTool("foo", `{"val":"boo"}`)
	require.NoError(err)
	require.Equal(out, "foo boo", "output")

}
