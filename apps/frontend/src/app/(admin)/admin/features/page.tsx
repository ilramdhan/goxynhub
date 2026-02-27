"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  listFeatures,
  createFeature,
  updateFeature,
  deleteFeature,
} from "@/lib/api/components.api";
import type { Feature } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

const COMMON_ICONS = [
  { value: "zap", label: "⚡ Zap" },
  { value: "shield", label: "🛡️ Shield" },
  { value: "plug", label: "🔌 Plug" },
  { value: "bar-chart", label: "📊 Chart" },
  { value: "users", label: "👥 Users" },
  { value: "headphones", label: "🎧 Support" },
  { value: "star", label: "⭐ Star" },
  { value: "lock", label: "🔒 Lock" },
  { value: "globe", label: "🌐 Globe" },
  { value: "code", label: "💻 Code" },
  { value: "rocket", label: "🚀 Rocket" },
  { value: "heart", label: "❤️ Heart" },
  { value: "cloud", label: "☁️ Cloud" },
  { value: "database", label: "🗄️ Database" },
  { value: "settings", label: "⚙️ Settings" },
];

export default function FeaturesPage() {
  const queryClient = useQueryClient();
  const [editingItem, setEditingItem] = useState<Feature | null>(null);
  const [showForm, setShowForm] = useState(false);

  const { data: features = [], isLoading } = useQuery({
    queryKey: ["features", SITE_ID],
    queryFn: () => listFeatures(SITE_ID),
  });

  const createMutation = useMutation({
    mutationFn: createFeature,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["features", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Feature> }) =>
      updateFeature(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["features", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteFeature,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["features", SITE_ID] }),
  });

  const handleSubmit = (data: Partial<Feature>) => {
    if (editingItem) {
      updateMutation.mutate({ id: editingItem.id, data });
    } else {
      createMutation.mutate({ ...data, site_id: SITE_ID });
    }
  };

  const iconMap: Record<string, string> = {
    zap: "⚡", shield: "🛡️", plug: "🔌", "bar-chart": "📊", users: "👥",
    headphones: "🎧", star: "⭐", lock: "🔒", globe: "🌐", code: "💻",
    rocket: "🚀", heart: "❤️", cloud: "☁️", database: "🗄️", settings: "⚙️",
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Features</h1>
          <p className="text-gray-600 mt-1">Manage your product features list</p>
        </div>
        <button
          onClick={() => { setEditingItem(null); setShowForm(true); }}
          className="px-4 py-2 bg-primary text-white rounded-lg font-semibold text-sm hover:bg-primary/90 transition-colors"
        >
          + Add Feature
        </button>
      </div>

      {isLoading ? (
        <div className="text-gray-500">Loading...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature) => (
            <div key={feature.id} className="bg-white rounded-xl border border-gray-200 p-6">
              <div className="flex items-start justify-between mb-3">
                <div
                  className="text-3xl"
                  style={{ color: feature.icon_color || undefined }}
                >
                  {(feature.icon && iconMap[feature.icon]) || "✨"}
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => { setEditingItem(feature); setShowForm(true); }}
                    className="text-xs text-primary hover:underline"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => confirm("Delete this feature?") && deleteMutation.mutate(feature.id)}
                    className="text-xs text-red-600 hover:underline"
                  >
                    Delete
                  </button>
                </div>
              </div>
              <h3 className="font-semibold text-gray-900 mb-1">{feature.title}</h3>
              {feature.description && (
                <p className="text-sm text-gray-600 line-clamp-2">{feature.description}</p>
              )}
              <div className="mt-3 flex gap-2">
                <span className={`text-xs px-2 py-0.5 rounded-full ${feature.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-600"}`}>
                  {feature.is_active ? "Active" : "Inactive"}
                </span>
                <span className="text-xs text-gray-400">Order: {feature.sort_order}</span>
              </div>
            </div>
          ))}
          {features.length === 0 && (
            <div className="col-span-3 text-center py-12 text-gray-500">
              No features yet. Add your first one!
            </div>
          )}
        </div>
      )}

      {showForm && (
        <FeatureFormModal
          item={editingItem}
          onSubmit={handleSubmit}
          onClose={() => { setShowForm(false); setEditingItem(null); }}
          isLoading={createMutation.isPending || updateMutation.isPending}
          icons={COMMON_ICONS}
        />
      )}
    </div>
  );
}

function FeatureFormModal({
  item,
  onSubmit,
  onClose,
  isLoading,
  icons,
}: {
  item: Feature | null;
  onSubmit: (data: Partial<Feature>) => void;
  onClose: () => void;
  isLoading: boolean;
  icons: { value: string; label: string }[];
}) {
  const [form, setForm] = useState({
    title: item?.title || "",
    description: item?.description || "",
    icon: item?.icon || "zap",
    icon_color: item?.icon_color || "",
    image_url: item?.image_url || "",
    image_alt: item?.image_alt || "",
    link_url: item?.link_url || "",
    link_text: item?.link_text || "",
    is_active: item?.is_active ?? true,
    sort_order: item?.sort_order || 0,
  });

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl max-h-[90vh] overflow-y-auto">
        <h3 className="text-lg font-bold text-gray-900 mb-4">
          {item ? "Edit Feature" : "Add Feature"}
        </h3>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Title *</label>
            <input
              type="text"
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Icon</label>
              <select
                value={form.icon}
                onChange={(e) => setForm({ ...form, icon: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              >
                {icons.map((icon) => (
                  <option key={icon.value} value={icon.value}>{icon.label}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Icon Color</label>
              <div className="flex gap-2">
                <input
                  type="color"
                  value={form.icon_color || "#6366f1"}
                  onChange={(e) => setForm({ ...form, icon_color: e.target.value })}
                  className="w-10 h-10 rounded border border-gray-300 cursor-pointer"
                />
                <input
                  type="text"
                  value={form.icon_color}
                  onChange={(e) => setForm({ ...form, icon_color: e.target.value })}
                  placeholder="#6366f1"
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Image URL (optional, overrides icon)</label>
            <input
              type="url"
              value={form.image_url}
              onChange={(e) => setForm({ ...form, image_url: e.target.value })}
              placeholder="https://..."
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Link URL</label>
              <input
                type="url"
                value={form.link_url}
                onChange={(e) => setForm({ ...form, link_url: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Link Text</label>
              <input
                type="text"
                value={form.link_text}
                onChange={(e) => setForm({ ...form, link_text: e.target.value })}
                placeholder="Learn more"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Sort Order</label>
              <input
                type="number"
                value={form.sort_order}
                onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div className="flex items-end pb-2">
              <label className="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                  className="w-4 h-4 text-primary rounded"
                />
                <span className="text-sm text-gray-700">Active</span>
              </label>
            </div>
          </div>
        </div>

        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50">
            Cancel
          </button>
          <button
            onClick={() => onSubmit(form)}
            disabled={isLoading || !form.title}
            className="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 disabled:opacity-50"
          >
            {isLoading ? "Saving..." : item ? "Update" : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
