package graph

import (
	"sync"

	"github.com/translucent-link/owl/graph/model"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	EventObservers map[string]chan []model.AnyEvent
	mu             sync.Mutex
}
