# Cotton Cloud Backend

A Golang backend API for the Cotton Cloud digital wardrobe iOS app.

## Features

- **Clothing Management** - CRUD operations for wardrobe items
- **Avatar System** - Digital twin creation and management
- **Outfit Logging** - Daily outfit journaling with calendar
- **AI Integration** - Proxy endpoints for Gemini AI features
- **Wear Tracking** - Laundry reminders based on wear count

## Tech Stack

- **Framework**: [Gin](https://gin-gonic.com/) - HTTP web framework
- **Database**: SQLite with [GORM](https://gorm.io/) ORM
- **Auth**: JWT tokens (to be implemented)
- **AI**: Google Gemini API proxy

## Getting Started

### Prerequisites

- Go 1.21 or later
- SQLite

### Installation

1. Clone and navigate to the project:
```bash
cd cotton-cloud-backand
```

2. Install dependencies:
```bash
go mod download
```

3. Copy environment file:
```bash
cp .env.example .env
```

4. Set your Gemini API key in `.env`:
```
GEMINI_API_KEY=your_api_key_here
```

5. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

### Health Check
- `GET /health` - Server health status

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Clothing
- `GET /api/v1/clothing` - List all items
- `POST /api/v1/clothing` - Create item
- `GET /api/v1/clothing/:id` - Get item
- `PUT /api/v1/clothing/:id` - Update item
- `DELETE /api/v1/clothing/:id` - Delete item
- `POST /api/v1/clothing/:id/wash` - Mark as washed
- `POST /api/v1/clothing/:id/wear` - Increment wear count

### Avatars
- `GET /api/v1/avatars` - List all avatars
- `POST /api/v1/avatars` - Create avatar
- `GET /api/v1/avatars/:id` - Get avatar
- `PUT /api/v1/avatars/:id` - Update avatar
- `DELETE /api/v1/avatars/:id` - Delete avatar
- `POST /api/v1/avatars/:id/activate` - Set as active

### Outfits
- `GET /api/v1/outfits` - List all records
- `POST /api/v1/outfits` - Log outfit
- `GET /api/v1/outfits/:date` - Get by date
- `PUT /api/v1/outfits/:id` - Update record
- `DELETE /api/v1/outfits/:id` - Delete record

### AI (Gemini Proxy)
- `POST /api/v1/ai/analyze` - Analyze clothing image
- `POST /api/v1/ai/cutout` - Generate cutout
- `POST /api/v1/ai/avatar` - Generate avatar
- `POST /api/v1/ai/collage` - Generate collage
- `POST /api/v1/ai/tryon` - Virtual try-on

## Project Structure

```
cotton-cloud-backand/
├── cmd/server/main.go         # Entry point
├── internal/
│   ├── api/
│   │   ├── router.go          # Route definitions
│   │   └── handlers/          # Request handlers
│   ├── models/                # Database models
│   ├── services/              # Business logic
│   └── database/              # DB configuration
├── configs/                   # Config files
├── go.mod
└── README.md
```

## License

MIT
