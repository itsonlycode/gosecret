package secrets

import (
	"github.com/itsonlycode/gosecret/pkg/gosecret"
)

// New creates a new secret
func New() gosecret.Secret {
	return NewKV()
}
