// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package targetallocator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
	"github.com/open-telemetry/opentelemetry-operator/apis/v1beta1"
	"github.com/open-telemetry/opentelemetry-operator/internal/config"
)

func TestDesiredServiceMonitors(t *testing.T) {
	ta := v1alpha1.TargetAllocator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-instance",
			Namespace: "my-namespace",
		},
		Spec: v1alpha1.TargetAllocatorSpec{
			OpenTelemetryCommonFields: v1beta1.OpenTelemetryCommonFields{
				Tolerations: testTolerationValues,
			},
		},
	}
	cfg := config.New()

	params := Params{
		TargetAllocator: ta,
		Config:          cfg,
		Log:             logger,
	}

	actual := ServiceMonitor(params)
	assert.NotNil(t, actual)
	assert.Equal(t, fmt.Sprintf("%s-targetallocator", params.TargetAllocator.Name), actual.Name)
	assert.Equal(t, params.TargetAllocator.Namespace, actual.Namespace)
	assert.Equal(t, "targetallocation", actual.Spec.Endpoints[0].Port)

}

func TestDesiredServiceMonitorsWithEmptyExtraLabels(t *testing.T) {
	ta := v1alpha1.TargetAllocator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-instance",
			Namespace: "my-namespace",
		},
		Spec: v1alpha1.TargetAllocatorSpec{
			OpenTelemetryCommonFields: v1beta1.OpenTelemetryCommonFields{
				Tolerations: testTolerationValues,
			},
			Observability: v1beta1.ObservabilitySpec{
				Metrics: v1beta1.MetricsConfigSpec{
					EnableMetrics: true,
					ExtraLabels:   map[string]string{},
				},
			},
		},
	}
	cfg := config.New()

	params := Params{
		TargetAllocator: ta,
		Config:          cfg,
		Log:             logger,
	}

	actual := ServiceMonitor(params)
	assert.NotNil(t, actual)
	assert.Equal(t, fmt.Sprintf("%s-targetallocator", params.TargetAllocator.Name), actual.Name)
	assert.Equal(t, params.TargetAllocator.Spec.Observability.Metrics.ExtraLabels, map[string]string{})
}

func TestDesiredServiceMonitorsWithExtraLabels(t *testing.T) {
	ta := v1alpha1.TargetAllocator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-instance",
			Namespace: "my-namespace",
		},
		Spec: v1alpha1.TargetAllocatorSpec{
			OpenTelemetryCommonFields: v1beta1.OpenTelemetryCommonFields{
				Tolerations: testTolerationValues,
			},
			Observability: v1beta1.ObservabilitySpec{
				Metrics: v1beta1.MetricsConfigSpec{
					EnableMetrics: true,
					ExtraLabels: map[string]string{
						"prometheus":    "kube-prometheus",
						"team":          "platform",
						"environment":   "production",
						"custom.io/key": "custom-value",
					},
				},
			},
		},
	}
	cfg := config.New()
	params := Params{
		TargetAllocator: ta,
		Config:          cfg,
		Log:             logger,
	}

	actual := ServiceMonitor(params)
	assert.NotNil(t, actual)
	expectedLabels := map[string]string{
		"app.kubernetes.io/component":  "opentelemetry-targetallocator",
		"app.kubernetes.io/instance":   "my-namespace.my-instance",
		"app.kubernetes.io/managed-by": "opentelemetry-operator",
		"app.kubernetes.io/part-of":    "opentelemetry",
		"app.kubernetes.io/name":       "my-instance-targetallocator",
		"app.kubernetes.io/version":    "latest",
		"prometheus":                   "kube-prometheus",
		"team":                         "platform",
		"environment":                  "production",
		"custom.io/key":                "custom-value",
	}
	assert.Equal(t, expectedLabels, actual.ObjectMeta.Labels)

	expectedSelectorLabels := map[string]string{
		"app.kubernetes.io/component":  "opentelemetry-targetallocator",
		"app.kubernetes.io/instance":   "my-namespace.my-instance",
		"app.kubernetes.io/managed-by": "opentelemetry-operator",
		"app.kubernetes.io/name":       "my-instance-targetallocator",
		"app.kubernetes.io/part-of":    "opentelemetry",
	}
	assert.Equal(t, expectedSelectorLabels, actual.Spec.Selector.MatchLabels)
}
