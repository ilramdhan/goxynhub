# 🚀 Dynamic Landing Page CMS

A full-stack monorepo project for a **dynamic landing page** with a CMS admin panel. All landing page content (title, description, images, links, OG meta, etc.) can be managed through the admin panel.

---

## 📁 Project Structure

```
/
├── apps/
│   ├── backend/          # Go + Gin REST API
│   └── frontend/         # Next.js 14 App Router
├── packages/
│   └── shared/           # Shared types/schemas (future)
├── docs/
│   ├── api/              # API documentation
│   └── erd/              # ERD diagrams
├── scripts/
│   └── migrations/       # SQL migration files
├── docker/
│   ├── backend/
│   └── frontend/
├── docker-compose.yml
├── docker-compose.prod.yml
├── README.md
└── RULES.md
```

---

## 🏗️ Architecture Overview

### Backend (Go + Gin)
- **Framework**: [Gin](https://github.com/gin-gonic/gin) - High performance HTTP framework
- **Database**: Supabase (PostgreSQL)
- **Auth**: JWT (Access Token + Refresh Token) with secure httpOnly cookies
- **Architecture**: Clean Architecture (Handler → Service → Repository)
- **ORM**: `pgx` + `sqlx` for raw SQL with type safety
- **Validation**: `go-playground/validator`
- **Config**: `viper` for environment management
- **Logging**: `zerolog` for structured logging
- **Rate Limiting**: Redis-based or in-memory rate limiter
- **Security**: CORS, Helmet headers, CSRF protection, input sanitization

### Frontend (Next.js 14)
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS + shadcn/ui
- **State Management**: Zustand + React Query (TanStack Query)
- **Auth**: NextAuth.js v5 or custom JWT handling
- **Forms**: React Hook Form + Zod validation
- **Image Optimization**: Next.js Image component
- **SEO**: Dynamic metadata API (OG, Twitter cards, etc.)

---

## 🗄️ Database Design (ERD Summary)

### Core Tables

| Table | Description |
|-------|-------------|
| `users` | Admin users with roles |
| `sites` | Multi-site support (one site per deployment or multiple) |
| `site_settings` | Global site settings (SEO, analytics, etc.) |
| `pages` | Landing pages |
| `page_sections` | Sections within a page (hero, features, pricing, etc.) |
| `section_contents` | Key-value content within sections |
| `media` | Uploaded media files (images, videos) |
| `navigation_menus` | Navigation menu structure |
| `navigation_items` | Individual menu items |
| `testimonials` | Testimonials/reviews |
| `faqs` | FAQ items |
| `pricing_plans` | Pricing plan cards |
| `features` | Feature list items |
| `cta_buttons` | Call-to-action buttons |
| `social_links` | Social media links |
| `audit_logs` | Admin action audit trail |
| `refresh_tokens` | JWT refresh token store |

---

## 🔐 Authentication & Security

### Auth Flow
1. Admin logs in with email + password
2. Server validates credentials, returns:
   - **Access Token** (JWT, 15min expiry) → stored in memory/localStorage
   - **Refresh Token** (JWT, 7 days) → stored in httpOnly secure cookie
3. Access token used for API calls
4. Refresh token used to get new access token silently
5. Logout invalidates refresh token in DB

### Security Measures
- Password hashing with `bcrypt` (cost factor 12)
- JWT with RS256 algorithm (asymmetric keys)
- Rate limiting on auth endpoints (5 req/min)
- CORS whitelist
- SQL injection prevention (parameterized queries)
- XSS prevention (input sanitization)
- CSRF protection
- Audit logging for all admin actions
- Role-based access control (RBAC)

---

## 🚀 Getting Started

### Prerequisites
- Go 1.22+
- Node.js 20+
- pnpm 9+
- Docker & Docker Compose
- Supabase account

### Environment Setup

#### Backend
```bash
cp apps/backend/.env.example apps/backend/.env
# Fill in your Supabase credentials and JWT secrets
```

#### Frontend
```bash
cp apps/frontend/.env.example apps/frontend/.env.local
# Fill in your API URL and other configs
```

### Database Setup
1. Create a new Supabase project
2. Run migrations in order:
```bash
# Run all migrations
psql $DATABASE_URL -f scripts/migrations/001_create_users.sql
psql $DATABASE_URL -f scripts/migrations/002_create_sites.sql
psql $DATABASE_URL -f scripts/migrations/003_create_pages.sql
psql $DATABASE_URL -f scripts/migrations/004_create_content.sql
psql $DATABASE_URL -f scripts/migrations/005_create_media.sql
psql $DATABASE_URL -f scripts/migrations/006_create_navigation.sql
psql $DATABASE_URL -f scripts/migrations/007_create_components.sql
psql $DATABASE_URL -f scripts/migrations/008_create_audit.sql
psql $DATABASE_URL -f scripts/migrations/009_seed_data.sql
```

### Running Locally

#### With Docker Compose
```bash
docker-compose up -d
```

#### Without Docker
```bash
# Backend
cd apps/backend
go mod download
go run cmd/api/main.go

# Frontend
cd apps/frontend
pnpm install
pnpm dev
```

---

## 📡 API Endpoints

### Public Endpoints
```
GET  /api/v1/public/site/:slug          # Get site info
GET  /api/v1/public/pages/:slug         # Get page with all content
GET  /api/v1/public/navigation/:siteId  # Get navigation
```

### Auth Endpoints
```
POST /api/v1/auth/login                 # Login
POST /api/v1/auth/logout                # Logout
POST /api/v1/auth/refresh               # Refresh access token
POST /api/v1/auth/change-password       # Change password
```

### Admin Endpoints (Protected)
```
# Sites
GET    /api/v1/admin/sites
POST   /api/v1/admin/sites
GET    /api/v1/admin/sites/:id
PUT    /api/v1/admin/sites/:id
DELETE /api/v1/admin/sites/:id

# Pages
GET    /api/v1/admin/pages
POST   /api/v1/admin/pages
GET    /api/v1/admin/pages/:id
PUT    /api/v1/admin/pages/:id
DELETE /api/v1/admin/pages/:id

# Sections
GET    /api/v1/admin/pages/:pageId/sections
POST   /api/v1/admin/pages/:pageId/sections
PUT    /api/v1/admin/sections/:id
DELETE /api/v1/admin/sections/:id
PATCH  /api/v1/admin/sections/reorder

# Content
GET    /api/v1/admin/sections/:sectionId/contents
POST   /api/v1/admin/sections/:sectionId/contents
PUT    /api/v1/admin/contents/:id
DELETE /api/v1/admin/contents/:id

# Media
GET    /api/v1/admin/media
POST   /api/v1/admin/media/upload
DELETE /api/v1/admin/media/:id

# Navigation
GET    /api/v1/admin/navigation
POST   /api/v1/admin/navigation
PUT    /api/v1/admin/navigation/:id
DELETE /api/v1/admin/navigation/:id

# Components (Testimonials, FAQs, Pricing, Features)
GET    /api/v1/admin/testimonials
POST   /api/v1/admin/testimonials
PUT    /api/v1/admin/testimonials/:id
DELETE /api/v1/admin/testimonials/:id

# Users (Super Admin only)
GET    /api/v1/admin/users
POST   /api/v1/admin/users
PUT    /api/v1/admin/users/:id
DELETE /api/v1/admin/users/:id

# Audit Logs
GET    /api/v1/admin/audit-logs
```

---

## 🌐 Frontend Pages

### Public
- `/` - Landing page (dynamic content from CMS)
- `/[slug]` - Dynamic pages

### Admin Panel
- `/admin/login` - Admin login
- `/admin/dashboard` - Dashboard overview
- `/admin/sites` - Site management
- `/admin/pages` - Page management
- `/admin/pages/[id]/sections` - Section management
- `/admin/media` - Media library
- `/admin/navigation` - Navigation management
- `/admin/testimonials` - Testimonials
- `/admin/faqs` - FAQs
- `/admin/pricing` - Pricing plans
- `/admin/features` - Features
- `/admin/settings` - Site settings
- `/admin/users` - User management
- `/admin/audit-logs` - Audit logs

---

## 🐳 Deployment

### Production with Docker
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Variables

#### Backend
| Variable | Description | Required |
|----------|-------------|----------|
| `APP_ENV` | Environment (development/production) | Yes |
| `APP_PORT` | Server port | Yes |
| `DATABASE_URL` | Supabase PostgreSQL URL | Yes |
| `JWT_ACCESS_SECRET` | JWT access token secret (RS256 private key) | Yes |
| `JWT_REFRESH_SECRET` | JWT refresh token secret | Yes |
| `JWT_ACCESS_EXPIRY` | Access token expiry (e.g., 15m) | Yes |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry (e.g., 168h) | Yes |
| `CORS_ORIGINS` | Allowed CORS origins | Yes |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_ANON_KEY` | Supabase anon key | Yes |
| `SUPABASE_SERVICE_KEY` | Supabase service role key | Yes |
| `STORAGE_BUCKET` | Supabase storage bucket name | Yes |
| `RATE_LIMIT_REQUESTS` | Rate limit requests per window | No |
| `RATE_LIMIT_WINDOW` | Rate limit window duration | No |

#### Frontend
| Variable | Description | Required |
|----------|-------------|----------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | Yes |
| `NEXT_PUBLIC_SITE_URL` | Frontend URL | Yes |
| `NEXTAUTH_SECRET` | NextAuth secret | Yes |
| `NEXTAUTH_URL` | NextAuth URL | Yes |

---

## 📝 License

MIT
