package bootstrap

import (
	"gitlab.com/kode4food/ale/api"
	"gitlab.com/kode4food/ale/internal/namespace"
)

type (
	bootstrap struct {
		manager *namespace.Manager
		funcMap funcMap
	}

	funcMap map[api.Name]*api.Function
)

// Into sets up initial built-ins and assets
func Into(manager *namespace.Manager) {
	b := &bootstrap{
		manager: manager,
		funcMap: funcMap{},
	}
	b.builtIns()
	b.assets()
}
