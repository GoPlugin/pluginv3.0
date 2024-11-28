package logeventcap

import _ "github.com/goplugin/plugin-common/pkg/capabilities/cli/cmd" // Required so that the tool is available to be run in go generate below.

//go:generate go run github.com/goplugin/plugin-common/pkg/capabilities/cli/cmd/generate-types --dir $GOFILE
