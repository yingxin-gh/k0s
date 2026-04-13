// SPDX-FileCopyrightText: 2022 k0s authors
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponents_SortedAndUnique(t *testing.T) {
	for key, components := range map[string][]string{"available": availableComponents, "deprecated": deprecatedComponents} {
		expected := slices.Clone(components)
		slices.Sort(expected)

		assert.Equal(t, expected, components, key+" components aren't sorted")

		expected = slices.Compact(expected)
		assert.Equal(t, expected, components, key+" components contain duplicates")
	}
}

func TestControllerOptions_Normalize(t *testing.T) {
	t.Run("failsOnUnknownComponent", func(t *testing.T) {
		disabled := []string{"i-dont-exist"}

		underTest := ControllerOptions{DisableComponents: disabled}
		err := underTest.Normalize()

		assert.ErrorContains(t, err, "unknown component i-dont-exist")
	})

	t.Run("acceptsDeprecatedComponents", func(t *testing.T) {
		disabled := []string{"worker-config"}

		underTest := ControllerOptions{DisableComponents: disabled}
		err := underTest.Normalize()

		require.NoError(t, err)
		assert.Equal(t, disabled, underTest.DisableComponents)
	})

	t.Run("removesDuplicateComponents", func(t *testing.T) {
		disabled := []string{"helm", "kube-proxy", "coredns", "kube-proxy", "autopilot"}
		expected := []string{"helm", "kube-proxy", "coredns", "autopilot"}

		underTest := ControllerOptions{DisableComponents: disabled}
		err := underTest.Normalize()

		require.NoError(t, err)
		assert.Equal(t, expected, underTest.DisableComponents)
	})
}

func TestLogLevelsFlagSet(t *testing.T) {
	t.Run("full_input", func(t *testing.T) {
		var underTest logLevelsFlag
		assert.NoError(t, underTest.Set("kubelet=a,kube-scheduler=b,kube-controller-manager=c,kube-apiserver=d,konnectivity-server=e,etcd=f,containerd=g"))
		assert.Equal(t, logLevelsFlag{
			Containerd:            "g",
			Etcd:                  "f",
			Konnectivity:          "e",
			KubeAPIServer:         "d",
			KubeControllerManager: "c",
			KubeScheduler:         "b",
			Kubelet:               "a",
		}, underTest)
		assert.Equal(t, "[containerd=g,etcd=f,konnectivity-server=e,kube-apiserver=d,kube-controller-manager=c,kube-scheduler=b,kubelet=a]", underTest.String())
	})

	t.Run("partial_input", func(t *testing.T) {
		var underTest logLevelsFlag
		assert.NoError(t, underTest.Set("[kubelet=a,etcd=b]"))
		assert.Equal(t, logLevelsFlag{
			Containerd:            "info",
			Etcd:                  "b",
			Konnectivity:          "1",
			KubeAPIServer:         "1",
			KubeControllerManager: "1",
			KubeScheduler:         "1",
			Kubelet:               "a",
		}, underTest)
	})

	for _, test := range []struct {
		name, input, msg string
	}{
		{"unknown_component", "random=debug", `unknown component name: "random"`},
		{"empty_component_name", "=info", "component name cannot be empty"},
		{"unknown_component_only", "random", `must be of format component=level: "random"`},
		{"no_equals", "kube-apiserver", `must be of format component=level: "kube-apiserver"`},
		{"mixed_valid_and_invalid", "containerd=info,random=debug", `unknown component name: "random"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			var underTest logLevelsFlag
			assert.ErrorContains(t, underTest.Set(test.input), test.msg)
		})
	}
}
