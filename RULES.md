# 📏 Development Rules & Guidelines

This document defines the coding standards, conventions, and best practices for this project. All contributors must follow these rules to maintain consistency and quality.

---

## 📋 Table of Contents
1. [General Rules](#general-rules)
2. [Git Conventions](#git-conventions)
3. [Backend Rules (Go)](#backend-rules-go)
4. [Frontend Rules (Next.js)](#frontend-rules-nextjs)
5. [Database Rules](#database-rules)
6. [API Design Rules](#api-design-rules)
7. [Security Rules](#security-rules)
8. [Testing Rules](#testing-rules)
9. [Documentation Rules](#documentation-rules)

---

## 1. General Rules

### Project Structure
- This is a **monorepo** — all apps live under `apps/`, shared packages under `packages/`
- Never mix backend and frontend code
- Each app must be independently runnable
- Use environment variables for ALL configuration — never hardcode secrets

### Code Quality
- No commented-out code in production branches
- No `TODO` comments without a linked issue
- All functions must have a single responsibility
- Maximum function length: 50 lines (split if longer)
- Maximum file length: 300 lines (split if longer)
- DRY principle: extract repeated logic into utilities

### Naming Conventions
- Use **descriptive names** — avoid abbreviations except well-known ones (ID, URL, HTTP, etc.)
- Boolean variables: prefix with `is`, `has`, `can`, `should` (e.g., `isActive`, `hasPermission`)
- Constants: `SCREAMING_SNAKE_CASE` in Go, `SCREAMING_SNAKE_CASE` in TypeScript
- Files: `snake_case` in Go, `kebab-case` in TypeScript/Next.js

---

## 2. Git Conventions

### Branch Naming
```
feature/short-description
fix/short-description
hotfix/short-description
chore/short-description
docs/short-description
refactor/short-description
```

### Commit Messages (Conventional Commits)
```
feat: add hero section content management
fix: resolve JWT refresh token expiry bug
chore: update dependencies
docs: add API documentation for auth endpoints
refactor: extract validation logic to middleware
test: add unit tests for user service
style: format code with gofmt
perf: optimize page content query with index
security: patch XSS vulnerability in content renderer
```

### Pull Request Rules
- PRs must have a description explaining what and why
- All CI checks must pass before merge
- At least 1 reviewer approval required
- Squash merge to keep history clean

---

## 3. Backend Rules (Go)

### Project Structure
```
apps/backend/
├── cmd/
│   └── api/
│       └── main.go           # Entry point only
├── internal/
│   ├── config/               # Configuration loading
│   ├── domain/               # Business entities (models)
│   ├── repository/           # Data access layer
│   ├── service/              # Business logic layer
│   ├── handler/              # HTTP handlers (controllers)
│   ├── middleware/           # HTTP middleware
│   ├── router/               # Route definitions
│   └── pkg/                  # Internal utilities
│       ├── auth/             # JWT utilities
│       ├── crypto/           # Hashing utilities
│       ├── logger/           # Logger setup
│       ├── response/         # Standard response helpers
│       ├── validator/        # Custom validators
│       └── errors/           # Custom error types
├── pkg/                      # Exportable utilities
├── migrations/               # SQL migrations (symlink to scripts/migrations)
├── .env.example
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### Architecture Rules (Clean Architecture)
```
Handler → Service → Repository → Database
```
- **Handler**: Only handles HTTP request/response, calls service
- **Service**: Contains business logic, calls repository
- **Repository**: Only handles database operations, returns domain models
- **Domain**: Pure Go structs, no framework dependencies

### Go Coding Standards

#### Error Handling
```go
// ✅ CORRECT - Always wrap errors with context
if err != nil {
    return fmt.Errorf("service.GetPage: %w", err)
}

// ❌ WRONG - Never ignore errors
result, _ := someFunction()

// ❌ WRONG - Never use panic in business logic
panic("something went wrong")
```

#### Interfaces
```go
// ✅ CORRECT - Define interfaces in the consumer package
// In service package:
type PageRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*domain.Page, error)
    FindAll(ctx context.Context, filter domain.PageFilter) ([]*domain.Page, error)
}

// ❌ WRONG - Don't define interfaces in the implementation package
```

#### Context
```go
// ✅ CORRECT - Always pass context as first parameter
func (s *pageService) GetPage(ctx context.Context, id uuid.UUID) (*domain.Page, error)

// ❌ WRONG - Never store context in struct
type Service struct {
    ctx context.Context // NEVER DO THIS
}
```

#### Dependency Injection
```go
// ✅ CORRECT - Constructor injection
func NewPageService(repo PageRepository, logger *zerolog.Logger) *PageService {
    return &PageService{repo: repo, logger: logger}
}

// ❌ WRONG - Global variables for dependencies
var globalDB *sql.DB
```

#### HTTP Handlers
```go
// ✅ CORRECT - Handler structure
func (h *PageHandler) GetPage(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        response.BadRequest(c, "invalid page ID")
        return
    }
    
    page, err := h.service.GetPage(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            response.NotFound(c, "page not found")
            return
        }
        response.InternalError(c, err)
        return
    }
    
    response.OK(c, page)
}
```

#### Response Format
```go
// ALL API responses must use this format:
{
    "success": true/false,
    "message": "Human readable message",
    "data": {...} or null,
    "errors": [...] or null,
    "meta": {
        "page": 1,
        "per_page": 20,
        "total": 100
    } // only for paginated responses
}
```

#### Validation
```go
// ✅ CORRECT - Use struct tags for validation
type CreatePageRequest struct {
    Title       string `json:"title" validate:"required,min=1,max=255"`
    Slug        string `json:"slug" validate:"required,slug"`
    Description string `json:"description" validate:"max=500"`
    IsPublished bool   `json:"is_published"`
}
```

#### Database Queries
```go
// ✅ CORRECT - Always use parameterized queries
query := `SELECT id, title FROM pages WHERE id = $1 AND site_id = $2`
row := db.QueryRowContext(ctx, query, id, siteID)

// ❌ WRONG - Never concatenate SQL strings
query := "SELECT * FROM pages WHERE id = " + id // SQL INJECTION!
```

#### Logging
```go
// ✅ CORRECT - Structured logging with context
logger.Info().
    Str("page_id", id.String()).
    Str("user_id", userID.String()).
    Msg("page retrieved successfully")

// ❌ WRONG - fmt.Println or log.Printf in production code
fmt.Println("got page:", page)
```

### Go File Naming
- `snake_case.go` for all Go files
- Test files: `snake_case_test.go`
- Interface files: `interface.go` or `{name}_interface.go`
- Mock files: `mock_{name}.go`

### Go Package Naming
- Short, lowercase, no underscores
- Descriptive: `handler`, `service`, `repository`, `domain`

---

## 4. Frontend Rules (Next.js)

### Project Structure
```
apps/frontend/
├── src/
│   ├── app/                  # Next.js App Router
│   │   ├── (public)/         # Public routes group
│   │   │   ├── page.tsx      # Landing page
│   │   │   └── [slug]/       # Dynamic pages
│   │   ├── (admin)/          # Admin routes group
│   │   │   ├── layout.tsx    # Admin layout with auth
│   │   │   └── admin/        # Admin pages
│   │   ├── api/              # API routes (if needed)
│   │   ├── layout.tsx        # Root layout
│   │   └── globals.css
│   ├── components/
│   │   ├── ui/               # shadcn/ui components
│   │   ├── common/           # Shared components
│   │   ├── landing/          # Landing page sections
│   │   └── admin/            # Admin panel components
│   ├── hooks/                # Custom React hooks
│   ├── lib/                  # Utilities and helpers
│   │   ├── api/              # API client functions
│   │   ├── auth/             # Auth utilities
│   │   └── utils/            # General utilities
│   ├── stores/               # Zustand stores
│   ├── types/                # TypeScript type definitions
│   └── constants/            # App constants
├── public/                   # Static assets
├── .env.example
├── next.config.ts
├── tailwind.config.ts
├── tsconfig.json
└── package.json
```

### TypeScript Rules
```typescript
// ✅ CORRECT - Always type function parameters and return types
async function getPage(id: string): Promise<Page> { ... }

// ✅ CORRECT - Use interfaces for objects, types for unions/primitives
interface Page {
  id: string;
  title: string;
  slug: string;
}

type PageStatus = 'draft' | 'published' | 'archived';

// ❌ WRONG - Never use `any`
const data: any = response.data; // FORBIDDEN

// ✅ CORRECT - Use unknown and narrow the type
const data: unknown = response.data;
if (isPage(data)) { ... }
```

### Component Rules
```typescript
// ✅ CORRECT - Functional components with explicit props interface
interface HeroSectionProps {
  title: string;
  subtitle: string;
  ctaText: string;
  ctaLink: string;
  backgroundImage?: string;
}

export function HeroSection({ title, subtitle, ctaText, ctaLink, backgroundImage }: HeroSectionProps) {
  return (...)
}

// ✅ CORRECT - Export named, not default (except page.tsx files)
export { HeroSection };

// ❌ WRONG - Default exports for components
export default function HeroSection() { ... }
```

### File Naming (Frontend)
- Components: `PascalCase.tsx` (e.g., `HeroSection.tsx`)
- Hooks: `use-kebab-case.ts` (e.g., `use-page-content.ts`)
- Utilities: `kebab-case.ts` (e.g., `api-client.ts`)
- Types: `kebab-case.types.ts` (e.g., `page.types.ts`)
- Constants: `kebab-case.constants.ts`
- Stores: `kebab-case.store.ts`

### State Management
```typescript
// ✅ CORRECT - Server state with React Query
const { data: page, isLoading } = useQuery({
  queryKey: ['page', slug],
  queryFn: () => getPage(slug),
});

// ✅ CORRECT - Client state with Zustand
const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  accessToken: null,
  setUser: (user) => set({ user }),
}));

// ❌ WRONG - Don't use useState for server data
const [page, setPage] = useState(null); // Use React Query instead
```

### API Calls
```typescript
// ✅ CORRECT - Centralized API client
// lib/api/client.ts
const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  withCredentials: true, // for httpOnly cookies
});

