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

package controller

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// error messages
	errMatchGVK = "failed to match GVK"

	fmtGVK = "%s/%s.%s"
)

// IsAPIEnabled returns `true` if the specified GVK matches any
// regular expression specified in enabledAPIs slice, or if enabledAPIs
// has zero-length.
func IsAPIEnabled(gvk schema.GroupVersionKind, enabledAPIs []string) (bool, error) {
	if len(enabledAPIs) == 0 {
		return true, nil
	}

	gvkStr := fmt.Sprintf(fmtGVK, gvk.Group, gvk.Version, gvk.Kind)
	for _, r := range enabledAPIs {
		ok, err := regexp.MatchString(r, gvkStr)
		if err != nil {
			return false, errors.Wrap(err, errMatchGVK)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
