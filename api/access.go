package api

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"

	"github.com/biztos/greenhead/rgxp"
)

// EncodeAuthKey takes an AuthKey and returns a client-facing value.
//
// The client-facing value should be sent by the client in an Authorization
// header, as:
//
//	Authorization: Bearer <encoded-key>
//
// This is used in the KeyAccess middleware.
//
// The standard implementation is an RFC 4648 "URL" Base64 encoding of the
// SHA256 hash.
var EncodeAuthKey = func(s string) string {
	h := sha256.Sum256([]byte(s))
	return base64.RawURLEncoding.EncodeToString(h[:])

}

// NotEncodeAuthKey is the non-encoder, used when RawKeys is configured.
var NotEncodeAuthKey = func(s string) string {
	return s
}

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
	roles      []*Role
	keys       []*Key
	keyMap     map[string]*Key
	keyRoles   map[*Key][]*Role
	keyEncoder func(string) string
}

var ErrBlankRoleName = errors.New("empty or blank role name")
var ErrBlankAuthKey = errors.New("empty or blank auth key")
var ErrBlankKeyName = errors.New("empty or blank key name")
var ErrDupeRoleName = errors.New("duplicate role")
var ErrDupeAuthKey = errors.New("duplicate auth key")
var ErrDupeKeyName = errors.New("duplicate key name")

// DefaultAccess creates an Access with one key that has full access to
// all endpoints and agents.
//
// The AuthKey is returned together with the Access.
func DefaultAccess(encoder func(string) string) (*Access, string) {

	allow_all, _ := rgxp.ParseOptional("/.*/")
	role := &Role{
		Name:      "default-all-access-role",
		Endpoints: []*rgxp.OptionalRgxp{allow_all},
		Agents:    []*rgxp.OptionalRgxp{allow_all},
	}
	key := &Key{
		AuthKey:   ulid.Make().String(),
		Name:      "default-all-access-key",
		RoleNames: []string{"default-all-access-role"},
	}
	acc, _ := NewAccess([]*Role{role}, []*Key{key}, encoder)
	return acc, key.AuthKey
}

// NewAccess creates an Access from roles and keys.
//
// Duplicates by Name or AuthKey are disallowed, as are empty/all-whitespace
// strings for both.  Name your roles and keys.
func NewAccess(roles []*Role, keys []*Key, encoder func(string) string) (*Access, error) {

	if encoder == nil {
		encoder = EncodeAuthKey
	}
	role_map := map[string]*Role{}
	for _, r := range roles {
		if strings.TrimSpace(r.Name) == "" {
			return nil, ErrBlankRoleName
		}
		if role_map[r.Name] != nil {
			return nil, fmt.Errorf("%w: %q", ErrDupeRoleName, r.Name)
		}
		role_map[r.Name] = r
	}
	have_key_name := map[string]bool{}
	key_map := map[string]*Key{}
	for _, k := range keys {
		if strings.TrimSpace(k.AuthKey) == "" {
			return nil, ErrBlankAuthKey
		}
		if strings.TrimSpace(k.Name) == "" {
			return nil, ErrBlankKeyName
		}
		if have_key_name[k.Name] {
			return nil, ErrDupeKeyName
		}
		have_key_name[k.Name] = true
		encoded := encoder(k.AuthKey)
		if key_map[encoded] != nil {
			// NB: don't put the key in the error message!
			return nil, fmt.Errorf("%w for %s", ErrDupeAuthKey, k.Name)
		}
		key_map[encoded] = k

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
	return &Access{
		roles:      roles,
		keys:       keys,
		keyMap:     key_map,
		keyRoles:   key_roles,
		keyEncoder: encoder,
	}, nil

}

// GetKey returns a Key for the provided AuthKey k, which is assumed to be
// encoded with the equivalent of the encoder function passed to NewAccess.
//
// If no key is found, nil is returned.
func (acc *Access) GetKey(k string) *Key {
	return acc.keyMap[k]
}

// EndpointAllowed checks whether any Role for the Key can access endpoint.
func (acc *Access) EndpointAllowed(key *Key, endpoint string) bool {
	for _, role := range acc.keyRoles[key] {
		if role.CanAccessURL(endpoint) {
			return true
		}
	}
	return false
}

// AgentAllowed checks whether any Role for the Key can use an agent with
// the given name.
func (acc *Access) AgentAllowed(key *Key, name string) bool {
	for _, role := range acc.keyRoles[key] {
		if role.CanUseAgent(name) {
			return true
		}
	}
	return false
}
