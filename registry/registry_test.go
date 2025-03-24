package registry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/registry"
	"github.com/biztos/greenhead/tools"
)

type TestInput struct {
	Val string
}

func testTool(name string) tools.Tooler {
	return tools.NewTool[TestInput, string](name, name+" ok\nyes!",
		func(ctx context.Context, in TestInput) (string, error) {
			return name + "_ok", nil
		})
}

func TestRegisterBlankNameFails(t *testing.T) {

	require := require.New(t)

	tool := testTool("   ")

	// Try to keep it clean...
	registry.Clear()
	defer registry.Clear()

	require.EqualError(registry.Register(tool), "empty name for tool")

}

func TestRegisterNewOK(t *testing.T) {

	require := require.New(t)

	tool := testTool("foo")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool), "register ok")

	// Did we "get" it?
	got, err := registry.Get("foo")
	require.NoError(err, "got ok")
	require.Equal(tool, got, "got the tool")

}

func TestRegisterReplaceOK(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("foo")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new ok")
	require.NoError(registry.Register(tool2), "replace ok")

	got, err := registry.Get("foo")
	require.NoError(err, "got ok")
	require.Equal(tool2, got, "got the replacement tool")

}

func TestRemoveOK(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("bar")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new foo ok")
	require.NoError(registry.Register(tool2), "new bar ok")
	require.NoError(registry.Remove("bar"), "remove bar ok")

	got, err := registry.Get("foo")
	require.NoError(err, "got ok")
	require.Equal(tool1, got, "got the remaining tool")
	_, err = registry.Get("bar")
	require.ErrorIs(err, registry.ErrNotRegistered, "get removed")

}

func TestRemoveErrNotRegistered(t *testing.T) {

	require := require.New(t)

	registry.Clear()
	defer registry.Clear()

	err := registry.Remove("nonesuch")
	require.ErrorIs(err, registry.ErrNotRegistered, "remove nonexistent")

}

func TestNamesOrdered(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("zoo")
	tool3 := testTool("foo")
	tool4 := testTool("aaa")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new foo ok")
	require.NoError(registry.Register(tool2), "new zoo ok")
	require.NoError(registry.Register(tool3), "replace foo ok")
	require.NoError(registry.Register(tool4), "new aaa ok")
	require.NoError(registry.Remove("aaa"), "remove aaa ok")

	require.Equal([]string{"zoo", "foo"}, registry.Names(),
		"names in registration order")

}

func TestDisplayLines(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("zoo")
	tool3 := testTool("gooboo")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new ok")
	require.NoError(registry.Register(tool2), "new ok")
	require.NoError(registry.Register(tool3), "new ok")

	exp := []string{
		"foo    - foo ok",
		"zoo    - zoo ok",
		"gooboo - gooboo ok",
	}
	require.Equal(exp, registry.DisplayLines(),
		"display lines formatted nice-like")

}