// ✅ CORRECT - Typed API functions
async function getPage(slug: string): Promise<ApiResponse<Page>> {
  const { data } = await apiClient.get<ApiResponse<Page>>(`/public/pages/${slug}`);
  return data;
}
```

### SEO & Metadata
```typescript
// ✅ CORRECT - Dynamic metadata from CMS
export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const page = await getPage(params.slug);
  return {
    title: page.seoTitle || page.title,
    description: page.seoDescription,
    openGraph: {
      title: page.ogTitle,
      description: page.ogDescription,
      images: [page.ogImage],
    },
  };
}
```

---

## 5. Database Rules

### Naming Conventions
- Tables: `snake_case`, plural (e.g., `pages`, `page_sections`)
- Columns: `snake_case` (e.g., `created_at`, `is_published`)
- Primary keys: always `id UUID DEFAULT gen_random_uuid()`
- Foreign keys: `{table_singular}_id` (e.g., `page_id`, `site_id`)
- Timestamps: always include `created_at` and `updated_at`
- Soft delete: use `deleted_at TIMESTAMPTZ` (nullable)
- Boolean columns: prefix with `is_` (e.g., `is_active`, `is_published`)

### Migration Rules
- One migration file per logical change
- Migrations are **irreversible** — always write rollback scripts
- File naming: `{number}_{description}.sql` (e.g., `001_create_users.sql`)
- Never modify existing migration files — create new ones
- Always test migrations on a copy of production data

### Index Rules
```sql
-- ✅ CORRECT - Index foreign keys
CREATE INDEX idx_pages_site_id ON pages(site_id);

