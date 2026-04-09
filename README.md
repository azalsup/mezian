# Daba — Classified Ads Platform

Daba is an online classified ads platform for the Moroccan market. It allows individuals and professionals to publish, manage, and browse ads, with the option to create an online shop.

---

## Technical stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22 + [Gin](https://github.com/gin-gonic/gin) |
| Frontend | Angular 19 (standalone components) + TypeScript SDK |
| Database | SQLite + [GORM](https://gorm.io) |
| Authentication | JWT + OTP (WhatsApp / SMS) |
| Images | Upload + thumbnails via [imaging](https://github.com/disintegration/imaging) |

---

## Project structure

```
daba/
├── backend/    # Go REST API (Gin)
├── front/      # Angular 19 application with integrated TypeScript SDK
├── build.sh    # Production build script
├── serve_backend.sh  # Run backend in dev
├── serve_front.sh    # Run frontend in dev
└── start.sh    # Run production build
```

---

## Development Setup

### Prerequisites

- Go 1.22+
- Node.js 18+
- npm

### Backend

1. Navigate to backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the server:
   ```bash
   go run ./cmd/server/main.go
   ```

The server starts on **http://localhost:8080**.

### Frontend

1. Navigate to frontend directory:
   ```bash
   cd front
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Run the development server:
   ```bash
   ng serve
   ```

The frontend starts on **http://localhost:4200** (proxied to backend on `/api`).

### Running Both in Development

Use the provided scripts:

- `./serve_backend.sh` - Start backend
- `./serve_front.sh` - Start frontend

Or run both manually in separate terminals.

---

## Production Build

### Build Script

Run the build script to create production artifacts:

```bash
./build.sh
```

This will:
- Build the frontend with production API URL
- Generate static files in `dist/front/browser`

### Run Production

Use the start script to run the production build:

```bash
./start.sh
```

This starts:
- Backend on production config
- Frontend static files served on http://localhost:4200

### Manual Production Run

1. Build backend binary:
   ```bash
   cd backend
   go build -o ../bin/server ./cmd/server/main.go
   ```

2. Run backend:
   ```bash
   APP_ENV=production ./bin/server
   ```

3. Serve frontend:
   ```bash
   npx serve dist/front/browser -s -l 4200
   ```

---

## Backend Architecture

The API follows a layered architecture:

```
Handler → Service → Repository → Database (SQLite)
```

- **Handler**: HTTP request processing (Gin)
- **Service**: Business logic
- **Repository**: Data access layer
- **Database**: SQLite with GORM

### API Endpoints

- `GET /api/v1/ads` - List ads
- `POST /api/v1/ads` - Create ad
- `GET /api/v1/categories` - List categories
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register

### Configuration

- `config/config.yaml` - Development config
- `config/config.production.yaml` - Production config

---

## Frontend Architecture

- **Standalone Components**: Angular 19 with standalone components
- **Routing**: Lazy-loaded modules for admin and features
- **State Management**: Signals for reactive state
- **Styling**: Tailwind CSS + SCSS

### Key Features

- User authentication with OTP
- Ad creation and management
- Admin panel for user and role management
- Responsive design

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes
4. Run tests
5. Submit a pull request

---

## License

This project is licensed under the MIT License.
- **Service**: business logic (Auth, Ad, Media, Shop, Notification)
- **Repository**: data access (GORM)
- **Middleware**: JWT authentication (`RequireAuth`, `OptionalAuth`)

### API endpoints (`/api/v1`)

#### Authentication
| Method | Route | Description |
|--------|-------|-------------|
| POST | `/auth/send-otp` | Send an OTP code |
| POST | `/auth/verify-otp` | Verify an OTP code |
| POST | `/auth/register` | Create an account |
| POST | `/auth/login` | Log in |
| POST | `/auth/refresh` | Refresh the token |
| POST | `/auth/logout` | Log out *(auth)* |
| GET | `/auth/me` | Current profile *(auth)* |
| PUT | `/auth/me` | Update profile *(auth)* |

#### Ads
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/ads` | List ads |
| GET | `/ads/:slug` | Ad details |
| POST | `/ads` | Create an ad *(auth)* |
| PUT | `/ads/:slug` | Update an ad *(auth)* |
| DELETE | `/ads/:slug` | Delete an ad *(auth)* |
| POST | `/ads/:id/media` | Add an image *(auth)* |
| POST | `/ads/:id/media/youtube` | Add a YouTube video *(auth)* |

#### Media
| Method | Route | Description |
|--------|-------|-------------|
| DELETE | `/media/:id` | Delete media *(auth)* |
| PUT | `/media/:id/cover` | Set cover media *(auth)* |
| PUT | `/media/:id/order` | Update media order *(auth)* |

#### Categories
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/categories` | List categories |
| GET | `/categories/:slug` | Category details |

#### Shops
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/shops/:slug` | View a shop |
| POST | `/shops` | Create a shop *(auth)* |
| PUT | `/shops/:slug` | Update a shop *(auth)* |

#### User
| Method | Route | Description |
|--------|-------|-------------|
| GET | `/users/me/ads` | My ads *(auth)* |
| GET | `/users/me/shop` | My shop *(auth)* |

### Configuration

The file `backend/config/config.yaml` centralizes all configuration:

```yaml
server:
  port: 8080
  mode: debug          # debug | release

jwt:
  access_ttl_minutes: 15
  refresh_ttl_days: 30

otp:
  length: 6
  ttl_minutes: 10
  max_attempts: 5
  rate_limit_per_hour: 3

notification:
  provider: mock       # mock | twilio | infobip | orange_ma

media:
  max_size_mb: 5
  max_per_ad: 10
  thumbnail_width: 400
  thumbnail_height: 300
  allowed_types: [image/jpeg, image/png, image/webp]
```

> **Important**: change `jwt.secret` in production.

### Pricing plans

| Plan | Price | Ads | Duration |
|------|-------|------|----------|
| Starter | 199 MAD/month | 50 | 30 days |
| Pro | 399 MAD/month | 200 | 30 days |
| Premium | 799 MAD/month | Unlimited | 30 days |

---

## Frontend

Angular 19 application with standalone components and lazy loading.

### Run the frontend

```bash
cd front
npm install
ng build --configuration=production --base-href=/ --output-mode static
```

The app starts on **http://localhost:4200**.

### Production build

```bash
npm run build
```

---

## Requirements

- **Go** >= 1.22
- **Node.js** >= 18
- **Angular CLI** >= 19

---

## Quick start

```bash
# Backend
cd backend && go run ./cmd/server/main.go

# Frontend (in another terminal)
cd front && npm install && npm start
```
