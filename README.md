# Daba — Classified Ads Platform

Daba is an online classified ads platform for the Moroccan market. It allows individuals and professionals to publish, manage, and browse ads, with the option to create an online shop.

---

## Technical stack

| Layer | Technology |
|-------|------------|
| Backend | Go 1.22 + [Gin](https://github.com/gin-gonic/gin) |
| Frontend | Angular 19 (standalone components) |
| SDK | TypeScript (`@classifieds/sdk`) |
| Database | SQLite + [GORM](https://gorm.io) |
| Authentication | JWT + OTP (WhatsApp / SMS) |
| Images | Upload + thumbnails via [imaging](https://github.com/disintegration/imaging) |

---

## Project structure

```
daba/
├── backend/    # Go REST API (Gin)
├── front/      # Angular 19 application
└── sdk/        # TypeScript SDK (@daba/sdk)
```

---

## Backend

### Run the server

```bash
cd backend
go run ./cmd/server/main.go
```

The server starts on **http://localhost:8080**.

### Architecture

The API follows a layered architecture:

```
Handler → Service → Repository → Database (SQLite)
```

- **Handler**: HTTP request processing (Gin)
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
npm start
```

The app starts on **http://localhost:4200**.

### Production build

```bash
npm run build
```

---

## TypeScript SDK

The `@daba/sdk` SDK allows you to interact with the API from any JavaScript/TypeScript app.

### Build

```bash
cd sdk
npm install
npm run build
```

### Usage

```typescript
import { ClassifiedsClient } from '@classifieds/sdk';

const client = new ClassifiedsClient({ baseURL: 'http://localhost:8080' });

// Authentication
await client.auth.sendOtp({ phone: '+212600000000' });
await client.auth.verifyOtp({ phone: '+212600000000', code: '123456' });

// Ads
const ads = await client.ads.list();
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
