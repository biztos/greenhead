package api_test

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/biztos/greenhead/api"
)

func TestAccessFromTomlOK(t *testing.T) {

	require := require.New(t)

	config_toml := `[[roles]]
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

	config := &api.Config{}
	err := toml.Unmarshal([]byte(config_toml), config)
	require.NoError(err, "unmarshal")

	// We can now construct based on our roles and keys.
	a, err := api.NewAccess(config.Roles, config.Keys, nil)
	require.NoError(err, "NewAccess")

	// Access checks.
	nokey := a.GetKey("nope")
	require.Nil(nokey)
	bosskey := a.GetKey(api.EncodeAuthKey("any-key-1"))
	require.NotNil(bosskey)
	can_url := a.EndpointAllowed(bosskey, "/any/url")
	require.True(can_url, "boss EndpointAllowed")
	can_agent := a.AgentAllowed(bosskey, "any_agent")
	require.True(can_agent, "boss AgentAllowed")
	workerkey := a.GetKey(api.EncodeAuthKey("any-key-2"))
	require.NotNil(workerkey)
	worker_urls := map[string]bool{
		"/foo":                    false,
		"/support/anything/hello": true,
		"/":                       true,
	}
	for url, can := range worker_urls {
		can_url := a.EndpointAllowed(workerkey, url)
		require.Equal(can, can_url, "worker "+url)
	}
	worker_agents := map[string]bool{
		"other":             false,
		"support_thaibox":   true,
		"therapist":         true,
		"therapist_thaibox": false,
	}
	for agent, can := range worker_agents {
		can_agent := a.AgentAllowed(workerkey, agent)
		require.Equal(can, can_agent, "worker "+agent)
	}

}
