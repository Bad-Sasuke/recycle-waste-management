# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Full-stack Recycle Waste Management System with:
- **Frontend**: Vue.js 3 + TypeScript (in `frontend/`)
- **Backend**: Go Fiber (in `backend/`)
- **Database**: MongoDB
- **Cache**: Redis
- **Storage**: AWS S3
- **Auth**: Firebase + GitHub OAuth + JWT

## Development Commands

### Frontend (`cd frontend`)

```bash
# Development
npm run dev              # Start dev server (Vite)
npm run build           # Build for production
npm run preview         # Preview production build

# Code Quality
npm run check           # Run type-check and lint together
npm run type-check      # TypeScript compilation check
npm run lint            # ESLint with auto-fix
npm run format          # Prettier formatting

# Testing
npm run test:unit       # Run Vitest tests

# Storybook
npm run storybook       # Start Storybook dev server on port 6006
npm run build-storybook # Build Storybook for production
```

### Backend (`cd backend`)

```bash
# Development
go run main.go          # Start server (reads .env.local)
go mod tidy            # Install/update dependencies

# The server runs on PORT env var (default: 8080)
# API docs at: http://localhost:{PORT}/api/docs (Scalar UI)
# Monitor at: http://localhost:{PORT}/api/monitor
```

## Architecture

### Backend Structure (Go)

The backend follows **clean architecture** with dependency injection:

```
backend/src/
├── domain/
│   ├── entities/      # Data structures (User, Shop, RecyclableItem, etc.)
│   └── datasources/   # Database connection (MongoDB)
├── repositories/      # Data access layer (MongoDB operations)
├── services/          # Business logic layer
├── gateways/          # HTTP handlers + route definitions
├── middlewares/       # JWT auth, logging, monitoring
├── websocket/         # WebSocket hub for chat
├── configuration/     # Fiber app config
└── infrastructure/    # External services (AWS S3, Redis)
```

**Dependency flow**: `main.go` → repositories → services → gateways → routes

**Key files**:
- [main.go](backend/main.go): Initializes dependencies, connects MongoDB, sets up routes
- [gateways/route.go](backend/src/gateways/route.go): Route definitions with middleware
- [domain/entities/](backend/src/domain/entities/): Core data structures

### Frontend Structure (Vue 3 + TypeScript)

```
frontend/src/
├── views/           # Page components (HomeView, MarketplaceView, etc.)
├── components/      # Reusable UI components
├── stores/          # Pinia stores (users, shop, wastes, etc.)
├── services/        # API calls and Firebase integration
├── router/          # Vue Router with auth guards
├── i18n/            # Translations (en.json, th.json)
├── types/           # TypeScript interfaces
└── config.ts        # Environment variable mapping
```

**State Management (Pinia)**:
- [stores/users.ts](frontend/src/stores/users.ts): User auth, profile, role, employee status
- [stores/shop.ts](frontend/src/stores/shop.ts): Shop data and ownership check
- [stores/wastes.ts](frontend/src/stores/wastes.ts): Recyclable waste items
- All stores use `pinia-plugin-persistedstate` for persistence

**Routing Guards** ([router/index.ts](frontend/src/router/index.ts:84-156)):
- Checks authentication via `useUsersStore().checkLogin()`
- Moderators without shops are redirected to `/shop/create`
- Route meta fields: `requiresAuth`, `restrictedRoles`

### User Roles & Permissions

Three roles defined in [backend/src/domain/entities/users.go](backend/src/domain/entities/users.go:6-12):
- `user`: Regular customers (can request waste pickup)
- `moderator`: Shop owners (manage shops, employees, purchases)
- `admin`: Full system access

**Shop-Moderator Relationship**:
- Moderators must create a shop before accessing most features
- Shop creation flow enforced in router guard
- `useShopStore().checkUserShop()` validates shop ownership
- Employees can act on behalf of shop owners

### API Communication

**Frontend → Backend**:
- Base URL: `import.meta.env.VITE_WEB_API` (from [config.ts](frontend/src/config.ts:2))
- Service files in `frontend/src/services/` handle HTTP requests
- JWT token stored in cookies, sent via `Authorization` header

**Authentication Flow**:
1. User signs in via Firebase (Google) or GitHub OAuth
2. Backend validates and issues JWT
3. JWT stored in cookie by [stores/cookie.ts](frontend/src/stores/cookie.ts)
4. Protected routes require JWT via [middlewares/jwtmid.go](backend/src/middlewares/jwtmid.go)

