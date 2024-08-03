package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/spf13/cobra"
)

var (
	allCaches bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache files",
	Long:  "This command cleans cache files from your system.",
	Run:   cleanCaches,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cache files",
	Long:  "This command lists cache files and their sizes.",
	Run:   listCaches,
}

func init() {
	cleanCmd.Flags().BoolVarP(&allCaches, "all", "a", false, "Clean all cache files")
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(listCmd)
}

func detectCacheDirectories() []string {
	var dirs []string
	switch runtime.GOOS {
	case "darwin":
		dirs = []string{
			filepath.Join(os.Getenv("HOME"), "Library", "Caches"),
		}
	case "linux":
		dirs = []string{
			filepath.Join(os.Getenv("HOME"), ".cache"),
		}
	case "windows":
		dirs = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp"),
		}
	}
	return dirs
}

func cleanCaches(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup
	caches := detectCacheDirectories()

	if allCaches {
		fmt.Println("Cleaning all cache files...")
		for _, cacheDir := range caches {
			files, err := filepath.Glob(filepath.Join(cacheDir, "*"))
			if err != nil {
				fmt.Printf("Error listing files in %s: %v\n", cacheDir, err)
				continue
			}

			for _, file := range files {
				wg.Add(1)
				go func(f string) {
					defer wg.Done()
					fmt.Printf("Cleaning %s...\n", f)
					if err := os.RemoveAll(f); err != nil {
						fmt.Printf("Error removing %s: %v\n", f, err)
					}
				}(file)
			}
		}
	} else {
		fmt.Println("Cleaning specified cache files...")
		for _, cache := range args {
			for _, cacheDir := range caches {
				files, err := filepath.Glob(filepath.Join(cacheDir, "*"+cache+"*"))
				if err != nil {
					fmt.Printf("Error listing files in %s: %v\n", cacheDir, err)
					continue
				}

				for _, file := range files {
					wg.Add(1)
					go func(f string) {
						defer wg.Done()
						fmt.Printf("Cleaning %s...\n", f)
						if err := os.RemoveAll(f); err != nil {
							fmt.Printf("Error removing %s: %v\n", f, err)
						}
					}(file)
				}
			}
		}
	}

	wg.Wait()
	fmt.Println("Cache cleaning completed.")
}

func listCaches(cmd *cobra.Command, args []string) {
	caches := detectCacheDirectories()

	for _, cacheDir := range caches {
		files, err := filepath.Glob(filepath.Join(cacheDir, "*"))
		if err != nil {
			fmt.Printf("Error listing files in %s: %v\n", cacheDir, err)
			continue
		}

		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				fmt.Printf("Error getting info for %s: %v\n", file, err)
				continue
			}
			fmt.Printf("%s: %d bytes\n", filepath.Base(file), info.Size())
		}
	}
}
