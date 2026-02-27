# 🚀 Dynamic Landing Page CMS

A production-ready, full-stack monorepo for a **dynamic landing page** with a complete CMS admin panel. Every piece of content on the landing page (title, description, hero text, features, testimonials, pricing, FAQs, navigation, OG meta, Twitter cards, etc.) is fully manageable through the admin panel.

---

## 📁 Project Structure

```
/
├── apps/
│   ├── backend/          # Go + Gin REST API (Clean Architecture)
│   └── frontend/         # Next.js 14 App Router
├── scripts/
│   └── migrations/       # 8 SQL migration files for Supabase
├── docker/
│   └── nginx/            # Nginx reverse proxy config
├── docker-compose.yml    # Development setup
├── docker-compose.prod.yml # Production setup
├── README.md
└── RULES.md
```

---

## 🏗️ Architecture

### Backend (Go + Gin) — Clean Architecture
```
Handler → Service → Repository → Database
```
- **Framework**: Gin (high-performance HTTP)
- **Database**: Supabase (PostgreSQL) via `pgx` + `sqlx`
- **Auth**: JWT (HS256) — access token (15min) + refresh token (7 days, httpOnly cookie)
- **Security**: RBAC, bcrypt (cost 12), rate limiting, CORS, security headers, account lockout
- **Logging**: zerolog (structured JSON)
- **Config**: viper
- **Validation**: go-playground/validator with custom validators
- **Testing**: Go standard testing + table-driven tests

### Frontend (Next.js 14) — App Router
- **Language**: TypeScript (strict mode)
- **Styling**: Tailwind CSS
- **State**: Zustand (client) + React Query (server)
- **Forms**: React Hook Form + Zod
- **Testing**: Vitest + Testing Library
- **Rendering**: SSR + ISR (60s revalidation) for landing page

---

## 🗄️ Database Schema (ERD)

### Tables

| Table | Description |
|-------|-------------|
| `users` | Admin users with roles (super_admin/admin/editor) |
| `refresh_tokens` | JWT refresh token store with revocation |
| `password_reset_tokens` | Password reset flow |
| `sites` | Multi-site support |
| `site_settings` | 30+ configurable settings (SEO, OG, social, analytics, appearance) |
| `pages` | Landing pages with full SEO/OG/Twitter meta |
| `page_sections` | Sections within pages (hero, features, pricing, etc.) |
| `section_contents` | Key-value content per section (flexible, extensible) |
| `media` | Uploaded files (Supabase Storage) |
| `navigation_menus` | Navigation menu containers |
| `navigation_items` | Hierarchical menu items |
| `features` | Product feature cards |
| `testimonials` | Customer reviews |
| `pricing_plans` | Pricing tiers with feature lists |
| `faqs` | FAQ items with categories |
| `cta_buttons` | Reusable CTA buttons |
| `social_links` | Social media links |
| `team_members` | Team member profiles |
| `logos` | Partner/client logos |
| `audit_logs` | Complete admin action audit trail |
| `schema_migrations` | Migration tracking |

---

## 🔐 Security

### Auth Flow
1. POST `/api/v1/auth/login` → returns access token (body) + refresh token (httpOnly cookie)
2. Access token stored in memory (Zustand), used for API calls
3. Refresh token in httpOnly cookie, used to silently refresh access token
4. On 401, axios interceptor auto-refreshes and retries
5. Logout revokes refresh token in DB

### Security Measures
- Passwords: bcrypt cost 12
- JWT: HS256, separate secrets for access/refresh
- Account lockout: 5 failed attempts → 15 min lock
- Token rotation: new refresh token on every refresh
- Rate limiting: 5 req/min (auth), 100 req/min (public), 200 req/min (admin)
- CORS whitelist
- Security headers: X-Frame-Options, X-Content-Type-Options, X-XSS-Protection, etc.
- Request ID for distributed tracing
- Audit logging for all admin actions

---

## 📡 API Reference

### Public Endpoints (no auth)
```
GET  /health                                    # Health check with DB status
GET  /api/v1/public/sites/:id                  # Site info + public settings
GET  /api/v1/public/site/:slug                 # Site by slug
GET  /api/v1/public/pages/:slug?site_id=...    # Page with sections + content (SSR)
GET  /api/v1/public/pages?site_id=...          # Homepage
GET  /api/v1/public/navigation/:siteId/:id     # Navigation menu tree
```

