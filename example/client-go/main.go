package main

import (
	"github.com/NautiluX/sloburn"
)

const WindowPlaceHolder string = ":window:"
const queryGood string = "sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[" + WindowPlaceHolder + "]))"
const queryValid string = "sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[" + WindowPlaceHolder + "]))"

func main() {
	alert := sloburn.NewBurnAlert(
		"APIServerAvailability",
		queryGood,
		queryValid,
		99.0,
		map[string]string{"prometheus": "k8s"},
	)
	alert.AddAlertLabels(map[string]string{"service": "API Server"})
	alert.SetNamespace("openshift-monitoring")
	sloburn.UpsertAlertsKube(&alert)
}
