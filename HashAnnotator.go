package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"sigs.k8s.io/kustomize/api/filters/annotations"
	"sigs.k8s.io/kustomize/api/ifc"
	"sigs.k8s.io/kustomize/api/provider"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type hashannotator struct {
	Target     *types.Selector   `json:"target,omitempty" yaml:"target,omitempty"`
	Resource   *types.Selector   `json:"resource,omitempty" yaml:"resource,omitempty"`
	FieldSpecs []types.FieldSpec `json:"fieldSpecs,omitempty" yaml:"fieldSpecs,omitempty"`
	hasher     ifc.KustHasher
}

func main() {
	manifest, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read in manifest: %q", err)
		os.Exit(1)
	}

	depProvider := provider.NewDefaultDepProvider()
	rf := depProvider.GetResourceFactory()
	rmf := resmap.NewFactory(rf)

	m, err := rmf.NewResMapFromBytes(manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse manifest: %q", err)
		os.Exit(1)
	}

	var p hashannotator

	defaultFieldSpecs := types.FieldSpec{Path: "metadata/annotations", CreateIfNotPresent: true}
	p.FieldSpecs = []types.FieldSpec{defaultFieldSpecs}
	p.hasher = rmf.RF().Hasher()

	config, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read in config: %q", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(config, &p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling config content: %q \n%s\n", err, config)
		os.Exit(1)
	}

	targets, err := m.Select(*p.Target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting target: %q", err)
		os.Exit(1)
	}

	resources, err := m.Select(*p.Resource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error selecting resource: %q", err)
		os.Exit(1)
	}

	an := make(map[string]string)
	for _, res := range resources {
		h, err := res.Hash(p.hasher)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error hashing resource: %q", err)
			os.Exit(1)
		}
		key := fmt.Sprintf("kustomize.config.k8s.io/hash-%s-%s-%s", res.GetNamespace(), strings.ToLower(res.GetKind()), res.GetName())
		an[key] = h
	}

	for _, target := range targets {
		err = target.ApplyFilter(annotations.Filter{
			Annotations: an,
			FsSlice:     p.FieldSpecs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error applying annotations filter: %q", err)
			os.Exit(1)
		}
	}

	output, err := m.AsYaml()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting resource to yaml: %q", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(output))
}