### Auth Endpoints (rate-limited: 5/min)
```
POST /api/v1/auth/login                        # Login → access token + refresh cookie
POST /api/v1/auth/logout                       # Logout → revoke refresh token
POST /api/v1/auth/refresh                      # Refresh access token
GET  /api/v1/auth/me                           # Current user (requires auth)
POST /api/v1/auth/change-password              # Change password (requires auth)
```

### Admin Endpoints (requires auth + role)

#### Sites (admin+)
```
GET    /api/v1/admin/sites
POST   /api/v1/admin/sites
GET    /api/v1/admin/sites/:id
PUT    /api/v1/admin/sites/:id
DELETE /api/v1/admin/sites/:id          # super_admin only
GET    /api/v1/admin/sites/:id/settings
PUT    /api/v1/admin/sites/:id/settings
PUT    /api/v1/admin/sites/:id/settings/:key
```

#### Pages (editor+)
```
GET    /api/v1/admin/pages
POST   /api/v1/admin/pages
GET    /api/v1/admin/pages/:id
PUT    /api/v1/admin/pages/:id
DELETE /api/v1/admin/pages/:id          # admin+
PATCH  /api/v1/admin/pages/:id/publish
PATCH  /api/v1/admin/pages/:id/unpublish
GET    /api/v1/admin/pages/:id/sections
POST   /api/v1/admin/pages/:id/sections
```

#### Sections & Content (editor+)
```
PUT    /api/v1/admin/sections/:id
DELETE /api/v1/admin/sections/:id
PATCH  /api/v1/admin/sections/reorder
GET    /api/v1/admin/sections/:id/contents
POST   /api/v1/admin/sections/:id/contents
POST   /api/v1/admin/sections/:id/contents/bulk
DELETE /api/v1/admin/contents/:id
```

#### Components (editor+)
```
GET/POST/PUT/DELETE /api/v1/admin/features
GET/POST/PUT/DELETE /api/v1/admin/testimonials
GET/POST/PUT/DELETE /api/v1/admin/pricing
GET/POST/PUT/DELETE /api/v1/admin/faqs
GET                 /api/v1/admin/navigation
POST                /api/v1/admin/navigation/:menuId/items
PUT/DELETE          /api/v1/admin/navigation/items/:id
GET/POST/PUT/DELETE /api/v1/admin/media
POST                /api/v1/admin/media/upload
```

#### Users & Audit (admin+)
```
GET/POST/PUT/DELETE /api/v1/admin/users
GET                 /api/v1/admin/audit-logs
```

---

## 🌐 Admin Panel Pages

| URL | Description |
|-----|-------------|
| `/admin/login` | Secure login |
| `/admin/dashboard` | Overview + quick actions |
| `/admin/pages` | Pages list with publish/unpublish |
| `/admin/pages/new` | Create page with SEO/OG fields |
| `/admin/pages/[id]` | **Page editor** — sections + inline content + SEO panel |
| `/admin/features` | Product features CRUD |
| `/admin/testimonials` | Customer reviews CRUD |
| `/admin/pricing` | Pricing plans CRUD |
| `/admin/faqs` | FAQ items CRUD |
| `/admin/navigation` | Navigation menus + hierarchical items |
| `/admin/media` | Media library with drag-and-drop upload |
| `/admin/settings` | Site settings (SEO, OG, social, analytics, appearance) |
| `/admin/users` | User management with RBAC |
| `/admin/audit-logs` | Filterable audit trail |

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Node.js 20+
- pnpm 9+
- Docker & Docker Compose (optional)
- Supabase account

### 1. Database Setup

Run all 8 migration files in order in your Supabase SQL editor:

```sql
-- Run in Supabase SQL Editor, in order:
-- 001_create_users.sql
-- 002_create_sites.sql
-- 003_create_pages.sql
-- 004_create_content.sql
-- 005_create_media.sql
-- 006_create_navigation.sql
-- 007_create_components.sql
-- 008_create_audit.sql
```

Or via psql:
```bash
for i in 001 002 003 004 005 006 007 008; do
  psql $DATABASE_URL -f scripts/migrations/${i}_*.sql
done
```

After running migrations, get your default site ID:
```sql
SELECT id FROM sites WHERE slug = 'default';
```

### 2. Backend Setup

