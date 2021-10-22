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

package pipeline

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/muvaf/typewriter/pkg/wrapper"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/crossplane-contrib/terrajet/pkg/pipeline/templates"
)

// NewSetupGenerator returns a new SetupGenerator.
func NewSetupGenerator(rootDir, modulePath string) *SetupGenerator {
	return &SetupGenerator{
		LocalDirectoryPath: filepath.Join(rootDir, "internal", "controller"),
		ModulePath:         modulePath,
	}
}

// SetupGenerator generates controller setup file.
type SetupGenerator struct {
	LocalDirectoryPath string
	ModulePath         string
}

type aliasGVK struct {
	Alias string
	GVK   schema.GroupVersionKind
}

// Generate writes the setup file with the content produced using given
// list of version packages.
func (rg *SetupGenerator) Generate(versionPkgList []string, gvkList []schema.GroupVersionKind) error {
	setupFile := wrapper.NewFile(filepath.Join(rg.ModulePath, "apis"), "apis", templates.SetupTemplate,
		wrapper.WithGenStatement(GenStatement),
		wrapper.WithHeaderPath("hack/boilerplate.go.txt"),
	)
	aliasesGVKs := make([]aliasGVK, len(versionPkgList))
	for i, pkgPath := range versionPkgList {
		aliasesGVKs[i] = aliasGVK{
			Alias: setupFile.Imports.UsePackage(pkgPath),
			GVK:   gvkList[i],
		}
	}
	sort.Slice(aliasesGVKs, func(i, j int) bool {
		return strings.Compare(aliasesGVKs[i].Alias, aliasesGVKs[j].Alias) == -1
	})
	vars := map[string]interface{}{
		"AliasesGVKs": aliasesGVKs,
	}
	filePath := filepath.Join(rg.LocalDirectoryPath, "zz_setup.go")
	return errors.Wrap(setupFile.Write(filePath, vars, os.ModePerm), "cannot write setup file")
}
