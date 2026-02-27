import apiClient, { extractApiError } from "./client";
import type {
  ApiResponse,
  Feature,
  Testimonial,
  PricingPlan,
  FAQ,
  NavigationMenu,
  Media,
  AuditLog,
} from "@/types/api.types";

// ─── Features ─────────────────────────────────────────────────────────────────

export async function listFeatures(siteId: string): Promise<Feature[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<Feature[]>>("/admin/features", {
      params: { site_id: siteId },
    });
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function createFeature(input: Partial<Feature>): Promise<Feature> {
  try {
    const { data } = await apiClient.post<ApiResponse<Feature>>("/admin/features", input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function updateFeature(id: string, input: Partial<Feature>): Promise<Feature> {
  try {
    const { data } = await apiClient.put<ApiResponse<Feature>>(`/admin/features/${id}`, input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function deleteFeature(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/features/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Testimonials ─────────────────────────────────────────────────────────────

export async function listTestimonials(siteId: string): Promise<Testimonial[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<Testimonial[]>>("/admin/testimonials", {
      params: { site_id: siteId },
    });
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function createTestimonial(input: Partial<Testimonial>): Promise<Testimonial> {
  try {
    const { data } = await apiClient.post<ApiResponse<Testimonial>>("/admin/testimonials", input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function updateTestimonial(id: string, input: Partial<Testimonial>): Promise<Testimonial> {
  try {
    const { data } = await apiClient.put<ApiResponse<Testimonial>>(`/admin/testimonials/${id}`, input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function deleteTestimonial(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/testimonials/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Pricing Plans ────────────────────────────────────────────────────────────

export async function listPricingPlans(siteId: string): Promise<PricingPlan[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<PricingPlan[]>>("/admin/pricing", {
      params: { site_id: siteId },
    });
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function createPricingPlan(input: Partial<PricingPlan>): Promise<PricingPlan> {
  try {
    const { data } = await apiClient.post<ApiResponse<PricingPlan>>("/admin/pricing", input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function updatePricingPlan(id: string, input: Partial<PricingPlan>): Promise<PricingPlan> {
  try {
    const { data } = await apiClient.put<ApiResponse<PricingPlan>>(`/admin/pricing/${id}`, input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function deletePricingPlan(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/pricing/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── FAQs ─────────────────────────────────────────────────────────────────────

export async function listFAQs(siteId: string): Promise<FAQ[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<FAQ[]>>("/admin/faqs", {
      params: { site_id: siteId },
    });
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function createFAQ(input: Partial<FAQ>): Promise<FAQ> {
  try {
    const { data } = await apiClient.post<ApiResponse<FAQ>>("/admin/faqs", input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function updateFAQ(id: string, input: Partial<FAQ>): Promise<FAQ> {
  try {
    const { data } = await apiClient.put<ApiResponse<FAQ>>(`/admin/faqs/${id}`, input);
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function deleteFAQ(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/faqs/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Navigation ───────────────────────────────────────────────────────────────

export async function getNavigation(siteId: string, identifier = "header"): Promise<NavigationMenu> {
  try {
    const { data } = await apiClient.get<ApiResponse<NavigationMenu>>(
      `/public/navigation/${siteId}/${identifier}`
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

export async function listNavigation(siteId: string): Promise<NavigationMenu[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<NavigationMenu[]>>("/admin/navigation", {
      params: { site_id: siteId },
    });
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Audit Logs ───────────────────────────────────────────────────────────────

export async function listAuditLogs(
  siteId?: string,
  params?: { action?: string; resource_type?: string; page?: number }
): Promise<{ data: AuditLog[]; meta: ApiResponse["meta"] }> {
  try {
    const { data } = await apiClient.get<ApiResponse<AuditLog[]>>("/admin/audit-logs", {
      params: { site_id: siteId, ...params },
    });
    return { data: data.data ?? [], meta: data.meta };
  } catch (error) {
    throw extractApiError(error);
  }
}
