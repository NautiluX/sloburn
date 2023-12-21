package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"strings"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type BurnAlert struct {
	Name                 string
	Namespace            string
	SLO                  float64
	QueryGood            string
	QueryValid           string
	WindowPlaceHolder    string
	PrometheusRuleLabels map[string]string
	AlertLabels          map[string]string
}

func (a *BurnAlert) SetNamespace(namespace string) {
	a.Namespace = namespace
}

func (a *BurnAlert) SetWindowPlaceholder(placeholder string) {
	a.WindowPlaceHolder = placeholder
}

func (a *BurnAlert) AddAlertLabels(labels map[string]string) {
	maps.Copy(a.AlertLabels, labels)
}

func (a *BurnAlert) CompilePrometheusRule() monitoringv1.PrometheusRule {
	return monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PrometheusRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
			Labels:    a.PrometheusRuleLabels,
		},
		Spec: monitoringv1.PrometheusRuleSpec{
			Groups: []monitoringv1.RuleGroup{
				{
					Name:  "slo-rules",
					Rules: a.GenerateAlertRules(),
				},
			},
		},
	}
}

func (a *BurnAlert) CompilePrometheusRuleString() string {
	p := a.CompilePrometheusRule()
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err := enc.Encode(&p)
	if err != nil {
		log.Fatalf("error generating prometheus rule yaml: %v", err)
	}
	return buf.String()
}

func (a *BurnAlert) compileQuery(query, window string) string {
	return strings.ReplaceAll(query, a.WindowPlaceHolder, window)
}

func (a *BurnAlert) compileSLIQuery(window string) string {
	return a.compileQuery(a.QueryGood+"/"+a.QueryValid, window)
}

func (a *BurnAlert) CompileAlertExpression(windowLong, windowShort string) string {
	return fmt.Sprintf("%s > %.3f and %s > %.3f",
		a.compileSLIQuery(windowShort), 1-a.SLO,
		a.compileSLIQuery(windowLong), 1-a.SLO,
	)
}

func (a *BurnAlert) GenerateAlertRules() []monitoringv1.Rule {
	durationCritical := monitoringv1.Duration("5m")
	durationWarning := monitoringv1.Duration("1h")
	rules := []monitoringv1.Rule{
		{
			Alert: "SLOBurn" + a.Name + "Critical",
			Expr:  intstr.FromString(a.CompileAlertExpression("1h", "5m")),
			For:   &durationCritical,
			Annotations: map[string]string{
				"message": fmt.Sprintf("High error budget burn for %s over the past 1h and 5m (current value: {{ $value }})", a.Name),
			},
			Labels: mergeLabels(a.AlertLabels, map[string]string{
				"prometheusrule": a.Name,
				"severity":       "critical",
				"longWindow":     "1h",
				"shortWindow":    "5m",
			}),
		},
		{
			Alert: "SLOBurn" + a.Name + "Critical",
			Expr:  intstr.FromString(a.CompileAlertExpression("6h", "30m")),
			For:   &durationCritical,
			Annotations: map[string]string{
				"message": fmt.Sprintf("High error budget burn for %s over the past 6h and 30m(current value: {{ $value }})", a.Name),
			},
			Labels: mergeLabels(a.AlertLabels, map[string]string{
				"prometheusrule": a.Name,
				"severity":       "critical",
				"longWindow":     "6h",
				"shortWindow":    "30m",
			}),
		},
		{
			Alert: "SLOBurn" + a.Name + "Warning",
			Expr:  intstr.FromString(a.CompileAlertExpression("24h", "3h")),
			For:   &durationWarning,
			Annotations: map[string]string{
				"message": fmt.Sprintf("Moderate error budget burn for %s over the past 24h and 3h(current value: {{ $value }})", a.Name),
			},
			Labels: mergeLabels(a.AlertLabels, map[string]string{
				"prometheusrule": a.Name,
				"severity":       "warning",
				"longWindow":     "24h",
				"shortWindow":    "3h",
			}),
		},
		{
			Alert: "SLOBurn" + a.Name + "Warning",
			Expr:  intstr.FromString(a.CompileAlertExpression("3d", "6h")),
			For:   &durationWarning,
			Annotations: map[string]string{
				"message": fmt.Sprintf("Moderate error budget burn for %s over the past 6h and 3d (current value: {{ $value }})", a.Name),
			},
			Labels: mergeLabels(a.AlertLabels, map[string]string{
				"prometheusrule": a.Name,
				"severity":       "warning",
				"longWindow":     "3d",
				"shortWindow":    "6h",
			}),
		},
	}
	return rules
}

func mergeLabels(m1, m2 map[string]string) map[string]string {
	l := map[string]string{}
	maps.Copy(l, m1)
	maps.Copy(l, m2)
	return l
}
