package main

import (
	"flag"
	"fmt"
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

	for _, v := range pods.Items {
		if v.Spec.NodeName == *nodeName {
			kill := true
			if len(v.ObjectMeta.OwnerReferences) > 0 {
				for _, owner := range v.ObjectMeta.OwnerReferences {
					if owner.Kind == "DaemonSet" {
						kill = false
					}
				}
			}
			if kill {
				killPods(v.ObjectMeta.Namespace, v.ObjectMeta.Name, clientset)
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

func killPods(namespace, name string, clientset *kubernetes.Clientset) {
	deleteOptions := &metav1.DeleteOptions{}
	graceSeconds := int64(0)
	deleteOptions.GracePeriodSeconds = &graceSeconds
	fmt.Printf("Killing pod %s in namespace %s\n", name, namespace)
	err := clientset.CoreV1().Pods(namespace).Delete(name, deleteOptions)
	if err != nil {
		panic(err.Error())
	}
}