-- ✅ CORRECT - Index frequently queried columns
CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_is_published ON pages(is_published) WHERE is_published = true;

-- ✅ CORRECT - Composite index for common query patterns
CREATE INDEX idx_page_sections_page_order ON page_sections(page_id, sort_order);
```

### Data Integrity
```sql
-- ✅ CORRECT - Always use constraints
ALTER TABLE pages ADD CONSTRAINT pages_slug_site_unique UNIQUE (slug, site_id);
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

-- ✅ CORRECT - Use CHECK constraints for enums
ALTER TABLE users ADD CONSTRAINT users_role_check 
  CHECK (role IN ('super_admin', 'admin', 'editor'));
```

### Soft Delete Pattern
```sql
-- All main tables should support soft delete
deleted_at TIMESTAMPTZ DEFAULT NULL,

-- Always filter in queries
WHERE deleted_at IS NULL
```

---

## 6. API Design Rules

### URL Structure
```
/api/v1/{scope}/{resource}/{id}/{sub-resource}

Examples:
GET  /api/v1/public/pages/home
GET  /api/v1/admin/pages
POST /api/v1/admin/pages
GET  /api/v1/admin/pages/123/sections
```

### HTTP Methods
- `GET` - Read (never modify state)
- `POST` - Create
- `PUT` - Full update (replace entire resource)
- `PATCH` - Partial update
- `DELETE` - Delete (soft delete preferred)

### Status Codes
- `200` - OK (GET, PUT, PATCH success)
- `201` - Created (POST success)
- `204` - No Content (DELETE success)
- `400` - Bad Request (validation error)
- `401` - Unauthorized (not authenticated)
- `403` - Forbidden (authenticated but no permission)
- `404` - Not Found
- `409` - Conflict (duplicate resource)
- `422` - Unprocessable Entity (business logic error)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

### Pagination
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### Filtering & Sorting
```
GET /api/v1/admin/pages?page=1&per_page=20&sort=created_at&order=desc&search=home&is_published=true
```

---

## 7. Security Rules

### Authentication
- NEVER store passwords in plain text — always bcrypt with cost ≥ 12
- NEVER store sensitive data in JWT payload (only user ID and role)
- Access tokens: short-lived (15 minutes)
- Refresh tokens: stored in httpOnly, Secure, SameSite=Strict cookie
- Invalidate refresh tokens on logout (store in DB)
- Implement token rotation on refresh

### Input Validation
- Validate ALL inputs on the server side (never trust client)
- Sanitize HTML content to prevent XSS
- Validate file uploads (type, size, content)
- Use parameterized queries — NEVER string concatenation in SQL

### Headers
```go
// Required security headers
c.Header("X-Content-Type-Options", "nosniff")
c.Header("X-Frame-Options", "DENY")
c.Header("X-XSS-Protection", "1; mode=block")
c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
// CSP header for production
c.Header("Content-Security-Policy", "default-src 'self'; ...")
```

### Rate Limiting
- Auth endpoints: 5 requests per minute per IP
- Public API: 100 requests per minute per IP
- Admin API: 200 requests per minute per user

### Secrets Management
- Use environment variables for ALL secrets
- Never commit `.env` files
- Rotate secrets regularly
- Use different secrets for each environment

---

## 8. Testing Rules

### Backend Testing
```
Unit tests:    apps/backend/internal/{package}/{file}_test.go
Integration:   apps/backend/tests/integration/
E2E:           apps/backend/tests/e2e/
```

- Minimum 80% code coverage for service layer
- Mock all external dependencies (DB, external APIs)
- Use table-driven tests for multiple scenarios
- Test both happy path and error cases

### Frontend Testing
```
Unit tests:    src/components/__tests__/
Integration:   src/app/__tests__/
E2E:           e2e/
```

- Test all custom hooks
- Test form validation
- E2E tests for critical user flows (login, content update)

---

## 9. Documentation Rules

### Code Documentation
```go
// ✅ CORRECT - Document exported functions
// GetPage retrieves a page by its slug for the given site.
// Returns ErrNotFound if the page does not exist or is not published.
func (s *pageService) GetPage(ctx context.Context, siteID uuid.UUID, slug string) (*domain.Page, error)
```

```typescript
/**
 * Fetches page content from the CMS API.
 * @param slug - The page slug identifier
 * @returns Page data with all sections and content
 * @throws ApiError if the page is not found or request fails
 */
