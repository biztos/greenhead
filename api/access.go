package api

import (
	"errors"

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
	Roles []*Role `toml:"roles"`
	Keys  []*Key  `toml:"keys"`
}

var ErrKeyNotFound = errors.New("key not found")

// GetKey returns a Key for the provided AuthKey.
func (acc *Access) GetKey(key string) (*Key, error) {
	for _, k := range acc.Keys {
		if k.AuthKey == key {
			return k, nil
		}
	}
	return nil, ErrKeyNotFound
}

var ErrNoKeyRoles = errors.New("no roles found for key")

// KeyRoles returns all roles for the provided Key.
//
// If none of the key's named roles is present, ErrNoKeyRoles is returned.
func (acc *Access) KeyRoles(key *Key) ([]*Role, error) {

	// TODO: consider caching the names or a map; problem is then updates.
	// We *expect* this is never enough items to care about the length.
	roles := make([]*Role, 0, len(key.RoleNames))
	for _, name := range key.RoleNames {
		for _, role := range acc.Roles {
			if role.Name == name {
				roles = append(roles, role)
			}
		}
	}
	if len(roles) == 0 {
		return nil, ErrNoKeyRoles
	}
	return roles, nil
}

// KeyCanAccessURL checks whether any Role for the Key can access url.
func (acc *Access) KeyCanAccessURL(key *Key, url string) (bool, error) {
	roles, err := acc.KeyRoles(key)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.CanAccessURL(url) {
			return true, nil
		}
	}
	return false, nil
}

// KeyCanUseAgent checks whether any Role for the Key can use an agent with
// the given name.
func (acc *Access) KeyCanUseAgent(key *Key, name string) (bool, error) {
	roles, err := acc.KeyRoles(key)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role.CanUseAgent(name) {
			return true, nil
		}
	}
	return false, nil
}
