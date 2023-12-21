package sloburn

import (
	"github.com/NautiluX/sloburn/alert"
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