```bash
cd apps/backend
cp .env.example .env
# Edit .env with your values:
# - DATABASE_URL (from Supabase: Settings → Database → Connection string)
# - JWT_ACCESS_SECRET (generate: openssl rand -base64 64)
# - JWT_REFRESH_SECRET (generate: openssl rand -base64 64)
# - SUPABASE_URL, SUPABASE_ANON_KEY, SUPABASE_SERVICE_KEY

# Install dependencies
go mod download

# Run with hot reload (install air first: go install github.com/air-verse/air@latest)
make dev

# Or run directly
make run
```

### 3. Frontend Setup

```bash
cd apps/frontend
cp .env.example .env.local
# Edit .env.local:
# - NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
# - NEXT_PUBLIC_DEFAULT_SITE_ID=<your-site-id-from-step-1>

# Install dependencies
pnpm install

# Run development server
pnpm dev
```

### 4. Access the App

- **Landing page**: http://localhost:3000
- **Admin panel**: http://localhost:3000/admin/login
- **API**: http://localhost:8080
- **Health check**: http://localhost:8080/health

**Default admin credentials:**
- Email: `admin@example.com`
- Password: `Admin@123456`
- ⚠️ **Change immediately in production!**

---

## 🐳 Docker Deployment

### Development
```bash
docker-compose up -d
```

### Production
```bash
# Copy and fill production env vars
cp apps/backend/.env.example apps/backend/.env.prod
cp apps/frontend/.env.example apps/frontend/.env.prod

docker-compose -f docker-compose.prod.yml up -d
```

---

## 🧪 Testing

### Backend Tests
```bash
cd apps/backend

# Run all tests
go test ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/service/...
go test ./internal/pkg/auth/...
```

### Frontend Tests
```bash
cd apps/frontend

# Run tests once
pnpm test

# Watch mode
pnpm test:watch

# With coverage
pnpm test:coverage
```

---

## 🔧 Environment Variables

### Backend (`apps/backend/.env`)

| Variable | Description | Required |
|----------|-------------|----------|
| `APP_ENV` | `development` or `production` | Yes |
| `APP_PORT` | Server port (default: 8080) | Yes |
| `DATABASE_URL` | Supabase PostgreSQL connection string | Yes |
| `JWT_ACCESS_SECRET` | JWT access token secret (min 32 chars) | Yes |
| `JWT_REFRESH_SECRET` | JWT refresh token secret (min 32 chars, different from access) | Yes |
| `JWT_ACCESS_EXPIRY` | Access token expiry (default: 15m) | No |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry (default: 168h) | No |
| `CORS_ORIGINS` | Comma-separated allowed origins | Yes |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_ANON_KEY` | Supabase anon key | Yes |
| `SUPABASE_SERVICE_KEY` | Supabase service role key (for storage) | Yes |
| `SUPABASE_STORAGE_BUCKET` | Storage bucket name (default: media) | No |
| `BCRYPT_COST` | bcrypt cost factor (default: 12) | No |
| `LOG_LEVEL` | Log level: debug/info/warn/error | No |
| `LOG_FORMAT` | `json` (production) or `console` (development) | No |
| `COOKIE_DOMAIN` | Cookie domain | Yes |
| `COOKIE_SECURE` | Use secure cookies (true in production) | Yes |

### Frontend (`apps/frontend/.env.local`)

| Variable | Description | Required |
|----------|-------------|----------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | Yes |
| `NEXT_PUBLIC_SITE_URL` | Frontend URL | Yes |
| `NEXT_PUBLIC_DEFAULT_SITE_ID` | Default site UUID from database | Yes |

---

## 📦 Adding New Content Types

The CMS is designed to be extensible. To add a new content type:

### 1. Database
Add a new migration file (e.g., `009_create_blog_posts.sql`):
```sql
CREATE TABLE blog_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    content TEXT,
    -- ... your fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 2. Backend
1. Add domain model in `apps/backend/internal/domain/`
2. Add repository methods in `apps/backend/internal/repository/`
3. Add service methods in `apps/backend/internal/service/`
4. Add handler in `apps/backend/internal/handler/`
5. Register routes in `apps/backend/internal/router/router.go`

### 3. Frontend
1. Add TypeScript types in `apps/frontend/src/types/api.types.ts`
2. Add API functions in `apps/frontend/src/lib/api/`
3. Add admin page in `apps/frontend/src/app/(admin)/admin/`
4. Add landing section component in `apps/frontend/src/components/landing/`

---

## 📝 License

MIT
