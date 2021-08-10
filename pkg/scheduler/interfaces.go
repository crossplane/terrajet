package scheduler

import (
	"context"

	"github.com/crossplane-contrib/terrajet/pkg/resource"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
)

// ApplyResult represents result of an apply operation
type ApplyResult struct {
	// Tells whether the apply operation is completed.
	Completed bool

	// Sensitive information that is available during creation/update.
	ConnectionDetails managed.ConnectionDetails
}

// DeletionResult represents result of a delete operation
type DeletionResult struct {
	// Tells whether the apply operation is completed.
	Completed bool
}

// A Scheduler is used to interact with terraform managed resources
type Scheduler interface {
	Exists(ctx context.Context, mg resource.Terraformed) (bool, error)
	UpdateStatus(ctx context.Context, mg resource.Terraformed) error
	LateInitialize(ctx context.Context, mg resource.Terraformed) (bool, error)
	IsReady(ctx context.Context, mg resource.Terraformed) (bool, error)
	IsUpToDate(ctx context.Context, mg resource.Terraformed) (bool, error)
	GetConnectionDetails(ctx context.Context, mg resource.Terraformed) (managed.ConnectionDetails, error)
	Apply(ctx context.Context, mg resource.Terraformed) (*ApplyResult, error)
	Delete(ctx context.Context, mg resource.Terraformed) (*DeletionResult, error)
}
