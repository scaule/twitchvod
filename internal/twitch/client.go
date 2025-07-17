package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Client represents a Twitch API client
type Client struct {
	httpClient *http.Client
	clientID   string
	baseURL    string
}

// New creates a new Twitch client
func New(clientID, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		clientID:   clientID,
		baseURL:    baseURL,
	}
}

// ExtractVideoID extracts the video ID from a Twitch URL
func (c *Client) ExtractVideoID(url string) (string, error) {
	// Remove @ symbol if present
	url = strings.TrimPrefix(url, "@")

	// Regex to match Twitch video URLs
	re := regexp.MustCompile(`twitch\.tv/videos/(\d+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Twitch video URL format")
	}

	return matches[1], nil
}

// GetVideoInfo makes the HTTP call to Twitch's GraphQL API and returns the streaming URL
func (c *Client) GetVideoInfo(videoID string) (string, error) {
	// Create the GraphQL request payload
	request := TwitchGraphQLRequest{
		{
			OperationName: "VideoPlayer_VODSeekbarPreviewVideo",
			Variables: GraphQLVariables{
				VideoID: videoID,
			},
			Extensions: GraphQLExtensions{
				PersistedQuery: PersistedQuery{
					Version:    1,
					Sha256Hash: "07e99e4d56c5a7c67117a154777b0baf85a5ffefa393b213f4bc712ccaf85dd6",
				},
			},
		},
		{
			OperationName: "ComscoreStreamingQuery",
			Variables: GraphQLVariables{
				VodID: videoID,
			},
			Extensions: GraphQLExtensions{
				PersistedQuery: PersistedQuery{
					Version:    1,
					Sha256Hash: "e1edae8122517d013405f237ffcc124515dc6ded82480a88daef69c83b53ac01",
				},
			},
		},
		{
			OperationName: "VodChannelLoginQuery",
			Variables: GraphQLVariables{
				VideoID: videoID,
			},
			Extensions: GraphQLExtensions{
				PersistedQuery: PersistedQuery{
					Version:    1,
					Sha256Hash: "0c5feea4dad2565508828f16e53fe62614edf015159df4b3bca33423496ce78e",
				},
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Client-Id", c.clientID)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Parse the JSON response
	var response TwitchGraphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract seekPreviewsURL from the first response (VideoPlayer_VODSeekbarPreviewVideo)
	if len(response) == 0 || response[0].Data.Video.SeekPreviewsURL == "" {
		return "", fmt.Errorf("seekPreviewsURL not found in response")
	}

	return strings.Split(response[0].Data.Video.SeekPreviewsURL, "/storyboards/")[0] + "/chunked/index-dvr.m3u8", nil
}
