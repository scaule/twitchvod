# Twitch VOD Stream URL Extractor - RESEARCH PROJECT

⚠️ **IMPORTANT DISCLAIMER** ⚠️

**THIS IS RESEARCH CODE ONLY - NOT FOR PRODUCTION USE**

This project is **EXCLUSIVELY** for educational and research purposes. It demonstrates how to interact with Twitch's GraphQL API and extract streaming URLs from video metadata.

## 🚨 LEGAL AND ETHICAL DISCLAIMERS

- **DO NOT USE THIS TOOL TO DOWNLOAD OR REDISTRIBUTE TWITCH CONTENT**
- **DO NOT USE THIS TOOL FOR COPYRIGHT INFRINGEMENT**
- **DO NOT USE THIS TOOL TO CIRCUMVENT TWITCH'S TERMS OF SERVICE**
- **RESPECT CONTENT CREATORS' RIGHTS AND TWITCH'S PLATFORM POLICIES**

This code is provided as-is for **ACADEMIC RESEARCH** and **LEARNING PURPOSES ONLY**. Users are responsible for ensuring their use complies with all applicable laws, Twitch's Terms of Service, and content creators' rights.

## Research Purpose

This tool demonstrates:
- How to extract video IDs from Twitch URLs using regex
- How to make authenticated requests to Twitch's GraphQL API
- How to parse GraphQL responses to extract streaming metadata
- Go programming techniques for HTTP clients and JSON handling

## Prerequisites

- Go 1.19 or higher
- Valid Twitch Client ID (obtain from Twitch Developer Console for research purposes)

## Building for Research

### Using Makefile (Recommended)
```bash
# Build both CLI and server versions
make build

# Build only CLI version
make build-cli

# Build only server version
make build-server

# Run tests
make test
```

### Manual Building
```bash
# Set your research Client ID
export TWITCH_CLIENT_ID="your_research_client_id"

# Build CLI version
go build -o twitchvod cmd/cli/main.go

# Build server version
go build -o twitchvod-server cmd/server/main.go
```

## Research Usage

### CLI Version
```bash
# Example research usage
./twitchvod https://www.twitch.tv/videos/1234567
```

### HTTP Server Version
```bash
# Set environment variables
export TWITCH_CLIENT_ID="your_research_client_id"
export SECRET="your_secret_key"

# Start the server
./twitchvod-server

# Make a request
curl "http://localhost:8080/extract?url=https://www.twitch.tv/videos/1234567&secret=your_secret_key"
```

Response format:
```json
{
  "success": true,
  "url": "https://fake-cdn.example.com/video/1234567/chunked/index-dvr.m3u8"
}
```

## Research Output

The tool outputs streaming metadata URLs for research analysis:

```
https://d2nvs31859zcd8.cloudfront.net/df8a614dfcddf26d99ca_zet_324493056124_1752771240/chunked/index-dvr.m3u8
```

## Research Features

- **URL Parsing Research**: Demonstrates regex-based video ID extraction
- **API Interaction Research**: Shows GraphQL API authentication and request patterns
- **Response Parsing Research**: Illustrates JSON response handling
- **Error Handling Research**: Demonstrates robust error management

## Educational Value

This code serves as a learning resource for:
- Go HTTP client implementation
- GraphQL API integration
- JSON marshaling/unmarshaling
- Regular expression usage
- Error handling patterns
- Environment variable configuration

## ⚖️ Legal Compliance

Users must:
- Obtain proper authorization from Twitch for API access
- Comply with Twitch's Developer Agreement
- Respect content creators' intellectual property rights
- Use only for legitimate research and educational purposes
- Not redistribute or modify Twitch content without permission

## 📚 Academic Use Only

This project is intended for:
- Computer science education
- API research and analysis
- Programming language learning
- Web scraping technique study
- **NOT for content downloading or redistribution**

---

**By using this code, you acknowledge that you understand and agree to use it only for legitimate research and educational purposes in compliance with all applicable laws and platform terms of service.** 