async function getPage(slug: string): Promise<Page>
```

### API Documentation
- All endpoints must be documented in `docs/api/`
- Use OpenAPI 3.0 format
- Include request/response examples
- Document all error responses

### README Updates
- Update README when adding new features
- Keep environment variable list up to date
- Update API endpoint list when adding new routes

---

## 🚫 Forbidden Practices

1. **Never** commit secrets, API keys, or passwords
2. **Never** use `SELECT *` in production queries
3. **Never** disable SSL/TLS in production
4. **Never** log sensitive data (passwords, tokens, PII)
5. **Never** use `fmt.Println` for logging in production Go code
6. **Never** use `console.log` in production TypeScript code
7. **Never** use `any` type in TypeScript
8. **Never** skip input validation
9. **Never** use string concatenation for SQL queries
10. **Never** store tokens in localStorage (use httpOnly cookies for refresh tokens)
11. **Never** expose internal error details to clients
12. **Never** use deprecated packages or APIs
13. **Never** merge to main without code review
14. **Never** deploy without running tests

---

## ✅ Pre-commit Checklist

Before committing code, verify:
- [ ] No hardcoded secrets or credentials
- [ ] All new functions have proper error handling
- [ ] All new API endpoints have validation
- [ ] Tests pass locally
- [ ] Code is formatted (`gofmt` for Go, `prettier` for TS)
- [ ] No `console.log` or `fmt.Println` left in code
- [ ] Environment variables documented in `.env.example`
- [ ] README updated if needed
