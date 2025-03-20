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

type TestInputCannotMarshal struct {
	Nope *testing.T
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
	require.Equal(tool, registry.Get("foo"), "got the tool")

}

func TestRegisterReplaceOK(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("foo")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new ok")
	require.NoError(registry.Register(tool2), "replace ok")

	require.Equal(tool2, registry.Get("foo"), "got the replacement tool")

}

func TestNamesOrdered(t *testing.T) {

	require := require.New(t)

	tool1 := testTool("foo")
	tool2 := testTool("zoo")
	tool3 := testTool("foo")

	registry.Clear()
	defer registry.Clear()

	require.NoError(registry.Register(tool1), "new ok")
	require.NoError(registry.Register(tool2), "new ok")
	require.NoError(registry.Register(tool3), "replace ok")

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
