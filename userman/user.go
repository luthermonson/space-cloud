package userman

import (
	"sync"

	"github.com/spaceuptech/space-cloud/auth"
	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/crud"
)

// Module is responsible for user management
type Module struct {
	sync.RWMutex
	methods map[string]struct{}
	crud    *crud.Module
	auth    *auth.Module
}

// Init creates a new instance of the user management object
func Init(crud *crud.Module, auth *auth.Module) *Module {
	return &Module{crud: crud, auth: auth}
}

// SetConfig set the config required by the user management module
func (m *Module) SetConfig(auth config.Auth) {
	m.Lock()
	defer m.Unlock()

	m.methods = make(map[string]struct{}, len(auth))

	for k := range auth {
		m.methods[k] = struct{}{}
	}
}

func (m *Module) isActive(method string) bool {
	m.RLock()
	defer m.RUnlock()

	_, p := m.methods[method]
	return p
}

func (m *Module) isEnabled() bool {
	m.RLock()
	defer m.RUnlock()

	return len(m.methods) > 0
}
