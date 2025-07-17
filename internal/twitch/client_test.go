package twitch

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	clientID := "test_client_id"
	baseURL := "https://gql.twitch.tv/gql"
	client := New(clientID, baseURL)

	if client.clientID != clientID {
		t.Errorf("Expected clientID %s, got %s", clientID, client.clientID)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized, got nil")
	}

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}
}

func TestExtractVideoID(t *testing.T) {
	client := New("test_client_id", "https://gql.twitch.tv/gql")

	tests := []struct {
		name     string
		url      string
		expected string
		hasError bool
	}{
		{
			name:     "Valid Twitch video URL",
			url:      "https://www.twitch.tv/videos/2515010841",
			expected: "2515010841",
			hasError: false,
		},
		{
			name:     "Valid Twitch video URL with @ prefix",
			url:      "@https://www.twitch.tv/videos/1234567890",
			expected: "1234567890",
			hasError: false,
		},
		{
			name:     "Invalid URL format",
			url:      "https://www.twitch.tv/channel/123",
			expected: "",
			hasError: true,
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: "",
			hasError: true,
		},
		{
			name:     "Non-Twitch URL",
			url:      "https://www.youtube.com/watch?v=123",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.ExtractVideoID(tt.url)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for URL %s, but got none", tt.url)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for URL %s: %v", tt.url, err)
				}
				if result != tt.expected {
					t.Errorf("Expected video ID %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestGetVideoInfo_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept header 'application/json', got %s", r.Header.Get("Accept"))
		}

		if r.Header.Get("Client-Id") != "test_client_id" {
			t.Errorf("Expected Client-Id header 'test_client_id', got %s", r.Header.Get("Client-Id"))
		}

		// Verify request body contains expected operations
		var request TwitchGraphQLRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		if len(request) != 3 {
			t.Errorf("Expected 3 operations, got %d", len(request))
		}

		// Check if the first operation is VideoPlayer_VODSeekbarPreviewVideo
		if request[0].OperationName != "VideoPlayer_VODSeekbarPreviewVideo" {
			t.Errorf("Expected first operation to be 'VideoPlayer_VODSeekbarPreviewVideo', got %s", request[0].OperationName)
		}

		// Return a mock response
		response := TwitchGraphQLResponse{
			{
				Data: GraphQLData{
					Video: GraphQLVideo{
						SeekPreviewsURL: "https://fake-cdn.example.com/video/1234567890/storyboards/1234567890-storyboard-0.jpg",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL
	client := New("test_client_id", server.URL)
	client.httpClient = server.Client()

	videoID := "2515010841"
	streamingURL, err := client.GetVideoInfo(videoID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedBase := "https://fake-cdn.example.com/video/1234567890"
	expectedURL := expectedBase + "/chunked/index-dvr.m3u8"

	if streamingURL != expectedURL {
		t.Errorf("Expected streaming URL %s, got %s", expectedURL, streamingURL)
	}
}

func TestGetVideoInfo_EmptyResponse(t *testing.T) {
	// Create a test server that returns empty response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TwitchGraphQLResponse{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New("test_client_id", server.URL)
	client.httpClient = server.Client()

	_, err := client.GetVideoInfo("2515010841")

	if err == nil {
		t.Error("Expected error for empty response, got none")
	}

	if !strings.Contains(err.Error(), "seekPreviewsURL not found in response") {
		t.Errorf("Expected error message to contain 'seekPreviewsURL not found in response', got %s", err.Error())
	}
}

func TestGetVideoInfo_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := New("test_client_id", server.URL)
	client.httpClient = server.Client()

	_, err := client.GetVideoInfo("2515010841")

	if err == nil {
		t.Error("Expected error for invalid JSON, got none")
	}

	if !strings.Contains(err.Error(), "error parsing response") {
		t.Errorf("Expected error message to contain 'error parsing response', got %s", err.Error())
	}
}

func TestGetVideoInfo_HTTPError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := New("test_client_id", server.URL)
	client.httpClient = server.Client()

	_, err := client.GetVideoInfo("2515010841")

	if err == nil {
		t.Error("Expected error for HTTP error, got none")
	}
}

func TestGetVideoInfo_EmptySeekPreviewsURL(t *testing.T) {
	// Create a test server that returns response with empty SeekPreviewsURL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TwitchGraphQLResponse{
			{
				Data: GraphQLData{
					Video: GraphQLVideo{
						SeekPreviewsURL: "",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New("test_client_id", server.URL)
	client.httpClient = server.Client()

	_, err := client.GetVideoInfo("2515010841")

	if err == nil {
		t.Error("Expected error for empty SeekPreviewsURL, got none")
	}

	if !strings.Contains(err.Error(), "seekPreviewsURL not found in response") {
		t.Errorf("Expected error message to contain 'seekPreviewsURL not found in response', got %s", err.Error())
	}
}

func TestGraphQLRequestStructure(t *testing.T) {
	videoID := "2515010841"

	// Create the request manually to test structure
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

	// Test JSON marshaling
	jsonData, err := json.Marshal(request)
	if err != nil {
		t.Errorf("Failed to marshal request: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledRequest TwitchGraphQLRequest
	if err := json.Unmarshal(jsonData, &unmarshaledRequest); err != nil {
		t.Errorf("Failed to unmarshal request: %v", err)
	}

	// Verify structure
	if len(unmarshaledRequest) != 3 {
		t.Errorf("Expected 3 operations, got %d", len(unmarshaledRequest))
	}

	expectedOperations := []string{
		"VideoPlayer_VODSeekbarPreviewVideo",
		"ComscoreStreamingQuery",
		"VodChannelLoginQuery",
	}

	for i, expectedOp := range expectedOperations {
		if unmarshaledRequest[i].OperationName != expectedOp {
			t.Errorf("Expected operation %s at index %d, got %s", expectedOp, i, unmarshaledRequest[i].OperationName)
		}
	}
}

func TestStreamingURLConstruction(t *testing.T) {
	tests := []struct {
		name            string
		seekPreviewsURL string
		expected        string
	}{
		{
			name:            "Standard URL",
			seekPreviewsURL: "https://fake-cdn.example.com/video/1234567890/storyboards/1234567890-storyboard-0.jpg",
			expected:        "https://fake-cdn.example.com/video/1234567890/chunked/index-dvr.m3u8",
		},
		{
			name:            "URL without storyboards",
			seekPreviewsURL: "https://example.com/video/123",
			expected:        "https://example.com/video/123/chunked/index-dvr.m3u8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Split(tt.seekPreviewsURL, "/storyboards/")[0] + "/chunked/index-dvr.m3u8"
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
