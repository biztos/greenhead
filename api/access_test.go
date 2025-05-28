package api_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/api"
	"github.com/biztos/greenhead/utils"
)

func TestAccessFromTomlOK(t *testing.T) {

	require := require.New(t)

	config := `[[roles]]
  name = "admin"
  description = "The Boss"
  endpoints = ["/.*/"]
  agents = ["/.*/"]

[[roles]]
  name = "support"
  description = "Poor You"
  endpoints = ["/^[/]support.*/", "/"]
  agents = ["/^support_/", "therapist"]

[[keys]]
  auth_key = "any-key-1"
  name = "The Bossman"
  roles = ["admin", "godhead"]

[[keys]]
  auth_key = "any-key-2"
  name = "The Worker"
  roles = ["support"]
`

	a := &api.Access{}
	err := toml.Unmarshal([]byte(config), a)
	require.NoError(err, "unmarshal")

	// Round-trip compare.
	back := utils.MustTomlString(a)
	require.Equal(config, back)

	// Access checks.
	_, err = a.GetKey("nope")
	require.ErrorIs(api.ErrKeyNotFound, err)
	bosskey, err := a.GetKey("any-key-1")
	require.NoError(err, "first key found")
	can_url, err := a.KeyCanAccessURL(bosskey, "/any/url")
	require.NoError(err, "boss KeyCanAccessURL")
	require.True(can_url, "boss KeyCanAccessURL")
	can_agent, err := a.KeyCanUseAgent(bosskey, "any_agent")
	require.NoError(err, "boss KeyCanUseAgent")
	require.True(can_agent, "boss KeyCanUseAgent")
	workerkey, err := a.GetKey("any-key-2")
	require.NoError(err, "second key found")
	worker_urls := map[string]bool{
		"/foo":                    false,
		"/support/anything/hello": true,
		"/":                       true,
	}
	for url, can := range worker_urls {
		can_url, err := a.KeyCanAccessURL(workerkey, url)
		require.NoError(err, "worker "+url)
		require.Equal(can, can_url, "worker "+url)
	}
	worker_agents := map[string]bool{
		"other":             false,
		"support_thaibox":   true,
		"therapist":         true,
		"therapist_thaibox": false,
	}
	for agent, can := range worker_agents {
		can_agent, err := a.KeyCanUseAgent(workerkey, agent)
		require.NoError(err, "worker "+agent)
		require.Equal(can, can_agent, "worker "+agent)
	}

}
