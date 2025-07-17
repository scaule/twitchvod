package main

import (
	"fmt"
	"os"

	"twitchvod/internal/twitch"
)

func main() {
	// Check if URL parameter is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: twitchvod <twitch_url>")
		fmt.Println("Example: twitchvod https://www.twitch.tv/videos/2515010841")
		os.Exit(1)
	}

	// Get the Twitch URL from command line argument
	twitchURL := os.Args[1]

	// Create a new Twitch client
	twitchClient := twitch.New(os.Getenv("TWITCH_CLIENT_ID"), "https://gql.twitch.tv/gql")

	// Extract video ID from Twitch URL
	videoID, err := twitchClient.ExtractVideoID(twitchURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting video ID: %v\n", err)
		os.Exit(1)
	}

	// Make request to Twitch GraphQL API and get streaming URL
	streamingURL, err := twitchClient.GetVideoInfo(videoID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling Twitch API: %v\n", err)
		os.Exit(1)
	}

	// Output the streaming URL
	fmt.Println(streamingURL)
}
