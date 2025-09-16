# Playlist Migrator 🎵

An application to migrate playlists between music services (Spotify → Apple Music, etc).

## 🚀 Quick Setup

### 1. Prerequisites
- Go 1.25+
- Docker and Docker Compose
- Make (optional - there are alternatives for each OS)

### 2. Installation
```bash
# Clone the project
git clone https://github.com/zandomed/sync-playlist-api
cd sync-playlist-api

# Method 1: With Make (Unix/Linux/macOS/Windows with Make)
make deps
make install-tools

# Method 2: Platform-specific scripts
## Windows PowerShell
.\scripts\makefile.ps1 setup

# Configure environment variables
cp .env.example .env
# Edit .env with your credentials
```

### 3. Start services
```bash
# With Make
make docker-up
make migrate-up

# Or with platform-specific scripts
## Windows PowerShell
.\scripts\makefile.ps1 docker-up
.\scripts\makefile.ps1 migrate
```

### 4. Run the application
```bash
# With Make
make dev

# Or with platform-specific scripts
## Windows PowerShell
.\scripts\makefile.ps1 dev
```

The application will be available at `http://localhost:8080`

## 📡 API Endpoints

### Authentication
- `GET /api/v1/auth/spotify` - Start Spotify OAuth
- `GET /api/v1/auth/spotify/callback` - Spotify OAuth callback
- `GET /api/v1/auth/apple` - Start Apple Music OAuth
- `POST /api/v1/auth/refresh` - Refresh token

### Users (Authenticated)
- `GET /api/v1/users/me` - Get profile
- `PUT /api/v1/users/me` - Update profile

### Playlists (Authenticated)
- `GET /api/v1/playlists` - List playlists
- `GET /api/v1/playlists/:id` - Get specific playlist

### Migrations (Authenticated)
- `POST /api/v1/migrations` - Start migration
- `GET /api/v1/migrations` - List user migrations
- `GET /api/v1/migrations/:id` - Migration status
- `GET /api/v1/migrations/:id/progress` - Detailed progress
- `DELETE /api/v1/migrations/:id` - Cancel migration

### WebSocket
- `WS /ws/migration/:id` - Real-time progress

## 🛠️ Development Commands

### With Make (Recommended - cross-platform)
```bash
# Development
make dev          # Run with live reload
make run          # Run once
make build        # Compile binary

# Database
make docker-up    # Start PostgreSQL and Redis
make docker-down  # Stop services
make migrate-up   # Apply migrations
make migrate-down # Rollback migration
make migrate-status # View status

# Testing and quality
make test         # Run tests
make test-coverage # Tests with coverage
make lint         # Linter
make format       # Format code
make check-env    # Check installed tools

# Complete Docker
make docker-full  # Start everything including the app
make docker-clean # Clean volumes
```

### Alternative platform-specific scripts

#### Windows PowerShell
```powershell
.\scripts\makefile.ps1 help      # View available commands
.\scripts\makefile.ps1 setup     # Initial setup
.\scripts\makefile.ps1 dev       # Development mode
.\scripts\makefile.ps1 build     # Compile
.\scripts\makefile.ps1 test      # Run tests
.\scripts\makefile.ps1 check     # Check tools
```

> 💡 **Note**: If you have Make installed, use it preferentially as it works the same on all operating systems.

## 📁 Project Structure

```
sync-playlist/
├── cmd/server/           # Entry point
├── internal/
│   ├── config/          # Configuration
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   ├── models/          # Data models
│   ├── repository/      # Data access
│   └── services/        # Business logic
├── pkg/
│   ├── database/        # DB connection
│   └── logger/          # Logging system
├── migrations/          # SQL migrations
├── scripts/            # Utility scripts
└── docs/               # Documentation
```

## 🔧 Service Configuration

### Spotify
1. Create app in [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Configure redirect URI: `http://localhost:8080/auth/spotify/callback`
3. Add credentials to `.env`:
```bash
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
```

### Apple Music
1. Create certificate in [Apple Developer Portal](https://developer.apple.com/)
2. Configure MusicKit
3. Add credentials to `.env`:
```bash
APPLE_TEAM_ID=your_team_id
APPLE_KEY_ID=your_key_id
APPLE_PRIVATE_KEY="-----BEGIN PRIVATE KEY-----..."
```

## 🧪 Testing

```bash
# Unit tests
make test

# Coverage
make test-coverage
open coverage.html

# Integration tests
go test ./internal/... -tags=integration
```

## 🚀 Deployment

### With Docker
```bash
# Build image
docker build -t sync-playlist .

# Run container
docker run -p 9000:8080 --env-file .env sync-playlist
```

### Production Environment Variables
```bash
# Server
PORT=9000
HOST=0.0.0.0
LOG_LEVEL=INFO

# Database
DB_HOST=your-postgres-host
DB_PASSWORD=secure-password

# JWT
JWT_SECRET=your-super-secure-secret

# OAuth credentials
SPOTIFY_CLIENT_ID=prod-client-id
SPOTIFY_CLIENT_SECRET=prod-client-secret
# etc...
```

## 📚 Upcoming Features

- [ ] Complete OAuth for Spotify and Apple Music
- [ ] Improved track matching algorithm
- [ ] Queue system with Redis
- [ ] YouTube Music support
- [ ] Smart rate limiting
- [ ] Metrics and monitoring
- [ ] Integration tests
- [ ] CI/CD pipeline

## 🤝 Contributing

1. Clone the project
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push branch (`git push origin feature/new-feature`)
5. Create Pull Request

## 📄 License

This project is under the MIT License - see [LICENSE](LICENSE) for details.