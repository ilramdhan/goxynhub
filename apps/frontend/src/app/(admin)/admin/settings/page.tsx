"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient, { extractApiError } from "@/lib/api/client";
import type { ApiResponse, SiteSetting } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

async function getSettings(siteId: string): Promise<SiteSetting[]> {
  const { data } = await apiClient.get<ApiResponse<SiteSetting[]>>(
    `/admin/sites/${siteId}/settings`
  );
  return data.data ?? [];
}

async function bulkUpdateSettings(
  siteId: string,
  settings: Record<string, string>
): Promise<void> {
  await apiClient.put(`/admin/sites/${siteId}/settings`, { settings });
}

const SETTING_GROUPS = [
  { key: "general", label: "General", icon: "⚙️" },
  { key: "seo", label: "SEO", icon: "🔍" },
  { key: "social", label: "Social Media", icon: "📱" },
  { key: "analytics", label: "Analytics", icon: "📊" },
  { key: "appearance", label: "Appearance", icon: "🎨" },
];

export default function SettingsPage() {
  const queryClient = useQueryClient();
  const [activeGroup, setActiveGroup] = useState("general");
  const [formValues, setFormValues] = useState<Record<string, string>>({});
  const [isDirty, setIsDirty] = useState(false);

  const { data: settings = [], isLoading } = useQuery({
    queryKey: ["settings", SITE_ID],
    queryFn: () => getSettings(SITE_ID),
    onSuccess: (data) => {
      const values: Record<string, string> = {};
      data.forEach((s) => {
        values[s.key] = s.value || "";
      });
      setFormValues(values);
    },
  } as any);

  const saveMutation = useMutation({
    mutationFn: () => bulkUpdateSettings(SITE_ID, formValues),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["settings", SITE_ID] });
      setIsDirty(false);
    },
  });

  const handleChange = (key: string, value: string) => {
    setFormValues((prev) => ({ ...prev, [key]: value }));
    setIsDirty(true);
  };

  const groupSettings = settings.filter((s) => s.group_name === activeGroup);

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Site Settings</h1>
          <p className="text-gray-600 mt-1">Manage your site configuration and SEO settings</p>
        </div>
        {isDirty && (
          <button
            onClick={() => saveMutation.mutate()}
            disabled={saveMutation.isPending}
            className="px-6 py-2 bg-primary text-white rounded-lg font-semibold hover:bg-primary/90 disabled:opacity-50 transition-colors"
          >
            {saveMutation.isPending ? "Saving..." : "Save Changes"}
          </button>
        )}
      </div>

      <div className="flex gap-8">
        {/* Group tabs */}
        <div className="w-48 flex-shrink-0">
          <nav className="space-y-1">
            {SETTING_GROUPS.map((group) => (
              <button
                key={group.key}
                onClick={() => setActiveGroup(group.key)}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2 ${
                  activeGroup === group.key
                    ? "bg-primary/10 text-primary"
                    : "text-gray-600 hover:bg-gray-100"
                }`}
              >
                <span>{group.icon}</span>
                {group.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Settings form */}
        <div className="flex-1">
          {isLoading ? (
            <div className="text-gray-500">Loading settings...</div>
          ) : (
            <div className="bg-white rounded-xl border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-6">
                {SETTING_GROUPS.find((g) => g.key === activeGroup)?.icon}{" "}
                {SETTING_GROUPS.find((g) => g.key === activeGroup)?.label} Settings
              </h2>

              <div className="space-y-5">
                {groupSettings.map((setting) => (
                  <div key={setting.key}>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {setting.label || setting.key}
                    </label>
                    {setting.description && (
                      <p className="text-xs text-gray-500 mb-1">{setting.description}</p>
                    )}
                    {setting.type === "html" ? (
                      <textarea
                        value={formValues[setting.key] || ""}
                        onChange={(e) => handleChange(setting.key, e.target.value)}
                        rows={6}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary"
                      />
                    ) : setting.type === "boolean" ? (
                      <label className="flex items-center gap-2 cursor-pointer">
                        <input
                          type="checkbox"
                          checked={formValues[setting.key] === "true"}
                          onChange={(e) =>
                            handleChange(setting.key, e.target.checked ? "true" : "false")
                          }
                          className="w-4 h-4 text-primary rounded"
                        />
                        <span className="text-sm text-gray-700">Enabled</span>
                      </label>
                    ) : setting.type === "color" ? (
                      <div className="flex items-center gap-3">
                        <input
                          type="color"
                          value={formValues[setting.key] || "#000000"}
                          onChange={(e) => handleChange(setting.key, e.target.value)}
                          className="w-10 h-10 rounded border border-gray-300 cursor-pointer"
                        />
                        <input
                          type="text"
                          value={formValues[setting.key] || ""}
                          onChange={(e) => handleChange(setting.key, e.target.value)}
                          className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
                          placeholder="#6366f1"
                        />
                      </div>
                    ) : (
                      <input
                        type="text"
                        value={formValues[setting.key] || ""}
                        onChange={(e) => handleChange(setting.key, e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
                      />
                    )}
                  </div>
                ))}

                {groupSettings.length === 0 && (
                  <p className="text-gray-500 text-sm">No settings in this group.</p>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
