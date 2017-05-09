package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeClient abstracts interaction with k8s.
type KubeClient struct {
	clientset *kubernetes.Clientset
}

// NewKubeClient is a constructor for the new kube client.
// configPath should be an absolute path.
func NewKubeClient(inCluster bool, configPath string) (*KubeClient, error) {
	clientset, err := newClientset(inCluster, configPath)
	if err != nil {
		return nil, err
	}

	return &KubeClient{
		clientset: clientset,
	}, nil
}

// GetDisks retrives all the disk names used by different persistent volumes.
func (k *KubeClient) GetDisks() ([]string, error) {
	pvs, err := k.clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var disks []string
	for _, pv := range pvs.Items {
		disks = append(disks, pv.Spec.GCEPersistentDisk.PDName)
	}
	return disks, nil
}

func newClientset(inCluster bool, configPath string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("can't build the *in cluster* config: %s", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, fmt.Errorf("can't build the *out of cluster* config: %s", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("can't build kubernetes client: %s", err)
	}

	return clientset, nil
}
