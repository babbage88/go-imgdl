package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	url    string
	output string
)

var rootCmd = &cobra.Command{
	Use:   "image-downloader",
	Short: "Downloads an image using custom headers",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Fprintln(os.Stderr, "Error: --url is required")
			cmd.Usage()
			os.Exit(1)
		}

		err := downloadImage(url, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Downloaded successfully to %s\n", output)
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the image to download (required)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "output.jpg", "Output filename (default: output.jpg)")

	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
		os.Exit(1)
	}
}

func downloadImage(url, output string) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Emulate browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://uhdpaper.com/")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
