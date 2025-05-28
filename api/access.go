package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/rgxp"
)

// Role defines a set of permissions for API Keys.
type Role struct {
	Name        string               `toml:"name"`        // Name of role.
	Description string               `toml:"description"` // Description of role.
	Endpoints   []*rgxp.OptionalRgxp `toml:"endpoints"`   // Endpoint access.
	Agents      []*rgxp.OptionalRgxp `toml:"agents"`      // Agents access.
}

// CanAccessURL checks that the Role can access url.
func (r *Role) CanAccessURL(url string) bool {
	for _, e := range r.Endpoints {
		if e.MatchOrEqualString(url) {
			return true
		}
	}
	return false
}

// CanUseAgent checks that Role can use an agent with the given name.
func (r *Role) CanUseAgent(name string) bool {
	for _, a := range r.Agents {
		if a.MatchOrEqualString(name) {
			return true
		}
	}
	return false
}

// Key defines an API Key that is attached to a Role by name.
type Key struct {
	AuthKey   string   `toml:"auth_key"` // Key string for client auth.
	Name      string   `toml:"name"`     // Name of the key user for logs/UI.
	RoleNames []string `toml:"roles"`    // Name of the role of this key.
}

// Access manages the Roles and Keys used to access the system.
type Access struct {
	roles    []*Role
	keys     []*Key
	keyMap   map[string]*Key
	keyRoles map[*Key][]*Role
}

var ErrBlankRole = errors.New("blank role name")
var ErrBlankKey = errors.New("blank auth key")
var ErrDupeRole = errors.New("duplicate role")
var ErrDupeKey = errors.New("duplicate key")

// DefaultAccess creates an Access with one key that has full access to
// all endpoints and agents.
//
// The AuthKey is returned together with the Access.
func DefaultAccess() (*Access, string) {

	allow_all, _ := rgxp.ParseOptional("/.*/")
	role := &Role{
		Name:      "default-all-access-role",
		Endpoints: []*rgxp.OptionalRgxp{allow_all},
		Agents:    []*rgxp.OptionalRgxp{allow_all},
	}
	key := &Key{
		AuthKey:   ulid.Make().String(),
		Name:      "default-all-access-user",
		RoleNames: []string{"default-all-access-role"},
	}
	acc, _ := NewAccess([]*Role{role}, []*Key{key})
	return acc, key.AuthKey
}

// NewAccess creates an Access from roles and keys.
//
// Duplicates by Name or AuthKey are disallowed, as are blank strings for
// both.
func NewAccess(roles []*Role, keys []*Key) (*Access, error) {

	role_map := map[string]*Role{}
	for _, r := range roles {
		if strings.TrimSpace(r.Name) == "" {
			return nil, ErrBlankRole
		}
		if role_map[r.Name] != nil {
			return nil, fmt.Errorf("%w: %q", ErrDupeRole, r.Name)
		}
		role_map[r.Name] = r
	}
	key_map := map[string]*Key{}
	for _, k := range keys {
		if strings.TrimSpace(k.AuthKey) == "" {
			return nil, ErrBlankKey
		}
		if key_map[k.AuthKey] != nil {
			// NB: don't put the key in the error message!
			return nil, fmt.Errorf("%w for %s", ErrDupeKey, k.Name)
		}
		key_map[k.AuthKey] = k
	}

	// map roles to keys so we don't have to do it at every check.
	key_roles := map[*Key][]*Role{}
	for _, k := range keys {
		roles := []*Role{}
		for _, name := range k.RoleNames {
			if role_map[name] != nil {
				roles = append(roles, role_map[name])
			}
		}
		key_roles[k] = roles
	}
	return &Access{roles, keys, key_map, key_roles}, nil

}

var ErrKeyNotFound = errors.New("key not found")

// GetKey returns a Key for the provided AuthKey.
func (acc *Access) GetKey(key string) (*Key, error) {
	if acc.keyMap[key] != nil {
		return acc.keyMap[key], nil
	}
	return nil, ErrKeyNotFound
}

// EndpointAllowed checks whether any Role for the Key can access endpoint.
func (acc *Access) EndpointAllowed(key *Key, endpoint string) (bool, error) {
	for _, role := range acc.keyRoles[key] {
		if role.CanAccessURL(endpoint) {
			return true, nil
		}
	}
	return false, nil
}

// AgentAllowed checks whether any Role for the Key can use an agent with
// the given name.
func (acc *Access) AgentAllowed(key *Key, name string) (bool, error) {
	for _, role := range acc.keyRoles[key] {
		if role.CanUseAgent(name) {
			return true, nil
		}
	}
	return false, nil
}
