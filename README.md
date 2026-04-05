# Mezian — Plateforme de petites annonces

Mezian est une plateforme de petites annonces en ligne destinée au marché marocain. Elle permet aux particuliers et aux professionnels de publier, gérer et consulter des annonces, avec la possibilité de créer une boutique en ligne.

---

## Stack technique

| Couche | Technologie |
|--------|-------------|
| Backend | Go 1.22 + [Gin](https://github.com/gin-gonic/gin) |
| Frontend | Angular 19 (standalone components) |
| SDK | TypeScript (`@mezian/sdk`) |
| Base de données | SQLite + [GORM](https://gorm.io) |
| Authentification | JWT + OTP (WhatsApp / SMS) |
| Images | Upload + miniatures via [imaging](https://github.com/disintegration/imaging) |

---

## Structure du projet

```
mezian/
├── backend/    # API REST en Go (Gin)
├── front/      # Application Angular 19
└── sdk/        # SDK TypeScript (@mezian/sdk)
```

---

## Backend

### Lancer le serveur

```bash
cd backend
go run ./cmd/server/main.go
```

Le serveur démarre sur **http://localhost:8080**.

### Architecture

L'API suit une architecture en couches :

```
Handler → Service → Repository → Database (SQLite)
```

- **Handler** : traitement des requêtes HTTP (Gin)
- **Service** : logique métier (Auth, Ad, Media, Shop, Notification)
- **Repository** : accès aux données (GORM)
- **Middleware** : authentification JWT (`RequireAuth`, `OptionalAuth`)

### Endpoints de l'API (`/api/v1`)

#### Authentification
| Méthode | Route | Description |
|---------|-------|-------------|
| POST | `/auth/send-otp` | Envoyer un code OTP |
| POST | `/auth/verify-otp` | Vérifier le code OTP |
| POST | `/auth/register` | Créer un compte |
| POST | `/auth/login` | Se connecter |
| POST | `/auth/refresh` | Rafraîchir le token |
| POST | `/auth/logout` | Se déconnecter *(auth)* |
| GET | `/auth/me` | Profil courant *(auth)* |
| PUT | `/auth/me` | Modifier le profil *(auth)* |

#### Annonces
| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/ads` | Lister les annonces |
| GET | `/ads/:slug` | Détail d'une annonce |
| POST | `/ads` | Créer une annonce *(auth)* |
| PUT | `/ads/:slug` | Modifier une annonce *(auth)* |
| DELETE | `/ads/:slug` | Supprimer une annonce *(auth)* |
| POST | `/ads/:id/media` | Ajouter une image *(auth)* |
| POST | `/ads/:id/media/youtube` | Ajouter une vidéo YouTube *(auth)* |

#### Médias
| Méthode | Route | Description |
|---------|-------|-------------|
| DELETE | `/media/:id` | Supprimer un média *(auth)* |
| PUT | `/media/:id/cover` | Définir comme couverture *(auth)* |
| PUT | `/media/:id/order` | Modifier l'ordre *(auth)* |

#### Catégories
| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/categories` | Lister les catégories |
| GET | `/categories/:slug` | Détail d'une catégorie |

#### Boutiques
| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/shops/:slug` | Voir une boutique |
| POST | `/shops` | Créer une boutique *(auth)* |
| PUT | `/shops/:slug` | Modifier une boutique *(auth)* |

#### Utilisateur
| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/users/me/ads` | Mes annonces *(auth)* |
| GET | `/users/me/shop` | Ma boutique *(auth)* |

### Configuration

Le fichier `backend/config/config.yaml` centralise toute la configuration :

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

> **Important** : changer `jwt.secret` en production.

### Plans tarifaires

| Plan | Prix | Annonces | Durée |
|------|------|----------|-------|
| Starter | 199 MAD/mois | 50 | 30 jours |
| Pro | 399 MAD/mois | 200 | 30 jours |
| Premium | 799 MAD/mois | Illimité | 30 jours |

---

## Frontend

Application Angular 19 avec composants standalone et lazy loading.

### Lancer le frontend

```bash
cd front
npm install
npm start
```

L'application démarre sur **http://localhost:4200**.

### Build de production

```bash
npm run build
```

---

## SDK TypeScript

Le SDK `@mezian/sdk` permet d'interagir avec l'API depuis n'importe quelle application JavaScript/TypeScript.

### Build

```bash
cd sdk
npm install
npm run build
```

### Utilisation

```typescript
import { MezianClient } from '@mezian/sdk';

const client = new MezianClient({ baseURL: 'http://localhost:8080' });

// Authentification
await client.auth.sendOtp({ phone: '+212600000000' });
await client.auth.verifyOtp({ phone: '+212600000000', code: '123456' });

// Annonces
const ads = await client.ads.list();
```

---

## Prérequis

- **Go** >= 1.22
- **Node.js** >= 18
- **Angular CLI** >= 19

---

## Démarrage rapide

```bash
# Backend
cd backend && go run ./cmd/server/main.go

# Frontend (dans un autre terminal)
cd front && npm install && npm start
```
