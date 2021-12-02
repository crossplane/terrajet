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

package exec

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
)

// Exec interface provides a Command interface. By this way, real or fake implementation of functions is determined.
type Exec interface {
	CommandContext(ctx context.Context, name string, arg ...string) Command
	SetEnv(env []string)
	SetDir(dir string)
}

// OsExec struct implements the CommandContext function that calls the function from the exec package.
type OsExec struct {
	dir string
	env []string
}

// SetDir lets you set the dir of OsExec.
func (r *OsExec) SetDir(dir string) {
	r.dir = dir
}

// SetEnv lets you set the env of OsExec.
func (r *OsExec) SetEnv(env []string) {
	r.env = env
}

// CommandContext calls the function from exec package, then configure the *Cmd object with `dir` and `env`.
// It returns a *Cmd object thus, the real implementation of CombinedOutput will be called.
func (r *OsExec) CommandContext(ctx context.Context, name string, arg ...string) Command {
	cmd := exec.CommandContext(ctx, name, arg...)

	configureCmd(cmd, r.dir, r.env)

	return cmd
}

// FakeExec struct implements the CommandContext function that calls fake one.
type FakeExec struct {
	dir    string
	env    []string
	stdErr []byte
	stdOut []byte
	err    error
}

// SetDir lets you set the dir of FakeExec.
func (m *FakeExec) SetDir(dir string) {
	m.dir = dir
}

// SetEnv lets you set the env of FakeExec.
func (m *FakeExec) SetEnv(env []string) {
	m.env = env
}

// CommandContext function manipulates the `StdErr` and `stdOut` .
// It returns a FakeCommand object thus, the fake implementation of CombinedOutput will be called.
func (m *FakeExec) CommandContext(_ context.Context, _ string, _ ...string) Command {
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

// FakeCommandOption lets you configure the FakeCommand.
type FakeCommandOption func(f *FakeCommand)

// WithStdOut lets you set the stdOut of FakeCommand.
func WithStdOut(stdOut string) FakeCommandOption {
	return func(f *FakeCommand) {
		f.out = []byte(stdOut)
	}
}

// WithStdErr lets you set the stdErr of FakeCommand.
func WithStdErr(stdErr string) FakeCommandOption {
	return func(f *FakeCommand) {
		f.out = []byte(stdErr)
	}
}

// WithErr lets you set the error of FakeCommand.
func WithErr(err string) FakeCommandOption {
	return func(f *FakeCommand) {
		f.err = errors.New(err)
	}
}

// NewFakeCommand returns a new FakeCommand.
func NewFakeCommand(opts ...FakeCommandOption) *FakeCommand {
	f := &FakeCommand{}

	for _, o := range opts {
		o(f)
	}

	return f
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
