package twitch

// TwitchGraphQLRequest represents the structure of the GraphQL request
type TwitchGraphQLRequest []GraphQLOperation

// GraphQLOperation represents a single GraphQL operation
type GraphQLOperation struct {
	OperationName string            `json:"operationName"`
	Variables     GraphQLVariables  `json:"variables"`
	Extensions    GraphQLExtensions `json:"extensions"`
}

// GraphQLVariables represents the variables for a GraphQL operation
type GraphQLVariables struct {
	VideoID string `json:"videoID,omitempty"`
	VodID   string `json:"vodID,omitempty"`
}

// GraphQLExtensions represents the extensions for a GraphQL operation
type GraphQLExtensions struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}

// PersistedQuery represents a persisted query extension
type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}

// TwitchGraphQLResponse represents the structure of the GraphQL response
type TwitchGraphQLResponse []GraphQLResponseItem

// GraphQLResponseItem represents a single GraphQL response item
type GraphQLResponseItem struct {
	Data GraphQLData `json:"data"`
}

// GraphQLData represents the data field in a GraphQL response
type GraphQLData struct {
	Video GraphQLVideo `json:"video"`
}

// GraphQLVideo represents video data in a GraphQL response
type GraphQLVideo struct {
	SeekPreviewsURL string `json:"seekPreviewsURL"`
}
