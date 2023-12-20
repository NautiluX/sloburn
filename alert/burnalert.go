package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type BurnAlert struct {
	Name              string
	Namespace         string
	SLO               float64
	QueryGood         string
	QueryValid        string
	WindowPlaceHolder string
	Labels            map[string]string
}

func (a *BurnAlert) SetNamespace(namespace string) {
	a.Namespace = namespace
}

func (a *BurnAlert) SetWindowPlaceholder(placeholder string) {
	a.WindowPlaceHolder = placeholder
}

func (a *BurnAlert) CompilePrometheusRule() string {
	p := monitoringv1.PrometheusRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PrometheusRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
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

func (a *BurnAlert) CompileCriticalAlertExpressionShort() string {
	return fmt.Sprintf("%s > (14.4*%.3f) and %s > (14.4*%.3f)",
		a.compileSLIQuery("1h"), 1-a.SLO,
		a.compileSLIQuery("5m"), 1-a.SLO,
	)
}

func (a *BurnAlert) CompileCriticalAlertExpressionLong() string {
	return fmt.Sprintf("%s > (6*%.3f) and %s > (6*%.3f)",
		a.compileSLIQuery("6h"), 1-a.SLO,
		a.compileSLIQuery("30m"), 1-a.SLO,
	)
}

func (a *BurnAlert) CompileWarningAlertExpressionShort() string {
	return fmt.Sprintf("%s > (3*%.3f) and %s > (3*%.3f)",
		a.compileSLIQuery("24h"), 1-a.SLO,
		a.compileSLIQuery("3h"), 1-a.SLO,
	)
}
func (a *BurnAlert) CompileWarningAlertExpressionLong() string {
	return fmt.Sprintf("%s > %.3f and %s > %.3f",
		a.compileSLIQuery("3d"), 1-a.SLO,
		a.compileSLIQuery("6h"), 1-a.SLO,
	)
}

func (a *BurnAlert) GenerateAlertRules() []monitoringv1.Rule {
	durationCritical := monitoringv1.Duration("5m")
	durationWarning := monitoringv1.Duration("1h")
	rules := []monitoringv1.Rule{
		{
			Alert: "SLOBurn" + a.Name + "Critical",
			Expr:  intstr.FromString(a.CompileCriticalAlertExpressionShort()),
			For:   &durationCritical,
			Annotations: map[string]string{
				"message": fmt.Sprintf("High error budget burn for %s over the past 1h and 5m (current value: {{ $value }})", a.Name),
			},
			Labels: map[string]string{
				"service":  a.Name,
				"severity": "critical",
			},
		},
		{
			Alert: "SLOBurn" + a.Name + "Critical",
			Expr:  intstr.FromString(a.CompileCriticalAlertExpressionLong()),
			For:   &durationCritical,
			Annotations: map[string]string{
				"message": fmt.Sprintf("High error budget burn for %s over the past 6h and 30m(current value: {{ $value }})", a.Name),
			},
			Labels: map[string]string{
				"service":  a.Name,
				"severity": "critical",
			},
		},
		{
			Alert: "SLOBurn" + a.Name + "Warning",
			Expr:  intstr.FromString(a.CompileWarningAlertExpressionShort()),
			For:   &durationWarning,
			Annotations: map[string]string{
				"message": fmt.Sprintf("Moderate error budget burn for %s over the past 24h and 3h(current value: {{ $value }})", a.Name),
			},
			Labels: map[string]string{
				"service":  a.Name,
				"severity": "warning",
			},
		},
		{
			Alert: "SLOBurn" + a.Name + "Warning",
			Expr:  intstr.FromString(a.CompileWarningAlertExpressionLong()),
			For:   &durationWarning,
			Annotations: map[string]string{
				"message": fmt.Sprintf("Moderate error budget burn for %s over the past 6h and 3d (current value: {{ $value }})", a.Name),
			},
			Labels: map[string]string{
				"service":  a.Name,
				"severity": "warning",
			},
		},
	}
	return rules
}
