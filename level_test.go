package main

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestXName(t *testing.T) {
	mp := map[string]PlatformDef{}
	mp["p1"] = PlatformDef{Fuel: 100}
	mp["p2"] = PlatformDef{Fuel: 200}

	data, err := yaml.Marshal(mp)
	require.NoError(t, err)
	t.Log(string(data))
}
