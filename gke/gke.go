package gke

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"path/filepath"

	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getGoogleClient() (*container.Service, error) {
	ctx := context.Background()
	client, err := container.NewService(ctx, option.WithCredentialsFile("path/to/your/service-account-file.json"))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateCluster(projectID, zone, clusterName string) {
	client, err := getGoogleClient()
	if err != nil {
		log.Fatalf("Failed to create Google client: %v", err)
	}

	cluster := &container.Cluster{
		Name:             clusterName,
		InitialNodeCount: 3,
		NodeConfig: &container.NodeConfig{
			MachineType: "n1-standard-1",
		},
	}

	req := &container.CreateClusterRequest{
		Cluster: cluster,
	}

	op, err := client.Projects.Zones.Clusters.Create(projectID, zone, req).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Failed to create cluster: %v", err)
	}

	fmt.Printf("Cluster creation in progress: %s\n", op.Name)
}

func DeleteCluster(projectID, zone, clusterName string) {
	client, err := getGoogleClient()
	if err != nil {
		log.Fatalf("Failed to create Google client: %v", err)
	}

	op, err := client.Projects.Zones.Clusters.Delete(projectID, zone, clusterName).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Failed to delete cluster: %v", err)
	}

	fmt.Printf("Cluster deletion in progress: %s\n", op.Name)
}

func getKubeClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func ListPods(namespace string) {
	clientset, err := getKubeClient()
	if err != nil {
		log.Fatalf("Failed to get Kubernetes client: %v", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list pods: %v", err)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}

func CheckPodsHealth(namespace string) {
	clientset, err := getKubeClient()
	if err != nil {
		log.Fatalf("Failed to get Kubernetes client: %v", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list pods: %v", err)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}

func BuildAndUploadDockerImage(projectID, imageName, dockerfilePath string) {
	cmd := exec.Command("docker", "build", "-t", imageName, dockerfilePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to build Docker image: %v\n%s", err, string(output))
	}

	gcrImage := fmt.Sprintf("gcr.io/%s/%s", projectID, imageName)
	cmd = exec.Command("docker", "tag", imageName, gcrImage)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to tag Docker image: %v\n%s", err, string(output))
	}

	cmd = exec.Command("docker", "push", gcrImage)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("Failed to push Docker image: %v\n%s", err, string(output))
	}

	fmt.Printf("Docker image pushed to: %s\n", gcrImage)
}
