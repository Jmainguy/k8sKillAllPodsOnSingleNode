package main

import (
    "fmt"
	"flag"
	"os"
	"path/filepath"
	//"time"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
    
    nodeName := flag.String("nodeName", "", "node to kill all pods on, forcibly")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

    deleteOptions := &metav1.DeleteOptions{}
    graceSeconds := int64(0)
    deleteOptions.GracePeriodSeconds = &graceSeconds
	for _, v := range pods.Items {
		if v.Spec.NodeName == *nodeName {
            fmt.Printf("Killing pod %s in namespace %s\n", v.ObjectMeta.Name, v.ObjectMeta.Namespace)
	        err := clientset.CoreV1().Pods(v.ObjectMeta.Namespace).Delete(v.ObjectMeta.Name, deleteOptions)
    	    if err != nil {
		        panic(err.Error())
	        }
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
