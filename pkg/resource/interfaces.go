/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// Observable structs can get and set observations in the form of Terraform JSON.
type Observable interface {
	GetObservation() (map[string]interface{}, error)
	SetObservation(map[string]interface{}) error
	GetID() string
}

// Parameterizable structs can get and set parameters of the managed resource
// using map form of Terraform JSON.
type Parameterizable interface {
	GetParameters() (map[string]interface{}, error)
	SetParameters(map[string]interface{}) error
}

// MetadataProvider provides Terraform metadata for the Terraform managed
// resource.
type MetadataProvider interface {
	GetTerraformResourceType() string
	GetTerraformSchemaVersion() int
	GetConnectionDetailsMapping() map[string]string
}

// LateInitializer late-initializes the managed resource from observed Terraform
// state.
type LateInitializer interface {
	// LateInitialize this Terraformed resource using its observed tfState.
	// returns True if the there are any spec changes for the resource.
	LateInitialize(attrs []byte) (bool, error)
}

// Terraformed is a Kubernetes object representing a concrete terraform managed
// resource.
type Terraformed interface {
	resource.Managed

	MetadataProvider
	Observable
	Parameterizable
	LateInitializer
}
