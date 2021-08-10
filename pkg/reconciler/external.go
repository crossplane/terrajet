package reconciler

import (
	"context"

	"github.com/pkg/errors"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/terrajet/pkg/resource"
	"github.com/crossplane-contrib/terrajet/pkg/scheduler"
)

const (
	errUnexpectedObject = "The managed resource is not an Terraformed resource"
)

// NewTerraformExternal returns a terraform external client
func NewTerraformExternal(e scheduler.Scheduler) *TerraformExternal {
	return &TerraformExternal{
		executor: e,
	}
}

// TerraformExternal manages lifecycle of a Terraform managed resource by implementing
// managed.ExternalClient interface.
type TerraformExternal struct {
	executor scheduler.Scheduler
}

// Observe does an observation for the Terraform managed resource.
func (e *TerraformExternal) Observe(ctx context.Context, mg xpresource.Managed) (managed.ExternalObservation, error) {
	tr, ok := mg.(resource.Terraformed)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedObject)
	}

	exists, err := e.executor.Exists(ctx, tr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to check if exists")
	}
	if !exists {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	lateInitialized, err := e.executor.LateInitialize(ctx, tr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to late init")
	}

	upToDate, err := e.executor.IsUpToDate(ctx, tr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to check if is up to date")
	}

	isReady, err := e.executor.IsReady(ctx, tr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to check if ready")
	}
	if isReady {
		mg.SetConditions(xpv1.Available())
	} else {
		mg.SetConditions(xpv1.Creating())
	}

	if err = e.executor.UpdateStatus(ctx, tr); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to update status")
	}

	conn, err := e.executor.GetConnectionDetails(ctx, tr)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, "failed to get connection details")
	}
	return managed.ExternalObservation{
		ResourceExists:          true,
		ResourceUpToDate:        upToDate,
		ResourceLateInitialized: lateInitialized,
		ConnectionDetails:       conn,
	}, nil
}

// Create creates the Terraform managed resource.
func (e *TerraformExternal) Create(ctx context.Context, mg xpresource.Managed) (managed.ExternalCreation, error) {
	u, err := e.Update(ctx, mg)
	return managed.ExternalCreation{
		ExternalNameAssigned: meta.GetExternalName(mg) != "",
		ConnectionDetails:    u.ConnectionDetails,
	}, err
}

// Update updates the Terraform managed resource.
func (e *TerraformExternal) Update(ctx context.Context, mg xpresource.Managed) (managed.ExternalUpdate, error) {
	tr, ok := mg.(resource.Terraformed)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedObject)
	}

	ar, err := e.executor.Apply(ctx, tr)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "failed to apply")
	}
	return managed.ExternalUpdate{
		ConnectionDetails: ar.ConnectionDetails,
	}, nil
}

// Delete deletes the Terraform managed resource.
func (e *TerraformExternal) Delete(ctx context.Context, mg xpresource.Managed) error {
	tr, ok := mg.(resource.Terraformed)
	if !ok {
		return errors.New(errUnexpectedObject)
	}

	dr, err := e.executor.Delete(ctx, tr)
	if err != nil {
		return errors.Wrap(err, "failed to delete")
	}
	if dr.Completed {
		return nil
	}
	return errors.Errorf("still deleting")
}
