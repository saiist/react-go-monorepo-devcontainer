# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a React-Go monorepo using contract-first API development with OpenAPI 3.0. The frontend uses React 19 with Vite and TypeScript, while the backend uses Go 1.23 with Chi router and GORM.

## Essential Commands

### Development
```bash
# Start all services (frontend + backend + databases)
pnpm dev

# Start specific services
pnpm backend:dev      # Backend only
pnpm --filter frontend dev  # Frontend only

# Docker services
pnpm docker:up       # Start PostgreSQL and Redis
pnpm docker:down     # Stop services
```

### Code Generation (Critical for API Contract)
```bash
# Generate both frontend and backend types from OpenAPI spec
pnpm generate

# This runs:
# - Frontend: Creates TypeScript client in apps/frontend/src/api/generated/
# - Backend: Creates Go types in apps/backend/internal/api/generated.go
```

### Database Operations
```bash
# Run migrations
pnpm db:migrate

# Create new migration
pnpm db:migrate:create NAME=migration_name

# Rollback last migration
pnpm db:migrate:down
```

### Testing and Quality
```bash
# Run all tests
pnpm test

# Run specific app tests
pnpm --filter frontend test
pnpm --filter backend test

# Linting
pnpm lint

# Format code
pnpm format
```

## Architecture Decisions

### API Contract-First Development
The OpenAPI specification at `/api/openapi.yaml` is the single source of truth. Always update this file first when changing API contracts, then run `pnpm generate` to update both frontend and backend code.

### Frontend Architecture
- **State Management**: TanStack Query for server state, Zustand for client state
- **API Client**: Auto-generated from OpenAPI using @hey-api/openapi-ts
- **Routing**: React Router v7 with type-safe routes
- **Styling**: TailwindCSS v4 with PostCSS

### Backend Architecture
- **HTTP Router**: Chi v5 with middleware chain
- **Database**: GORM with PostgreSQL, migrations using golang-migrate
- **API Generation**: oapi-codegen generates Chi server interfaces
- **Configuration**: Environment-based using internal/config package

### Authentication Flow
JWT-based authentication with refresh tokens:
1. Login endpoint returns access token (24h) and refresh token (7d)
2. Frontend stores tokens and includes bearer token in API requests
3. Backend validates JWT using middleware

## Development Environment

### Required Environment Variables
The `.env` file must include:
- `DATABASE_URL` - PostgreSQL connection with `?sslmode=disable` for local dev
- `JWT_SECRET` - For token signing
- `VITE_API_URL` - Frontend API base URL (default: http://localhost:8080/api/v1)

### Service Ports
- Frontend: 3000
- Backend: 8080
- PostgreSQL: 5432
- Redis: 6379

## Common Gotchas

1. **Always regenerate types** after modifying `/api/openapi.yaml`
2. **Frontend proxy**: Vite proxies `/api` requests to backend (configured in vite.config.ts)
3. **Database migrations**: Run after any schema changes in `/apps/backend/migrations/`
4. **Air configuration**: Backend uses Air for hot reload - check `.air.toml` for watched directories
5. **Turbo caching**: Use `turbo run dev --force` if you encounter stale cache issues

## Testing Approach

### Frontend Testing
- Component tests use Vitest + React Testing Library
- Test files co-located with components as `*.test.tsx`
- Mock API responses using MSW (Mock Service Worker)

### Backend Testing
- Table-driven tests preferred
- Test files as `*_test.go` in same package
- Use testify for assertions
- Database tests use test transactions

## Build and Deployment

### Local Build
```bash
pnpm build  # Builds all apps
```

### Docker Build
Both apps have production Dockerfiles:
- Frontend: Multi-stage build with nginx
- Backend: Multi-stage build with minimal Alpine image

## Security Considerations

- All API endpoints except `/health` and `/auth/*` require authentication
- CORS configured for frontend URL only
- Rate limiting enabled (100 requests/minute per IP)
- Input validation using OpenAPI schema constraints