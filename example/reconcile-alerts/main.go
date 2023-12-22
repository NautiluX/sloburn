package main

import (
	"time"

	"github.com/NautiluX/sloburn"
	"github.com/NautiluX/sloburn/alert"
)

func main() {
	reconcileAlerts(getActiveAlerts(), getOutdatedAlerts())
}

func reconcileAlerts(alerts, outdatedAlerts []alert.BurnAlert) {
	for {
		for _, a := range outdatedAlerts {
			sloburn.DeleteAlertsKube(&a)
		}
		for _, a := range alerts {
			sloburn.UpsertAlertsKube(&a)
		}
		time.Sleep(time.Second * 10)
	}
}

func getActiveAlerts() []alert.BurnAlert {
	alerts := []alert.BurnAlert{}
	alerts = append(alerts, getApiServerAlerts()...)
	return alerts
}

func getApiServerAlerts() []alert.BurnAlert {
	alerts := []alert.BurnAlert{}
	apiAlert := sloburn.NewBurnAlert(
		"APIServerAvailability",
		"sum(rate(apiserver_request_total{job=\"kube-apiserver\", code=~\"5..\"}[:window:]))",
		"sum(rate(apiserver_request_total{job=\"kube-apiserver\"}[:window:]))",
		99.0,
		map[string]string{"prometheus": "k8s"},
	)
	apiAlert.AddAlertLabels(map[string]string{"service": "API Server"})
	apiAlert.SetNamespace("openshift-monitoring")
	alerts = append(alerts, apiAlert)
	return alerts
}

func getOutdatedAlerts() []alert.BurnAlert {
	outdatedAlerts := []alert.BurnAlert{}
	outdatedAlerts = append(outdatedAlerts, sloburn.NewBurnAlertWithNamespace("default", "SomeOldSLO", "", "", 0, map[string]string{}))
	return outdatedAlerts
}
