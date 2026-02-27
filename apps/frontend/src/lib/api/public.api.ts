/**
 * Public API functions - no authentication required
 * Used for SSR on the landing page
 */
import type {
  Site,
  Page,
  Feature,
  Testimonial,
  PricingPlan,
  FAQ,
  NavigationMenu,
  ApiResponse,
} from "@/types/api.types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function fetchPublic<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    next: { revalidate: 60 }, // ISR: revalidate every 60 seconds
    headers: { Accept: "application/json" },
  });

  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }

  const json: ApiResponse<T> = await res.json();
  if (!json.success || !json.data) {
    throw new Error(json.message || "API returned no data");
  }
  return json.data;
}

async function fetchPublicSafe<T>(path: string, fallback: T): Promise<T> {
  try {
    return await fetchPublic<T>(path);
  } catch {
    return fallback;
  }
}

// Get site info with public settings
export async function getPublicSite(siteId: string): Promise<Site> {
  return fetchPublic<Site>(`/public/sites/${siteId}`);
}

// Get page with all sections and content
export async function getPublicPage(siteId: string, slug?: string): Promise<Page> {
  const url = slug ? `/public/pages/${slug}?site_id=${siteId}` : `/public/pages?site_id=${siteId}`;
  return fetchPublic<Page>(url);
}

// Get navigation menu
export async function getPublicNavigation(siteId: string, identifier = "header"): Promise<NavigationMenu> {
  return fetchPublic<NavigationMenu>(`/public/navigation/${siteId}/${identifier}`);
}

// Get features for a site (public - used in landing page)
// Note: These use admin endpoints without auth for public data
// In production, you should create dedicated public endpoints
export async function getPublicFeatures(siteId: string): Promise<Feature[]> {
  const res = await fetch(`${API_URL}/admin/features?site_id=${siteId}&is_active=true`, {
    next: { revalidate: 60 },
    headers: { Accept: "application/json" },
  });
  if (!res.ok) return [];
  const json: ApiResponse<Feature[]> = await res.json();
  return json.data ?? [];
}

// Get testimonials for a site (public)
export async function getPublicTestimonials(siteId: string): Promise<Testimonial[]> {
  const res = await fetch(`${API_URL}/admin/testimonials?site_id=${siteId}&is_active=true`, {
    next: { revalidate: 60 },
    headers: { Accept: "application/json" },
  });
  if (!res.ok) return [];
  const json: ApiResponse<Testimonial[]> = await res.json();
  return json.data ?? [];
}

// Get pricing plans for a site (public)
export async function getPublicPricingPlans(siteId: string): Promise<PricingPlan[]> {
  const res = await fetch(`${API_URL}/admin/pricing?site_id=${siteId}&is_active=true`, {
    next: { revalidate: 60 },
    headers: { Accept: "application/json" },
  });
  if (!res.ok) return [];
  const json: ApiResponse<PricingPlan[]> = await res.json();
  return json.data ?? [];
}

// Get FAQs for a site (public)
export async function getPublicFAQs(siteId: string): Promise<FAQ[]> {
  const res = await fetch(`${API_URL}/admin/faqs?site_id=${siteId}&is_active=true`, {
    next: { revalidate: 60 },
    headers: { Accept: "application/json" },
  });
  if (!res.ok) return [];
  const json: ApiResponse<FAQ[]> = await res.json();
  return json.data ?? [];
}

// Export fetchPublicSafe for use in other modules
export { fetchPublicSafe };
