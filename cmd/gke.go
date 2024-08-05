package cmd

import (
	"github.com/Eldrago12/advanced-cli-tool/gke"

	"github.com/spf13/cobra"
)

var (
	clusterName    string
	podName        string
	imageName      string
	dockerfilePath string
	namespace      string
)

var gkeCmd = &cobra.Command{
	Use:   "gke",
	Short: "GKE operations",
}

var gkeCreateClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create a GKE cluster",
	Run: func(cmd *cobra.Command, args []string) {
		gke.CreateCluster(projectID, zone, clusterName)
	},
}

var gkeDeleteClusterCmd = &cobra.Command{
	Use:   "delete-cluster",
	Short: "Delete a GKE cluster",
	Run: func(cmd *cobra.Command, args []string) {
		gke.DeleteCluster(projectID, zone, clusterName)
	},
}

var gkeListPodsCmd = &cobra.Command{
	Use:   "list-pods",
	Short: "List GKE pods",
	Run: func(cmd *cobra.Command, args []string) {
		gke.ListPods(namespace)
	},
}

var gkeCheckPodsHealthCmd = &cobra.Command{
	Use:   "check-pods-health",
	Short: "Check GKE pods health",
	Run: func(cmd *cobra.Command, args []string) {
		gke.CheckPodsHealth(namespace)
	},
}

var gkeBuildAndUploadDockerImageCmd = &cobra.Command{
	Use:   "build-upload-docker",
	Short: "Build and upload Docker image to GCR",
	Run: func(cmd *cobra.Command, args []string) {
		gke.BuildAndUploadDockerImage(projectID, imageName, dockerfilePath)
	},
}

func init() {
	gkeCreateClusterCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Cluster name")
	gkeDeleteClusterCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Cluster name")
	gkeListPodsCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace")
	gkeCheckPodsHealthCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace")
	gkeBuildAndUploadDockerImageCmd.Flags().StringVarP(&imageName, "image", "i", "", "Docker image name")
	gkeBuildAndUploadDockerImageCmd.Flags().StringVarP(&dockerfilePath, "dockerfile", "d", ".", "Dockerfile path")

	gkeCmd.AddCommand(gkeCreateClusterCmd)
	gkeCmd.AddCommand(gkeDeleteClusterCmd)
	gkeCmd.AddCommand(gkeListPodsCmd)
	gkeCmd.AddCommand(gkeCheckPodsHealthCmd)
	gkeCmd.AddCommand(gkeBuildAndUploadDockerImageCmd)

	rootCmd.AddCommand(gkeCmd)
}
