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
export async function getPublicFeatures(siteId: string): Promise<Feature[]> {
  try {
    // Features are loaded via section contents or component API
    // For now we use the admin API with no auth for public data
    const res = await fetch(`${API_URL}/admin/features?site_id=${siteId}&is_active=true`, {
      next: { revalidate: 60 },
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return [];
    const json: ApiResponse<Feature[]> = await res.json();
    return json.data ?? [];
  } catch {
    return [];
  }
}

// Get testimonials for a site (public)
export async function getPublicTestimonials(siteId: string): Promise<Testimonial[]> {
  try {
    const res = await fetch(`${API_URL}/admin/testimonials?site_id=${siteId}&is_active=true`, {
      next: { revalidate: 60 },
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return [];
    const json: ApiResponse<Testimonial[]> = await res.json();
    return json.data ?? [];
  } catch {
    return [];
  }
}

// Get pricing plans for a site (public)
export async function getPublicPricingPlans(siteId: string): Promise<PricingPlan[]> {
  try {
    const res = await fetch(`${API_URL}/admin/pricing?site_id=${siteId}&is_active=true`, {
      next: { revalidate: 60 },
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return [];
    const json: ApiResponse<PricingPlan[]> = await res.json();
    return json.data ?? [];
  } catch {
    return [];
  }
}

// Get FAQs for a site (public)
export async function getPublicFAQs(siteId: string): Promise<FAQ[]> {
  try {
    const res = await fetch(`${API_URL}/admin/faqs?site_id=${siteId}&is_active=true`, {
      next: { revalidate: 60 },
      headers: { Accept: "application/json" },
    });
    if (!res.ok) return [];
    const json: ApiResponse<FAQ[]> = await res.json();
    return json.data ?? [];
  } catch {
    return [];
  }
}
