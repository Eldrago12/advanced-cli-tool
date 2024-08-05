package cmd

import (
	"github.com/Eldrago12/advanced-cli-tool/gcs"
	"github.com/spf13/cobra"
)

var bucketName string

var gcsCmd = &cobra.Command{
	Use:   "gcs",
	Short: "GCS operations",
}

var gcsUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files to GCS bucket",
	Run: func(cmd *cobra.Command, args []string) {
		gcs.UploadFilesToBucket(bucketName)
	},
}

func init() {
	gcsUploadCmd.Flags().StringVarP(&bucketName, "bucket", "b", "", "GCS bucket name")
	gcsUploadCmd.MarkFlagRequired("bucket")

	gcsCmd.AddCommand(gcsUploadCmd)
	rootCmd.AddCommand(gcsCmd)
}
