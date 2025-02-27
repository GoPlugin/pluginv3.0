package cmd

import "github.com/goplugin/plugin-common/pkg/capabilities"

type TypeInfo struct {
	CapabilityType   capabilities.CapabilityType
	RootType         string
	SchemaID         string
	SchemaOutputFile string
}
