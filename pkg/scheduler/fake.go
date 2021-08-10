package scheduler

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"

	"github.com/crossplane-contrib/terrajet/pkg/resource"
)

// Fake is a fake Scheduler implementation
type Fake struct{}

// Exists is a fake implementation for Exists function of Scheduler interface.
func (f *Fake) Exists(ctx context.Context, tr resource.Terraformed) (bool, error) {
	fmt.Println("Fake executor checking if exists...")
	return true, nil
}

// UpdateStatus is a fake implementation for UpdateStatus function of Scheduler interface.
func (f *Fake) UpdateStatus(ctx context.Context, tr resource.Terraformed) error {
	/*if err := tr.SetObservation([]byte{}); err != nil {
		return errors.Wrap(err, "failed to set observation")
	}*/
	fmt.Println("Fake executor updating status...")
	return nil
}

// LateInitialize is a fake implementation for LateInitialize function of Scheduler interface.
func (f *Fake) LateInitialize(ctx context.Context, tr resource.Terraformed) (bool, error) {
	/*if err := tr.SetParameters([]byte{}); err != nil {
		return false, errors.Wrap(err, "failed to set parameters")
	}*/
	fmt.Println("Fake executor late initializing...")
	return true, nil
}

// IsReady is a fake implementation for IsReady function of Scheduler interface.
func (f *Fake) IsReady(ctx context.Context, tr resource.Terraformed) (bool, error) {
	fmt.Println("Fake executor checking if ready...")
	return true, nil
}

// IsUpToDate is a fake implementation for IsUpToDate function of Scheduler interface.
func (f *Fake) IsUpToDate(ctx context.Context, tr resource.Terraformed) (bool, error) {
	fmt.Println("Fake executor checking if up to date...")
	return false, nil
}

// GetConnectionDetails is a fake implementation for GetConnectionDetails function of Scheduler interface.
func (f *Fake) GetConnectionDetails(ctx context.Context, tr resource.Terraformed) (managed.ConnectionDetails, error) {
	fmt.Println("Fake executor returning connection details...")
	return managed.ConnectionDetails{}, nil
}

// Apply is a fake implementation for Apply function of Scheduler interface.
func (f *Fake) Apply(ctx context.Context, tr resource.Terraformed) (*ApplyResult, error) {
	b, err := tr.GetParameters()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get parameters")
	}
	fmt.Printf("Fake executor applying with parameters %s\n", b)
	return &ApplyResult{Completed: true}, nil
}

// Delete is a fake implementation for Delete function of Scheduler interface.
func (f *Fake) Delete(ctx context.Context, tr resource.Terraformed) (*DeletionResult, error) {
	fmt.Println("Fake executor deleting...")
	return &DeletionResult{Completed: true}, nil
}
