package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	m := new()
	require.NotNil(t, m)
	require.NotNil(t, m.Registry)
	require.NotNil(t, m.TotalLatency)
	require.NotNil(t, m.TokensTotal)
	require.NotNil(t, m.RequestsTotal)
}

func TestGetOrCreate(t *testing.T) {
	m1 := GetOrCreate()
	m2 := GetOrCreate()
	require.Same(t, m1, m2)
}

func TestGetRegistry(t *testing.T) {
	reg := GetRegistry()
	require.NotNil(t, reg)
	require.Same(t, GetOrCreate().Registry, reg)
}
