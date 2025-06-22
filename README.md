# wifi-golang-backend

A backend server written in Go for managing WiFi network data, user authentication, statistics, and AI-powered travel recommendations, with a focus on secure OAuth flows and location-aware WiFi access.

---

## Features

- **OAuth Authentication**: Implements Civic OAuth2 with PKCE for secure user authentication.
- **WiFi Management**:
  - Scan and add new WiFi networks (with location and description).
  - Connect to WiFi networks with distance validation (ensures user is near the network).
  - List nearby WiFi networks based on geolocation.
  - Fetch all available and saved WiFi networks.
- **AI-Powered Recommendations**:
  - Integrates with Google Gemini AI to recommend interesting stops (landmarks, attractions, etc.) between two locations.
  - For each recommended stop, fetches all available WiFi networks nearby.
  - Enables users to plan routes and discover both places and connectivity along the way.
- **Statistics**: Endpoints to get and update user/network statistics.
- **MongoDB Integration**: Uses MongoDB for persistent storage with geospatial queries.
- **Modular & Testable**: Handlers are structured for dependency injection and easy testing.

---

## Project Structure

```
cmd/server/main.go         # Application entrypoint
config/                    # Configuration loading (env, secrets)
internal/auth/             # OAuth logic and authentication middleware
internal/db/               # MongoDB connection and collections
internal/models/           # Data models (WiFi, Location, etc.)
internal/routes/           # HTTP route handlers (including AI recommendation)
internal/utils/            # Utility functions
```

---

## Getting Started

### Prerequisites

- Go 1.18+
- MongoDB (local or Atlas)
- Civic OAuth credentials (Client ID, Secret, Redirect URI)
- Google Gemini API Key (for recommendations)

### Environment Variables

Create a `.env` file in the project root with the following:

```
MONGO_URI=mongodb+srv://<user>:<pass>@cluster0.mongodb.net/wifi_db?retryWrites=true&w=majority
OAUTH_CLIENT_ID=your_civic_client_id
OAUTH_CLIENT_SECRET=your_civic_client_secret
GEMINI_API_KEY=your_gemini_api_key
```

### Installation and Running

```bash
go mod download
go run cmd/server/main.go
```

Server will start on `:8080` by default.

---

## API Overview

### Authentication Endpoints

- `GET /api/auth/me` — Get current user info (requires auth)
- `POST /api/auth/upgrade` — Upgrade verification level (requires auth)
- `GET /api/auth/civic` — Begin Civic OAuth2 flow
- `GET /api/auth/civic/callback` — OAuth2 callback handler

### WiFi Endpoints

- `POST /api/wifi/scan` — Add new WiFi network (requires auth)
- `POST /api/wifi/connect` — Connect to WiFi (requires auth, location-based)
- `GET /api/wifi/nearby` — List nearby networks (latitude/longitude required)
- `GET /api/wifi/all` — List all WiFi networks
- `GET /api/wifi/saved` — List saved WiFi networks (requires auth)
- `POST /api/wifi/nearby/stops` — Given a list of stops, returns all WiFi networks near each stop

### AI Recommendation Endpoints

- `GET /api/gemini/recommendstops`  
  Returns 5 recommended stops (landmarks, attractions, etc.) between two coordinates using Gemini AI.  
  **Query parameters:** `start_lat`, `start_lng`, `end_lat`, `end_lng`

- `GET /api/gemini/recommendstopswifi`  
  Returns 5 recommended stops between two coordinates, and for each stop, lists all nearby WiFi networks.  
  **Query parameters:** `start_lat`, `start_lng`, `end_lat`, `end_lng`

---

## Development Notes

- Handlers return robust validation errors for malformed requests or missing data.
- Auth middleware is designed for upgradeability (currently checks a test token, easily extended for real JWT/OAuth).
- Geospatial queries and distance checks use MongoDB’s `$geoWithin` and the Haversine formula.
- **Mobile/remote DB connection:** When connecting from a mobile device, ensure your public IP is whitelisted in your MongoDB instance. Avoid `0.0.0.0/0` in production.
- **Gemini AI Integration:**  
  The backend uses Google Gemini AI for intelligent, real-world stop recommendations along a route. This enables users to discover both interesting places and available WiFi networks for a seamless travel experience.

---

## Troubleshooting

- **MongoDB connection from mobile:** Check your public IP (`whatismyip.com`) and whitelist it in MongoDB Atlas. Be aware of carrier NAT and dynamic IPs.
- **OAuth redirect issues:** The redirect URI in Civic app settings must exactly match the backend route.
- **Port errors:** If deploying in the cloud, check your provider’s port and firewall settings.
- **Gemini API issues:** Ensure your `GEMINI_API_KEY` is set and valid. Check logs for AI response or parsing errors.

---

## License

MIT

---

## Contributing

Feel free to open issues or PRs for bugs, feature requests, or improvements!

---