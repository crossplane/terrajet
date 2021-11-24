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

package terraform

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

// Exec interface provides a Command interface. By this way, real or fake implementation of functions is determined.
type Exec interface {
	CommandContext(ctx context.Context, name string, dir string, env []string, arg ...string) Command
}

// RealExec struct implements the CommandContext function that calls the function from the exec package.
type RealExec struct{}

// CommandContext calls the function from exec package, then configure the *Cmd object with `dir` and `env`.
// It returns a *Cmd object thus, the real implementation of CombinedOutput will be called.
func (r *RealExec) CommandContext(ctx context.Context, name string, dir string, env []string, arg ...string) Command {
	cmd := exec.CommandContext(ctx, name, arg...)

	configureCmd(cmd, dir, env)

	return cmd
}

// FakeExec struct implements the CommandContext function that calls fake one.
type FakeExec struct {
	stdErr []byte
	stdOut []byte
	err    error
}

// CommandContext function manipulates the `StdErr` and `stdOut` .
// It returns a FakeCommand object thus, the fake implementation of CombinedOutput will be called.
func (m *FakeExec) CommandContext(_ context.Context, _ string, _ string, _ []string, _ ...string) Command {
	b := bytes.Buffer{}

	if m.stdOut != nil {
		if _, err := b.Write(m.stdOut); err != nil {
			panic(err)
		}
	}

	if m.stdErr != nil {
		if _, err := b.Write(m.stdErr); err != nil {
			panic(err)
		}
	}

	return &FakeCommand{
		out: b.Bytes(),
		err: m.err,
	}
}

// Command interface
type Command interface {
	CombinedOutput() ([]byte, error)
}

// FakeCommand struct implements the CombinedOutput function that calls fake one.
type FakeCommand struct {
	out []byte
	err error
}

// CombinedOutput returns the output and error.
func (c *FakeCommand) CombinedOutput() ([]byte, error) {
	return c.out, c.err
}

func configureCmd(cmd *exec.Cmd, dir string, env []string) {
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}

	cmd.Env = append(cmd.Env, env...)
	cmd.Dir = dir
}
