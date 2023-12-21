package kube

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type SLOAlert interface {
	CompilePrometheusRule() monitoringv1.PrometheusRule
	CompilePrometheusRuleString() string
	GetNamespace() string
}

func initKubeClient() *dynamic.DynamicClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("failed to get in-cluster config, trying local config next: %v\n", err.Error())
		err = nil
	}

	if config == nil {
		kubeconfig := getLocalKubeconfig()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return client
}

func getLocalKubeconfig() string {
	var kubeconfig *string
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env
	}
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()
		return *kubeconfig
	}
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	return *kubeconfig

}

func CreateAlerts(a SLOAlert) {
	client := initKubeClient()
	unstructuredPrometheusRule := unstructured.Unstructured{}
	json.Unmarshal(([]byte)(a.CompilePrometheusRuleString()), &unstructuredPrometheusRule)
	res, err := client.Resource(monitoringv1.SchemeGroupVersion.WithResource("prometheusrules")).Namespace(a.GetNamespace()).
		Create(context.Background(), &unstructuredPrometheusRule, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Printf("PrometheusRule created: %v", res)

}
