"use client";

import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import {
  getPage,
  updatePage,
  listSections,
  createSection,
  updateSection,
  deleteSection,
  reorderSections,
  listContents,
  upsertContent,
  bulkUpsertContents,
} from "@/lib/api/pages.api";
import type { Page, PageSection, SectionContent, SectionType } from "@/types/api.types";

const SECTION_TYPES: { value: SectionType; label: string; icon: string }[] = [
  { value: "hero", label: "Hero", icon: "🦸" },
  { value: "features", label: "Features", icon: "⚡" },
  { value: "stats", label: "Stats", icon: "📊" },
  { value: "testimonials", label: "Testimonials", icon: "💬" },
  { value: "pricing", label: "Pricing", icon: "💰" },
  { value: "faq", label: "FAQ", icon: "❓" },
  { value: "cta", label: "CTA", icon: "🎯" },
  { value: "about", label: "About", icon: "ℹ️" },
  { value: "team", label: "Team", icon: "👥" },
  { value: "gallery", label: "Gallery", icon: "🖼️" },
  { value: "logos", label: "Logos", icon: "🏢" },
  { value: "contact", label: "Contact", icon: "📧" },
  { value: "video", label: "Video", icon: "🎥" },
  { value: "custom", label: "Custom", icon: "✏️" },
  { value: "html", label: "HTML", icon: "💻" },
];