### WebSocket (Chat System)

- Server: [websocket/hub.go](backend/src/websocket/hub.go) manages connections
- Endpoint: `ws://localhost:{PORT}/ws/chat`
- Frontend store: [stores/chat.ts](frontend/src/stores/chat.ts)

## Environment Setup

### Backend `.env.local` (required variables):

```bash
DATABASE_NAME=your_db_name
MONGODB_URI=mongodb://localhost:27017
PORT=1818
JWT_SECRET_KEY=your_secret
JWT_REFESH_SECRET_KEY=your_refresh_secret

# AWS S3 for image uploads
AWS_BUCKET_NAME=your_bucket
AWS_REGION=ap-southeast-1
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret

# Redis (optional - for caching)
REDIS_URI=redis://localhost:6379
REDISREAD_URI=redis://localhost:6379

# OAuth
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_secret
GITHUB_REDIRECT_URI=http://localhost:1818/api/auth/github/callback

# Firebase service account (for token validation)
FIREBASE_SERVICE_ACCOUNT='{"type": "service_account", ...}'

FRONTEND_URL=http://localhost:5173
VAT_RATE=0.07
```

### Frontend `.env` (required variables):

```bash
VITE_WEB_API=http://127.0.0.1:1818

# GitHub OAuth
VITE_GITHUB_CLIENT_ID=your_client_id

# Firebase config (from Firebase console)
VITE_API_KEY=your_api_key
VITE_AUTH_DOMAIN=your_domain
VITE_PROJECT_ID=your_project_id
VITE_STORAGE_BUCKET=your_bucket
VITE_MESSAGING_SENDER_ID=your_sender_id
VITE_APP_ID=your_app_id
VITE_MEASUREMENT_ID=your_measurement_id

VITE_VAT_RATE=0.07
```

## Key Features & Implementation

### Shop Management
- Shop creation: [frontend/src/views/CreateShopView.vue](frontend/src/views/CreateShopView.vue)
- Shop service: [backend/src/services/shops.go](backend/src/services/shops.go)
- Employees: [backend/src/services/employee.go](backend/src/services/employee.go)

### Waste Purchase Flow
1. Customer creates request via [services/customer_request.ts](frontend/src/services/customer_request.ts)
2. Shop accepts request
3. Shop creates receipt via [backend/src/services/receipt.go](backend/src/services/receipt.go)
4. Stock updated via [backend/src/services/stock.go](backend/src/services/stock.go)

### Map Integration
- Library: Leaflet ([package.json](frontend/package.json:27))
- Components: [components/map/](frontend/src/components/map/)
- Shows shops and customer locations

### Component Documentation
- Storybook stories in `frontend/src/components/**/*.stories.ts`
- Run `npm run storybook` to view component library

## Common Patterns

### Adding a New API Endpoint (Backend)

1. Define entity in `domain/entities/`
2. Create repository in `repositories/` (implements interface)
3. Create service in `services/` (business logic)
4. Create gateway handler in `gateways/`
5. Add route in `gateways/route.go` or create new route function
6. Initialize in [main.go](backend/main.go)

**Example**: See [gateways/employee.go](backend/src/gateways/employee.go) + [services/employee.go](backend/src/services/employee.go)

### Adding a New Frontend Page

1. Create view component in `views/`
2. Add route in [router/index.ts](frontend/src/router/index.ts)
3. If requiring auth, add `meta: { requiresAuth: true }`
4. If role-restricted, add `meta: { restrictedRoles: ['role'] }`
5. Create Pinia store if state needed in `stores/`
6. Create API service in `services/`

### Path Alias
Frontend uses `@` alias for `src/` directory (configured in [vite.config.ts](frontend/vite.config.ts:16-18))

```typescript
import { useUsersStore } from '@/stores/users'
```

## TypeScript & Type Safety

- Frontend: Strict TypeScript enabled
- Type definitions in `frontend/src/types/`
- Interfaces for API responses: [types/response.ts](frontend/src/types/response.ts)
- Backend: Go's static typing, entities in `domain/entities/`

## Testing

- Frontend: Vitest for unit tests
- Component testing: `@vue/test-utils` available
- Run tests: `npm run test:unit` (in `frontend/`)
- Storybook for visual testing

## API Documentation

Access Scalar API docs when backend is running:
- URL: `http://localhost:{PORT}/api/docs`
- Spec: [backend/src/docs/swagger.yaml](backend/src/docs/swagger.yaml)
