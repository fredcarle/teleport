/*
Copyright 2022 Gravitational, Inc.

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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gravitational/teleport/tool/tbot/destination"
	"github.com/gravitational/trace"
	"gopkg.in/yaml.v3"
)

// DestinationDirectory is a Destination that writes to the local filesystem
type DestinationDirectory struct {
	Path string `yaml:"path,omitempty"`
}

func (dd *DestinationDirectory) UnmarshalYAML(node *yaml.Node) error {
	var path string
	if err := node.Decode(&path); err == nil {
		dd.Path = path
		return nil
	}

	// Shenanigans to prevent UnmarshalYAML from recursing back to this
	// override (we want to use standard unmarshal behavior for the full
	// struct)
	type rawDirectory DestinationDirectory
	if err := node.Decode((*rawDirectory)(dd)); err != nil {
		return err
	}

	return nil
}

func (dd *DestinationDirectory) CheckAndSetDefaults() error {
	if dd.Path == "" {
		return trace.BadParameter("destination path must not be empty")
	}

	return nil
}

func (d *DestinationDirectory) Write(name string, data []byte, modeHint destination.ModeHint) error {
	// TODO: honor modeHint?
	if err := os.WriteFile(filepath.Join(d.Path, name), data, 0600); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (d *DestinationDirectory) Read(name string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(d.Path, name))
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return b, nil
}

func (d *DestinationDirectory) String() string {
	return fmt.Sprintf("directory %s", d.Path)
}
