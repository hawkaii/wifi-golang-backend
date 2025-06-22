package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Coordinate represents a geographical coordinate
type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
}

// RecommendationRequest represents the input for finding stops
type RecommendationRequest struct {
	StartCoordinate Coordinate `json:"start_coordinate"`
	EndCoordinate   Coordinate `json:"end_coordinate"`
	StopType        string     `json:"stop_type"` // e.g., "restaurants", "gas_stations", "tourist_attractions"
	MaxStops        int        `json:"max_stops"`
}

// RecommendationResponse represents the response with recommended stops
type RecommendationResponse struct {
	Stops []Coordinate `json:"stops"`
	Route string       `json:"route_description"`
}

// LocationRecommender handles AI-powered location recommendations
type LocationRecommender struct {
	client *genai.Client
}

// NewLocationRecommender creates a new instance of LocationRecommender
func NewLocationRecommender(apiKey string) (*LocationRecommender, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &LocationRecommender{
		client: client,
	}, nil
}

// FindStops uses Gemini AI to find recommended stops between two coordinates
func (lr *LocationRecommender) FindStops(ctx context.Context, req RecommendationRequest) (*RecommendationResponse, error) {
	model := lr.client.GenerativeModel("gemini-1.5-flash")

	// Configure the model for structured output
	model.SetTemperature(0.7)

	// Create a detailed prompt for Gemini
	prompt := fmt.Sprintf(`You are a travel recommendation AI. Given two coordinates, find interesting stops along or near the route.

START COORDINATE: %f, %f
END COORDINATE: %f, %f
STOP TYPE: %s
MAX STOPS: %d

Please recommend stops of type "%s" between or near these coordinates. For each stop, provide:
1. The name of the location
2. The exact latitude and longitude coordinates
3. A brief reason why it's recommended

Return your response in the following JSON format only (no additional text):
{
  "stops": [
    {
      "latitude": 0.0,
      "longitude": 0.0,
      "name": "Location Name"
    }
  ],
  "route_description": "Brief description of the route and recommendations"
}

Important guidelines:
- Provide real, existing locations with accurate coordinates
- Consider the geographical path between start and end points
- Ensure coordinates are realistic for the region
- Limit to %d stops maximum
- Focus on popular, well-known locations of the requested type`,
		req.StartCoordinate.Latitude, req.StartCoordinate.Longitude,
		req.EndCoordinate.Latitude, req.EndCoordinate.Longitude,
		req.StopType, req.MaxStops, req.StopType, req.MaxStops)

	// Generate content using Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	// Extract the response text
	responseText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			responseText += string(textPart)
		}
	}

	// Clean the response to extract JSON
	responseText = strings.TrimSpace(responseText)

	// Find JSON content (remove any markdown formatting)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonContent := responseText[jsonStart:jsonEnd]

	// Parse the JSON response
	var recommendation RecommendationResponse
	if err := json.Unmarshal([]byte(jsonContent), &recommendation); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &recommendation, nil
}

// FindStopsBetween returns 5 recommended stops as JSON between two coordinates
func (lr *LocationRecommender) FindStopsBetween(ctx context.Context, startLat, startLng, endLat, endLng float64) ([]byte, error) {
	model := lr.client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.7)

	prompt := fmt.Sprintf(`You are a travel recommendation AI. Given two coordinates, find 5 interesting stops along or near the route.

START COORDINATE: %f, %f
END COORDINATE: %f, %f
STOP TYPE: any
MAX STOPS: 5

Please recommend 5 stops between or near these coordinates. For each stop, provide:
1. The name of the location
2. The exact latitude and longitude coordinates
3. A brief reason why it's recommended

Return your response in the following JSON format only (no additional text):
{
  "stops": [
    {
      "latitude": 0.0,
      "longitude": 0.0,
      "name": "Location Name"
    }
  ]
}

Important guidelines:
- Provide real, existing locations with accurate coordinates
- Consider the geographical path between start and end points
- Ensure coordinates are realistic for the region
- Limit to 5 stops maximum
- Focus on popular, well-known locations`,
		startLat, startLng, endLat, endLng)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response generated")
	}

	responseText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			responseText += string(textPart)
		}
	}
	responseText = strings.TrimSpace(responseText)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}") + 1
	if jsonStart == -1 || jsonEnd <= jsonStart {
		return nil, fmt.Errorf("no valid JSON found in response")
	}
	jsonContent := responseText[jsonStart:jsonEnd]
	return []byte(jsonContent), nil
}

// Close closes the client connection
func (lr *LocationRecommender) Close() {
	if lr.client != nil {
		lr.client.Close()
	}
}

// Helper function to parse coordinate from string
func parseCoordinate(coordStr string) (Coordinate, error) {
	parts := strings.Split(coordStr, ",")
	if len(parts) != 2 {
		return Coordinate{}, fmt.Errorf("invalid coordinate format, expected 'lat,lng'")
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid longitude: %w", err)
	}

	return Coordinate{Latitude: lat, Longitude: lng}, nil
}
