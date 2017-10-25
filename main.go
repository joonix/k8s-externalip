package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var errNotFound = errors.New("not found")

func externalIP(c *kubernetes.Clientset) (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	node, err := c.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	for _, addr := range node.Status.Addresses {
		if addr.Type == v1.NodeExternalIP {
			return addr.Address, nil
		}
	}
	return "", errNotFound
}

func updateConfig(c *kubernetes.Clientset, namespace, configmap, filename, placeholder, ip string) error {
	conf, err := c.CoreV1().ConfigMaps(namespace).Get(configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}
	input := conf.Data[filepath.Base(filename)]
	output := strings.Replace(string(input), placeholder, ip, -1)
	if len(output) > 0 {
		if err := ioutil.WriteFile(filename, []byte(output), 0); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	namespace := flag.String("namespace", "default", "Which namespace to look for configmap in")
	configmap := flag.String("configmap", "", "Which config map to read config from")
	filename := flag.String("filename", "", "The file name to write config to")
	placeholder := flag.String("placeholder", "K8S_EXTERNALADDRESS", "What string to search and replace from config")
	flag.Parse()

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ip, err := externalIP(clientset)
	if err != nil {
		panic(err.Error())
	}

	if err = updateConfig(clientset, *namespace, *configmap, *filename, *placeholder, ip); err != nil {
		panic(err.Error())
	}
	fmt.Println("External IP set to:", ip)
}
