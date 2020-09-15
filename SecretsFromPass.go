// Copyright 2020 me!
// SPDX-License-Identifier: Apache-2.0

//go:generate pluginator
package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"sigs.k8s.io/kustomize/api/kv"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	h                *resmap.PluginHelpers
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// directory in pass where our secrets are located
	PassDir string `json:"passdir,omitempty" yaml:"passdir,omitempty"`
	// List of keys to use in database lookups
	Keys []string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

//nolint: golint
//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) passParse() (map[string]string, error) {
	var pass = make(map[string]string)

	var out bytes.Buffer
	path, err := exec.LookPath("pass")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(path, p.PassDir)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	one := strings.Split(out.String(), "\n")
	for _, field := range one {
		if field == p.PassDir {
			continue
		}
		if field == "" {
			continue
		}
		// get rid of that pesky └── & ├── pass likes to output
		two := strings.Split(field, "─")
		three := strings.TrimSpace(two[len(two)-1])
		// get the value of our pass key
		var res bytes.Buffer
		value := exec.Command(path, fmt.Sprintf("%s/%s", p.PassDir, three))
		value.Stdout = &res
		err := value.Run()
		if err != nil {
			return nil, err
		}
		pass[three] = strings.TrimSpace(res.String())
	}

	return pass, nil
}

func (p *plugin) Config(h *resmap.PluginHelpers, c []byte) error {
	p.h = h
	return yaml.Unmarshal(c, p)
}

// The plan here is to convert the plugin's input
// into the format used by the builtin secret generator plugin.
func (p *plugin) Generate() (resmap.ResMap, error) {
	pass, err := p.passParse()
	if err != nil {
		return nil, err
	}
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace
	for _, k := range p.Keys {
		if v, ok := pass[k]; ok {
			args.LiteralSources = append(
				args.LiteralSources, k+"="+v)
		}
	}
	return p.h.ResmapFactory().FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), args)
}
