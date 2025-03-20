// Note that these tests use the registry package. ONLY do this when you MUST
// have access to internal stuff, which USUALLY is never the case.
//
// Here we need it in order to reset the locks, for which there is no API.
package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/tools"
)

// minimal dummy tool, we don't do anything but register it.
func testTool(name string) tools.Tooler {
	return tools.NewTool[int, int](name, "",
		func(c context.Context, i int) (int, error) { return 0, nil })
}

func resetLocksAndClear() {
	lockedForNew = false
	lockedForReplace = false
	Clear()
}

func TestLock(t *testing.T) {

	require := require.New(t)

	resetLocksAndClear()
	defer resetLocksAndClear()

	// Register one:
	require.NoError(Register(testTool("foo")), "new ok")

	// Register dupe:
	require.NoError(Register(testTool("foo")), "replace ok")

	Lock()
	require.EqualError(Register(testTool("bar")),
		"registry is locked for new tools", "new blocked")
	require.EqualError(Register(testTool("foo")),
		"registry is locked for replacement tools", "replace blocked")

	// "Other" lock not affected by specific locks.
	LockForReplace()
	require.EqualError(Register(testTool("foo")),
		"registry is locked for replacement tools", "replace blocked")
	LockForNew()
	require.EqualError(Register(testTool("bar")),
		"registry is locked for new tools", "new blocked")

}

func TestLockForNew(t *testing.T) {

	require := require.New(t)

	resetLocksAndClear()
	defer resetLocksAndClear()

	// Register one:
	require.NoError(Register(testTool("foo")), "new ok")

	// Register dupe:
	require.NoError(Register(testTool("foo")), "replace ok")

	// Deny new:
	LockForNew()
	require.EqualError(Register(testTool("bar")),
		"registry is locked for new tools", "replace blocked")

	// Still can do dupe:
	require.NoError(Register(testTool("foo")), "replace ok")

}

func TestLockForReplace(t *testing.T) {

	require := require.New(t)

	resetLocksAndClear()
	defer resetLocksAndClear()

	// Register one:
	require.NoError(Register(testTool("foo")), "new ok")

	// Register dupe:
	require.NoError(Register(testTool("foo")), "replace ok")

	// Deny dupe:
	LockForReplace()
	require.EqualError(Register(testTool("foo")),
		"registry is locked for replacement tools", "replace blocked")

	// Still can do new:
	require.NoError(Register(testTool("bar")), "new ok")

}
