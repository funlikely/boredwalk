package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	// Replace with the URL of the image board page you want to monitor
	imageBoardURL := "https://example.com/imageboard"

	// Perform initial fetch
	knownPosts := make(map[string]bool)
	fetchAndDownload(imageBoardURL, knownPosts)

	// Periodically check for new posts
	for {
		sleepTime := generateRandomSleepTime(120, 240)
		fmt.Printf("Sleeping for %d seconds...\n", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)

		fetchAndDownload(imageBoardURL, knownPosts)
	}
}

func fetchAndDownload(url string, knownPosts map[string]bool) {
	// Fetch the image board page content
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching page:", err)
		return
	}
	defer resp.Body.Close()

	// Parse HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Extract post information and download new images or videos
	extractAndDownload(doc, knownPosts)
}

func extractAndDownload(n *html.Node, knownPosts map[string]bool) {
	if n.Type == html.ElementNode && n.Data == "img" {
		// Extract image URL from img tags
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				imageURL := attr.Val
				postID := strings.TrimPrefix(imageURL, "https://example.com/image/")
				if !knownPosts[postID] {
					// Download the image
					downloadFile(imageURL, "downloads/"+postID+".jpg")
					knownPosts[postID] = true
					fmt.Println("Downloaded new image:", imageURL)
				}
			}
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractAndDownload(c, knownPosts)
	}
}

func downloadFile(url, filepath string) error {
	// Create or open the file for writing
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Make HTTP request to download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	return err
}

func generateRandomSleepTime(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
