package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// uses the current context in kubeconfig
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: ""}
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	client, err := kubernetes.NewForConfig(config)
	fmt.Println(client)

	// port forward tiller
	tillerTunnel, _ := portforwarder.New("kube-system", client, config)
	host := fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)

	helmClient := helm.NewClient(helm.Host(host))
	resp, _ := helmClient.ListReleases()
	for _, release := range resp.Releases {
		fmt.Println(release.GetName())
	}
	if err != nil {
		panic(err.Error())
	}
	// pods, err := clientset.CoreV1().Pods("k8ssta").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	//     panic(err.Error())
	// }
	// fmt.Println(pods.APIVersion)
	// time.Sleep(10 * time.Second)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
