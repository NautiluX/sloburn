# SLO Burn

SLO Burn is the missing abstraction layer to simplify and solidify the creation of SLO-based alerts.

Turn your error budget alerts into Go code and never care about the details again.

# SLO Definition

Define a Service Level Objective (SLO) using Go code by setting a target availability and the input PromQL formula that calculates the SLI.

# Multiwindow multi-burn-rate alerts

That's what SLO Burn will generate for you. It's most likely non-readable PromQL which is why you won't want to craft it by hand.

# Example

See example/burnalert. Current state:

Go code:

```
package main

import (
	"fmt"

	"github.com/NautiluX/sloburn"
)

const WindowPlaceHolder string = ":window:"
const queryGood string = "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[" + WindowPlaceHolder + "]))"
const queryValid string = "sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[" + WindowPlaceHolder + "]))"

func main() {
	alert := sloburn.NewBurnAlert("APIServerAvailability", queryGood, queryValid, 99.0, map[string]string{"prometheus": "prometheus-k8s"})
	alert.SetWindowPlaceholder(WindowPlaceHolder)
	fmt.Println(alert.CompilePrometheusRule())
}
```

Output:

```
{
  "kind": "PrometheusRule",
  "apiVersion": "monitoring.coreos.com/v1",
  "metadata": {
    "name": "APIServerAvailability",
    "creationTimestamp": null
  },
  "spec": {
    "groups": [
      {
        "name": "slo-rules",
        "rules": [
          {
            "alert": "SLOBurnAPIServerAvailabilityCritical",
            "expr": "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[1h]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[1h])) \u003e (14.4*0.010) and sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[5m]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[5m])) \u003e (14.4*0.010)",
            "for": "5m",
            "labels": {
              "service": "APIServerAvailability",
              "severity": "critical"
            },
            "annotations": {
              "message": "High error budget burn for APIServerAvailability over the past 1h and 5m (current value: {{ $value }})"
            }
          },
          {
            "alert": "SLOBurnAPIServerAvailabilityCritical",
            "expr": "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[6h]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[6h])) \u003e (6*0.010) and sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[30m]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[30m])) \u003e (6*0.010)",
            "for": "5m",
            "labels": {
              "service": "APIServerAvailability",
              "severity": "critical"
            },
            "annotations": {
              "message": "High error budget burn for APIServerAvailability over the past 6h and 30m(current value: {{ $value }})"
            }
          },
          {
            "alert": "SLOBurnAPIServerAvailabilityWarning",
            "expr": "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[24h]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[24h])) \u003e (3*0.010) and sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[3h]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[3h])) \u003e (3*0.010)",
            "for": "1h",
            "labels": {
              "service": "APIServerAvailability",
              "severity": "warning"
            },
            "annotations": {
              "message": "Moderate error budget burn for APIServerAvailability over the past 24h and 3h(current value: {{ $value }})"
            }
          },
          {
            "alert": "SLOBurnAPIServerAvailabilityWarning",
            "expr": "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[3d]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[3d])) \u003e 0.010 and sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[6h]))/sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[6h])) \u003e 0.010",
            "for": "1h",
            "labels": {
              "service": "APIServerAvailability",
              "severity": "warning"
            },
            "annotations": {
              "message": "Moderate error budget burn for APIServerAvailability over the past 6h and 3d (current value: {{ $value }})"
            }
          }
        ]
      }
    ]
  }
}

```
