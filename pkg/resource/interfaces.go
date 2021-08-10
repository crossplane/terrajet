package resource

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// JSONIterable can get or set parameters or observations of terraform managed resources
type JSONIterable interface {
	GetObservation() ([]byte, error)
	SetObservation(data []byte) error

	GetParameters() ([]byte, error)
	SetParameters(data []byte) error
}

// Terraformed is a Kubernetes object representing a concrete terraform managed resource
type Terraformed interface {
	resource.Managed

	JSONIterable
}
