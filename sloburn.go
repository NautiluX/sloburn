package sloburn

import (
	"github.com/NautiluX/sloburn/alert"
	"github.com/NautiluX/sloburn/kube"
)

func NewBurnAlert(name string, queryGood string, queryValid string, slo float64, labels map[string]string) alert.BurnAlert {
	return alert.BurnAlert{
		Name:                 name,
		QueryGood:            queryGood,
		QueryValid:           queryValid,
		SLO:                  slo / 100,
		PrometheusRuleLabels: labels,
		AlertLabels:          map[string]string{},
		WindowPlaceHolder:    ":window:",
	}
}

func NewBurnAlertWithNamespace(namespace, name string, queryGood string, queryValid string, slo float64, labels map[string]string) alert.BurnAlert {
	return alert.BurnAlert{
		Name:                 name,
		Namespace:            namespace,
		QueryGood:            queryGood,
		QueryValid:           queryValid,
		SLO:                  slo / 100,
		PrometheusRuleLabels: labels,
		AlertLabels:          map[string]string{},
		WindowPlaceHolder:    ":window:",
	}
}

func UpsertAlertsKube(a kube.SLOAlert) {
	kube.UpsertAlerts(a)
}

func DeleteAlertsKube(a kube.SLOAlert) {
	kube.DeleteAlerts(a)
}