export default function PageEditorPage() {
  const params = useParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const pageId = params.id as string;

  const [activeSection, setActiveSection] = useState<string | null>(null);
  const [showAddSection, setShowAddSection] = useState(false);
  const [editingPage, setEditingPage] = useState(false);

  // Fetch page data
  const { data: page, isLoading: pageLoading } = useQuery({
    queryKey: ["page", pageId],
    queryFn: () => getPage(pageId),
  });

  // Fetch sections
  const { data: sections = [], isLoading: sectionsLoading } = useQuery({
    queryKey: ["sections", pageId],
    queryFn: () => listSections(pageId),
  });

  // Fetch contents for active section
  const { data: contents = [] } = useQuery({
    queryKey: ["contents", activeSection],
    queryFn: () => listContents(activeSection!),
    enabled: !!activeSection,
  });

  // Mutations
  const updatePageMutation = useMutation({
    mutationFn: (data: Partial<Page>) => updatePage(pageId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["page", pageId] });
      setEditingPage(false);
    },
  });

  const createSectionMutation = useMutation({
    mutationFn: (data: Partial<PageSection>) => createSection(pageId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sections", pageId] });
      setShowAddSection(false);
    },
  });

  const deleteSectionMutation = useMutation({
    mutationFn: deleteSection,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sections", pageId] });
      setActiveSection(null);
    },
  });

  const upsertContentMutation = useMutation({
    mutationFn: ({ sectionId, data }: { sectionId: string; data: Partial<SectionContent> }) =>
      upsertContent(sectionId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contents", activeSection] });
    },
  });

  if (pageLoading) {
    return (
      <div className="p-8 flex items-center justify-center">
        <div className="text-gray-500">Loading page editor...</div>
      </div>
    );
  }

  if (!page) {
    return (
      <div className="p-8">
        <p className="text-red-600">Page not found.</p>
        <Link href="/admin/pages" className="text-primary hover:underline">
          Back to pages
        </Link>
      </div>
    );
  }

  const activeSectionData = sections.find((s) => s.id === activeSection);
  const contentsMap = contents.reduce((acc, c) => {
    acc[c.key] = c;
    return acc;
  }, {} as Record<string, SectionContent>);

  return (
    <div className="flex h-screen overflow-hidden">
      {/* Left sidebar - Sections list */}
      <div className="w-72 bg-white border-r border-gray-200 flex flex-col">
        {/* Page header */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center gap-2 mb-2">
            <Link href="/admin/pages" className="text-gray-400 hover:text-gray-600">
              ← Back
            </Link>
          </div>
          <h2 className="font-semibold text-gray-900 truncate">{page.title}</h2>
          <div className="flex items-center gap-2 mt-1">
            <span
              className={`text-xs px-2 py-0.5 rounded-full ${
                page.status === "published"
                  ? "bg-green-100 text-green-700"
                  : "bg-yellow-100 text-yellow-700"
              }`}
            >
              {page.status}
            </span>
            <span className="text-xs text-gray-500">/{page.slug}</span>
          </div>
        </div>

        {/* Page settings button */}
        <button
          onClick={() => setEditingPage(!editingPage)}
          className={`mx-4 mt-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
            editingPage
              ? "bg-primary/10 text-primary"
              : "text-gray-600 hover:bg-gray-100"
          }`}
        >
          ⚙️ Page Settings & SEO
        </button>

        {/* Sections list */}
        <div className="flex-1 overflow-y-auto p-4">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
              Sections
            </h3>
            <button
              onClick={() => setShowAddSection(true)}
              className="text-xs text-primary hover:underline"
            >
              + Add
            </button>
          </div>

          <div className="space-y-1">
            {sections.map((section) => (
              <button
                key={section.id}
                onClick={() => {
                  setActiveSection(section.id);
                  setEditingPage(false);
                }}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-colors flex items-center gap-2 ${
                  activeSection === section.id
                    ? "bg-primary/10 text-primary"
                    : "text-gray-700 hover:bg-gray-100"
                }`}
              >
                <span className="text-base">
                  {SECTION_TYPES.find((t) => t.value === section.type)?.icon || "📄"}
                </span>
                <div className="flex-1 min-w-0">
                  <div className="font-medium truncate">{section.name}</div>
                  <div className="text-xs text-gray-400 capitalize">{section.type}</div>
                </div>
                {!section.is_visible && (
                  <span className="text-xs text-gray-400">Hidden</span>
                )}
              </button>
            ))}
          </div>

          {sections.length === 0 && (
            <p className="text-sm text-gray-500 text-center py-4">
              No sections yet.
              <br />
              <button
                onClick={() => setShowAddSection(true)}
                className="text-primary hover:underline"
              >
                Add your first section
              </button>
            </p>
          )}
        </div>
      </div>

      {/* Main content area */}
      <div className="flex-1 overflow-y-auto bg-gray-50">
        {/* Page Settings Panel */}
        {editingPage && (
          <PageSettingsPanel
            page={page}
            onSave={(data) => updatePageMutation.mutate(data)}
            isSaving={updatePageMutation.isPending}
          />
        )}

        {/* Section Content Editor */}
        {activeSection && activeSectionData && !editingPage && (
          <SectionEditor
            section={activeSectionData}
            contents={contentsMap}
            onSaveContent={(key, value, type) =>
              upsertContentMutation.mutate({
                sectionId: activeSection,
                data: { key, value, type: type || "text" },
              })
            }
            onDeleteSection={() => {
              if (confirm(`Delete section "${activeSectionData.name}"?`)) {
                deleteSectionMutation.mutate(activeSection);
              }
            }}
            isSaving={upsertContentMutation.isPending}
          />
        )}

        {/* Empty state */}
        {!activeSection && !editingPage && (
          <div className="flex items-center justify-center h-full">
            <div className="text-center">
              <div className="text-6xl mb-4">📄</div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                Select a section to edit
              </h3>
              <p className="text-gray-500">
                Choose a section from the left panel to start editing its content.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Add Section Modal */}
      {showAddSection && (
        <AddSectionModal
          onAdd={(type, name) =>
            createSectionMutation.mutate({
              type,
              name,
              is_visible: true,
              sort_order: sections.length + 1,
            })
          }
          onClose={() => setShowAddSection(false)}
          isLoading={createSectionMutation.isPending}
        />
      )}
    </div>
  );
}

// ─── Page Settings Panel ──────────────────────────────────────────────────────

function PageSettingsPanel({
  page,
  onSave,
  isSaving,
}: {
  page: Page;
  onSave: (data: Partial<Page>) => void;
  isSaving: boolean;
}) {
  const [formData, setFormData] = useState({
    title: page.title,
    slug: page.slug,
    description: page.description || "",
    seo_title: page.seo_title || "",
    seo_description: page.seo_description || "",
    seo_keywords: page.seo_keywords || "",
    og_title: page.og_title || "",
    og_description: page.og_description || "",
    og_image: page.og_image || "",
    twitter_title: page.twitter_title || "",
    twitter_description: page.twitter_description || "",
    twitter_image: page.twitter_image || "",
    canonical_url: page.canonical_url || "",
    robots_meta: page.robots_meta || "index, follow",
    custom_head: page.custom_head || "",
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(formData);
  };

  return (
    <div className="p-8 max-w-3xl">
      <h2 className="text-xl font-bold text-gray-900 mb-6">Page Settings & SEO</h2>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Info */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4">Basic Information</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Page Title</label>
              <input
                type="text"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Slug</label>
              <input
                type="text"
                value={formData.slug}
                onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>
        </div>

        {/* SEO */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4">🔍 SEO Settings</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                SEO Title <span className="text-gray-400 font-normal">(max 60 chars)</span>
              </label>
              <input
                type="text"
                value={formData.seo_title}
                onChange={(e) => setFormData({ ...formData, seo_title: e.target.value })}
                maxLength={60}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
              <p className="text-xs text-gray-400 mt-1">{formData.seo_title.length}/60</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                SEO Description <span className="text-gray-400 font-normal">(max 160 chars)</span>
              </label>
              <textarea
                value={formData.seo_description}
                onChange={(e) => setFormData({ ...formData, seo_description: e.target.value })}
                maxLength={160}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
              <p className="text-xs text-gray-400 mt-1">{formData.seo_description.length}/160</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Keywords</label>
              <input
                type="text"
                value={formData.seo_keywords}
                onChange={(e) => setFormData({ ...formData, seo_keywords: e.target.value })}
                placeholder="keyword1, keyword2, keyword3"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Robots Meta</label>
              <select
                value={formData.robots_meta}
                onChange={(e) => setFormData({ ...formData, robots_meta: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              >
                <option value="index, follow">index, follow (default)</option>
                <option value="noindex, follow">noindex, follow</option>
                <option value="index, nofollow">index, nofollow</option>
                <option value="noindex, nofollow">noindex, nofollow</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Canonical URL</label>
              <input
                type="url"
                value={formData.canonical_url}
                onChange={(e) => setFormData({ ...formData, canonical_url: e.target.value })}
                placeholder="https://example.com/page"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>
        </div>

        {/* Open Graph */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4">📱 Open Graph (Social Sharing)</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">OG Title</label>
              <input
                type="text"
                value={formData.og_title}
                onChange={(e) => setFormData({ ...formData, og_title: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">OG Description</label>
              <textarea
                value={formData.og_description}
                onChange={(e) => setFormData({ ...formData, og_description: e.target.value })}
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">OG Image URL</label>
              <input
                type="url"
                value={formData.og_image}
                onChange={(e) => setFormData({ ...formData, og_image: e.target.value })}
                placeholder="https://example.com/og-image.jpg (1200x630px recommended)"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>
        </div>

        {/* Twitter Card */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4">🐦 Twitter Card</h3>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Twitter Title</label>
              <input
                type="text"
                value={formData.twitter_title}
                onChange={(e) => setFormData({ ...formData, twitter_title: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Twitter Description</label>
              <textarea
                value={formData.twitter_description}
                onChange={(e) => setFormData({ ...formData, twitter_description: e.target.value })}
                rows={2}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Twitter Image URL</label>
              <input
                type="url"
                value={formData.twitter_image}
                onChange={(e) => setFormData({ ...formData, twitter_image: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>
        </div>

        {/* Custom Head */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="font-semibold text-gray-900 mb-4">💻 Custom Head Tags</h3>
          <textarea
            value={formData.custom_head}
            onChange={(e) => setFormData({ ...formData, custom_head: e.target.value })}
            rows={5}
            placeholder="<!-- Add custom meta tags, scripts, or styles here -->"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <button
          type="submit"
          disabled={isSaving}
          className="w-full py-3 bg-primary text-white rounded-lg font-semibold hover:bg-primary/90 disabled:opacity-50 transition-colors"
        >
          {isSaving ? "Saving..." : "Save Page Settings"}
        </button>
      </form>
    </div>
  );
}

// ─── Section Editor ───────────────────────────────────────────────────────────

function SectionEditor({
  section,
  contents,
  onSaveContent,
  onDeleteSection,
  isSaving,
}: {
  section: PageSection;
  contents: Record<string, SectionContent>;
  onSaveContent: (key: string, value: string, type?: string) => void;
  onDeleteSection: () => void;
  isSaving: boolean;
}) {
  const contentFields = getContentFieldsForSection(section.type);

  return (
    <div className="p-8 max-w-3xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-xl font-bold text-gray-900">{section.name}</h2>
          <p className="text-sm text-gray-500 capitalize">{section.type} section</p>
        </div>
        <button
          onClick={onDeleteSection}
          className="text-sm text-red-600 hover:text-red-700 px-3 py-1.5 rounded-lg hover:bg-red-50 transition-colors"
        >
          🗑️ Delete Section
        </button>
      </div>

      <div className="space-y-4">
        {contentFields.map((field) => (
          <ContentField
            key={field.key}
            field={field}
            currentValue={contents[field.key]?.value || ""}
            onSave={(value) => onSaveContent(field.key, value, field.type)}
            isSaving={isSaving}
          />
        ))}

        {contentFields.length === 0 && (
          <div className="bg-white rounded-xl border border-gray-200 p-6 text-center text-gray-500">
            <p>This section type uses component data (managed separately).</p>
            <p className="text-sm mt-1">
              Go to the dedicated management page for this content type.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

// ─── Content Field ────────────────────────────────────────────────────────────

function ContentField({
  field,
  currentValue,
  onSave,
  isSaving,
}: {
  field: { key: string; label: string; type: string; placeholder?: string; multiline?: boolean };
  currentValue: string;
  onSave: (value: string) => void;
  isSaving: boolean;
}) {
  const [value, setValue] = useState(currentValue);
  const [isDirty, setIsDirty] = useState(false);

  const handleChange = (newValue: string) => {
    setValue(newValue);
    setIsDirty(newValue !== currentValue);
  };

  const handleSave = () => {
    onSave(value);
    setIsDirty(false);
  };

  return (
    <div className="bg-white rounded-xl border border-gray-200 p-5">
      <div className="flex items-center justify-between mb-2">
        <label className="text-sm font-medium text-gray-700">{field.label}</label>
        <span className="text-xs text-gray-400 font-mono">{field.key}</span>
      </div>

      {field.type === "image" ? (
        <div className="space-y-2">
          <input
            type="url"
            value={value}
            onChange={(e) => handleChange(e.target.value)}
            placeholder={field.placeholder || "https://example.com/image.jpg"}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          />
          {value && (
            <img
              src={value}
              alt="Preview"
              className="w-full max-h-40 object-cover rounded-lg border border-gray-200"
              onError={(e) => {
                (e.target as HTMLImageElement).style.display = "none";
              }}
            />
          )}
        </div>
      ) : field.multiline || field.type === "html" ? (
        <textarea
          value={value}
          onChange={(e) => handleChange(e.target.value)}
          placeholder={field.placeholder}
          rows={field.type === "html" ? 8 : 4}
          className={`w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary ${
            field.type === "html" ? "font-mono" : ""
          }`}
        />
      ) : (
        <input
          type={field.type === "link" ? "url" : "text"}
          value={value}
          onChange={(e) => handleChange(e.target.value)}
          placeholder={field.placeholder}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
        />
      )}

      {isDirty && (
        <div className="flex justify-end mt-2">
          <button
            onClick={handleSave}
            disabled={isSaving}
            className="px-4 py-1.5 bg-primary text-white text-sm rounded-lg hover:bg-primary/90 disabled:opacity-50 transition-colors"
          >
            {isSaving ? "Saving..." : "Save"}
          </button>
        </div>
      )}
    </div>
  );
}

// ─── Add Section Modal ────────────────────────────────────────────────────────

function AddSectionModal({
  onAdd,
  onClose,
  isLoading,
}: {
  onAdd: (type: SectionType, name: string) => void;
  onClose: () => void;
  isLoading: boolean;
}) {
  const [selectedType, setSelectedType] = useState<SectionType>("hero");
  const [name, setName] = useState("");

  const handleAdd = () => {
    const sectionName = name || SECTION_TYPES.find((t) => t.value === selectedType)?.label || selectedType;
    onAdd(selectedType, sectionName);
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl">
        <h3 className="text-lg font-bold text-gray-900 mb-4">Add New Section</h3>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">Section Type</label>
          <div className="grid grid-cols-3 gap-2 max-h-60 overflow-y-auto">
            {SECTION_TYPES.map((type) => (
              <button
                key={type.value}
                onClick={() => setSelectedType(type.value)}
                className={`p-2 rounded-lg border text-sm text-center transition-colors ${
                  selectedType === type.value
                    ? "border-primary bg-primary/10 text-primary"
                    : "border-gray-200 hover:border-gray-300"
                }`}
              >
                <div className="text-xl mb-1">{type.icon}</div>
                <div className="font-medium">{type.label}</div>
              </button>
            ))}
          </div>
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-1">Section Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={SECTION_TYPES.find((t) => t.value === selectedType)?.label}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
          />
        </div>

        <div className="flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleAdd}
            disabled={isLoading}
            className="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
          >
            {isLoading ? "Adding..." : "Add Section"}
          </button>
        </div>
      </div>
    </div>
  );
}

// ─── Content field definitions per section type ───────────────────────────────

function getContentFieldsForSection(type: SectionType) {
  const fields: { key: string; label: string; type: string; placeholder?: string; multiline?: boolean }[] = [];

  switch (type) {
    case "hero":
      return [
        { key: "badge_text", label: "Badge Text", type: "text", placeholder: "🚀 New Release v2.0" },
        { key: "title", label: "Main Title", type: "text", placeholder: "Build Something Amazing" },
        { key: "title_highlight", label: "Highlighted Word", type: "text", placeholder: "Amazing" },
        { key: "subtitle", label: "Subtitle", type: "text", multiline: true, placeholder: "Your compelling subtitle..." },
        { key: "cta_primary_text", label: "Primary CTA Text", type: "text", placeholder: "Get Started Free" },
        { key: "cta_primary_link", label: "Primary CTA Link", type: "link", placeholder: "#pricing" },
        { key: "cta_secondary_text", label: "Secondary CTA Text", type: "text", placeholder: "Watch Demo" },
        { key: "cta_secondary_link", label: "Secondary CTA Link", type: "link", placeholder: "#demo" },
        { key: "hero_image", label: "Hero Image URL", type: "image", placeholder: "https://..." },
        { key: "hero_image_alt", label: "Hero Image Alt Text", type: "text", placeholder: "Product screenshot" },
        { key: "social_proof_text", label: "Social Proof Text", type: "text", placeholder: "Trusted by 10,000+ teams" },
        { key: "video_url", label: "Demo Video URL", type: "link", placeholder: "https://youtube.com/..." },
      ];
    case "features":
      return [
        { key: "badge_text", label: "Badge Text", type: "text", placeholder: "Features" },
        { key: "title", label: "Section Title", type: "text", placeholder: "Everything you need" },
        { key: "subtitle", label: "Section Subtitle", type: "text", multiline: true },
      ];
    case "stats":
      return [
        { key: "stat_1_value", label: "Stat 1 Value", type: "text", placeholder: "10,000+" },
        { key: "stat_1_label", label: "Stat 1 Label", type: "text", placeholder: "Active Users" },
        { key: "stat_2_value", label: "Stat 2 Value", type: "text", placeholder: "99.9%" },
        { key: "stat_2_label", label: "Stat 2 Label", type: "text", placeholder: "Uptime SLA" },
        { key: "stat_3_value", label: "Stat 3 Value", type: "text", placeholder: "50+" },
        { key: "stat_3_label", label: "Stat 3 Label", type: "text", placeholder: "Countries" },
        { key: "stat_4_value", label: "Stat 4 Value", type: "text", placeholder: "24/7" },
        { key: "stat_4_label", label: "Stat 4 Label", type: "text", placeholder: "Support" },
      ];
    case "testimonials":
      return [
        { key: "badge_text", label: "Badge Text", type: "text", placeholder: "Testimonials" },
        { key: "title", label: "Section Title", type: "text", placeholder: "Loved by thousands" },
        { key: "subtitle", label: "Section Subtitle", type: "text", multiline: true },
      ];
    case "pricing":
      return [
        { key: "badge_text", label: "Badge Text", type: "text", placeholder: "Pricing" },
        { key: "title", label: "Section Title", type: "text", placeholder: "Simple, transparent pricing" },
        { key: "subtitle", label: "Section Subtitle", type: "text", multiline: true },
      ];
    case "faq":
      return [
        { key: "badge_text", label: "Badge Text", type: "text", placeholder: "FAQ" },
        { key: "title", label: "Section Title", type: "text", placeholder: "Frequently asked questions" },
        { key: "subtitle", label: "Section Subtitle", type: "text", multiline: true },
      ];
    case "cta":
      return [
        { key: "title", label: "CTA Title", type: "text", placeholder: "Ready to get started?" },
        { key: "subtitle", label: "CTA Subtitle", type: "text", multiline: true },
        { key: "cta_primary_text", label: "Primary CTA Text", type: "text", placeholder: "Start Free Trial" },
        { key: "cta_primary_link", label: "Primary CTA Link", type: "link", placeholder: "/signup" },
        { key: "cta_secondary_text", label: "Secondary CTA Text", type: "text", placeholder: "Contact Sales" },
        { key: "cta_secondary_link", label: "Secondary CTA Link", type: "link", placeholder: "/contact" },
      ];
    case "about":
      return [
        { key: "badge_text", label: "Badge Text", type: "text" },
        { key: "title", label: "Title", type: "text" },
        { key: "content", label: "Content", type: "html", multiline: true },
        { key: "image", label: "Image URL", type: "image" },
        { key: "image_alt", label: "Image Alt Text", type: "text" },
      ];
    case "contact":
      return [
        { key: "title", label: "Title", type: "text", placeholder: "Get in touch" },
        { key: "subtitle", label: "Subtitle", type: "text", multiline: true },
        { key: "email", label: "Contact Email", type: "text", placeholder: "hello@example.com" },
        { key: "phone", label: "Contact Phone", type: "text" },
        { key: "address", label: "Address", type: "text", multiline: true },
      ];
    case "custom":
    case "html":
      return [
        { key: "title", label: "Title", type: "text" },
        { key: "content", label: "HTML Content", type: "html", multiline: true },
      ];
    default:
      return [
        { key: "title", label: "Title", type: "text" },
        { key: "subtitle", label: "Subtitle", type: "text", multiline: true },
      ];
  }
}
