// Package docs provides Swagger API documentation for the Landing CMS API.
//
// # Landing CMS API
//
// A headless CMS API for managing landing pages, content, and components.
//
// ## Authentication
//
// The API uses JWT Bearer tokens for authentication. Include the token in the
// Authorization header: `Authorization: Bearer <token>`
//
// ## Base URL
//
// All API endpoints are prefixed with `/api/v1`
//
// ## Endpoints Overview
//
// ### Public Endpoints (no auth required)
//   - GET /api/v1/public/sites/:id - Get site info with public settings
//   - GET /api/v1/public/site/:slug - Get site by slug
//   - GET /api/v1/public/pages/:slug - Get published page with content
//   - GET /api/v1/public/pages - Get homepage
//   - GET /api/v1/public/navigation/:siteId/:identifier - Get navigation menu
//   - GET /api/v1/public/navigation/:siteId - Get default navigation
//
// ### Auth Endpoints
//   - POST /api/v1/auth/login - Login with email/password
//   - POST /api/v1/auth/logout - Logout (revoke refresh token)
//   - POST /api/v1/auth/refresh - Refresh access token
//   - GET /api/v1/auth/me - Get current user info (requires auth)
//   - POST /api/v1/auth/change-password - Change password (requires auth)
//
// ### Admin Endpoints (requires auth + role)
//
// #### Sites (admin+)
//   - GET /api/v1/admin/sites - List all sites
//   - POST /api/v1/admin/sites - Create a new site
//   - GET /api/v1/admin/sites/:id - Get site by ID
//   - PUT /api/v1/admin/sites/:id - Update site
//   - DELETE /api/v1/admin/sites/:id - Delete site (super_admin only)
//   - GET /api/v1/admin/sites/:id/settings - Get site settings
//   - PUT /api/v1/admin/sites/:id/settings - Bulk update settings
//   - PUT /api/v1/admin/sites/:id/settings/:key - Update single setting
//
// #### Pages (editor+)
//   - GET /api/v1/admin/pages - List pages
//   - POST /api/v1/admin/pages - Create page
//   - GET /api/v1/admin/pages/:id - Get page
//   - PUT /api/v1/admin/pages/:id - Update page
//   - DELETE /api/v1/admin/pages/:id - Delete page (admin+)
//   - PATCH /api/v1/admin/pages/:id/publish - Publish page
//   - PATCH /api/v1/admin/pages/:id/unpublish - Unpublish page
//   - GET /api/v1/admin/pages/:id/sections - List page sections
//   - POST /api/v1/admin/pages/:id/sections - Create section
//
// #### Sections (editor+)
//   - PUT /api/v1/admin/sections/:id - Update section
//   - DELETE /api/v1/admin/sections/:id - Delete section
//   - PATCH /api/v1/admin/sections/reorder - Reorder sections
//   - GET /api/v1/admin/sections/:id/contents - List section contents
//   - POST /api/v1/admin/sections/:id/contents - Upsert content
//   - POST /api/v1/admin/sections/:id/contents/bulk - Bulk upsert contents
//
// #### Contents (editor+)
//   - DELETE /api/v1/admin/contents/:id - Delete content
//
// #### Features (editor+)
//   - GET /api/v1/admin/features - List features
//   - POST /api/v1/admin/features - Create feature
//   - PUT /api/v1/admin/features/:id - Update feature
//   - DELETE /api/v1/admin/features/:id - Delete feature
//
// #### Testimonials (editor+)
//   - GET /api/v1/admin/testimonials - List testimonials
//   - POST /api/v1/admin/testimonials - Create testimonial
//   - PUT /api/v1/admin/testimonials/:id - Update testimonial
//   - DELETE /api/v1/admin/testimonials/:id - Delete testimonial
//
// #### Pricing Plans (editor+)
//   - GET /api/v1/admin/pricing - List pricing plans
//   - POST /api/v1/admin/pricing - Create pricing plan
//   - PUT /api/v1/admin/pricing/:id - Update pricing plan
//   - DELETE /api/v1/admin/pricing/:id - Delete pricing plan
//
// #### FAQs (editor+)
//   - GET /api/v1/admin/faqs - List FAQs
//   - POST /api/v1/admin/faqs - Create FAQ
//   - PUT /api/v1/admin/faqs/:id - Update FAQ
//   - DELETE /api/v1/admin/faqs/:id - Delete FAQ
//
// #### Navigation (editor+)
//   - GET /api/v1/admin/navigation - List navigation menus
//   - POST /api/v1/admin/navigation/:menuId/items - Create navigation item
//   - PUT /api/v1/admin/navigation/items/:id - Update navigation item
//   - DELETE /api/v1/admin/navigation/items/:id - Delete navigation item
//
// #### Media (editor+)
//   - GET /api/v1/admin/media - List media files
//   - POST /api/v1/admin/media/upload - Upload media file
//   - PUT /api/v1/admin/media/:id - Update media metadata
//   - DELETE /api/v1/admin/media/:id - Delete media file
//
// #### Users (admin+)
//   - GET /api/v1/admin/users - List users
//   - POST /api/v1/admin/users - Create user
//   - GET /api/v1/admin/users/:id - Get user
//   - PUT /api/v1/admin/users/:id - Update user
//   - DELETE /api/v1/admin/users/:id - Delete user (super_admin only)
//
// #### Audit Logs (admin+)
//   - GET /api/v1/admin/audit-logs - List audit logs
//
// ## Response Format
//
// All responses follow this structure:
//
//	{
//	  "success": true,
//	  "message": "success",
//	  "data": {...},
//	  "meta": {...}  // for paginated responses
//	}
//
// ## Error Response Format
//
//	{
//	  "success": false,
//	  "message": "error description"
//	}
//
// ## Pagination
//
// Paginated endpoints accept `page` and `per_page` query parameters.
// Response includes `meta` with pagination info:
//
//	{
//	  "page": 1,
//	  "per_page": 20,
//	  "total": 100,
//	  "total_pages": 5
//	}
//
// ## Roles
//
//   - editor: Can manage content (pages, sections, features, etc.)
//   - admin: Can manage users, sites, and all content
//   - super_admin: Full access including delete operations
package docs
