import apiClient, { extractApiError } from "./client";
import type {
  ApiResponse,
  Page,
  PageSection,
  SectionContent,
  PageQueryParams,
  PaginatedResponse,
} from "@/types/api.types";

// ─── Public API ──────────────────────────────────────────────────────────────

// Get public page with all content (for landing page rendering)
export async function getPublicPage(
  siteId: string,
  slug?: string
): Promise<Page> {
  try {
    const url = slug ? `/public/pages/${slug}` : "/public/pages";
    const { data } = await apiClient.get<ApiResponse<Page>>(url, {
      params: { site_id: siteId },
    });
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Admin API ────────────────────────────────────────────────────────────────

// List pages
export async function listPages(
  params?: PageQueryParams
): Promise<PaginatedResponse<Page>> {
  try {
    const { data } = await apiClient.get<ApiResponse<Page[]>>("/admin/pages", {
      params,
    });
    return {
      data: data.data ?? [],
      meta: data.meta ?? {
        page: 1,
        per_page: 20,
        total: 0,
        total_pages: 0,
      },
    };
  } catch (error) {
    throw extractApiError(error);
  }
}

// Get page by ID
export async function getPage(id: string): Promise<Page> {
  try {
    const { data } = await apiClient.get<ApiResponse<Page>>(
      `/admin/pages/${id}`
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Create page
export async function createPage(
  input: Partial<Page>
): Promise<Page> {
  try {
    const { data } = await apiClient.post<ApiResponse<Page>>(
      "/admin/pages",
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Update page
export async function updatePage(
  id: string,
  input: Partial<Page>
): Promise<Page> {
  try {
    const { data } = await apiClient.put<ApiResponse<Page>>(
      `/admin/pages/${id}`,
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Delete page
export async function deletePage(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/pages/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// Publish page
export async function publishPage(id: string): Promise<Page> {
  try {
    const { data } = await apiClient.patch<ApiResponse<Page>>(
      `/admin/pages/${id}/publish`
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Unpublish page
export async function unpublishPage(id: string): Promise<Page> {
  try {
    const { data } = await apiClient.patch<ApiResponse<Page>>(
      `/admin/pages/${id}/unpublish`
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Sections ─────────────────────────────────────────────────────────────────

// List sections for a page
export async function listSections(pageId: string): Promise<PageSection[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<PageSection[]>>(
      `/admin/pages/${pageId}/sections`
    );
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

// Create section
export async function createSection(
  pageId: string,
  input: Partial<PageSection>
): Promise<PageSection> {
  try {
    const { data } = await apiClient.post<ApiResponse<PageSection>>(
      `/admin/pages/${pageId}/sections`,
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Update section
export async function updateSection(
  id: string,
  input: Partial<PageSection>
): Promise<PageSection> {
  try {
    const { data } = await apiClient.put<ApiResponse<PageSection>>(
      `/admin/sections/${id}`,
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Delete section
export async function deleteSection(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/sections/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}

// Reorder sections
export async function reorderSections(
  sections: Array<{ id: string; sort_order: number }>
): Promise<void> {
  try {
    await apiClient.patch("/admin/sections/reorder", { sections });
  } catch (error) {
    throw extractApiError(error);
  }
}

// ─── Contents ─────────────────────────────────────────────────────────────────

// List contents for a section
export async function listContents(
  sectionId: string
): Promise<SectionContent[]> {
  try {
    const { data } = await apiClient.get<ApiResponse<SectionContent[]>>(
      `/admin/sections/${sectionId}/contents`
    );
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

// Upsert content
export async function upsertContent(
  sectionId: string,
  input: Partial<SectionContent>
): Promise<SectionContent> {
  try {
    const { data } = await apiClient.post<ApiResponse<SectionContent>>(
      `/admin/sections/${sectionId}/contents`,
      input
    );
    if (!data.data) throw new Error("No data returned");
    return data.data;
  } catch (error) {
    throw extractApiError(error);
  }
}

// Bulk upsert contents
export async function bulkUpsertContents(
  sectionId: string,
  inputs: Partial<SectionContent>[]
): Promise<SectionContent[]> {
  try {
    const { data } = await apiClient.post<ApiResponse<SectionContent[]>>(
      `/admin/sections/${sectionId}/contents/bulk`,
      inputs
    );
    return data.data ?? [];
  } catch (error) {
    throw extractApiError(error);
  }
}

// Delete content
export async function deleteContent(id: string): Promise<void> {
  try {
    await apiClient.delete(`/admin/contents/${id}`);
  } catch (error) {
    throw extractApiError(error);
  }
}
