package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	node       string
)

func init() {
	if home := homeDir(); home != "" {
		flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&node, "node", "minikube", "Which node name to use for testing")
	flag.Parse()
}

func TestExternalIP(t *testing.T) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	// creates the clientset
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	ip, err := externalIP(c, node)
	if err != nil {
		t.Error(err)
	}
	if ip == "" {
		t.Error("no IP")
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
