"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listNavigation } from "@/lib/api/components.api";
import apiClient from "@/lib/api/client";
import type { NavigationMenu, NavigationItem } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

export default function NavigationPage() {
  const queryClient = useQueryClient();
  const [selectedMenu, setSelectedMenu] = useState<NavigationMenu | null>(null);
  const [showAddItem, setShowAddItem] = useState(false);
  const [editingItem, setEditingItem] = useState<NavigationItem | null>(null);

  const { data: menus = [], isLoading } = useQuery({
    queryKey: ["navigation", SITE_ID],
    queryFn: () => listNavigation(SITE_ID),
    onSuccess: (data) => {
      if (data.length > 0 && !selectedMenu) {
        setSelectedMenu(data[0]);
      }
    },
  } as any);

  const createItemMutation = useMutation({
    mutationFn: (data: Partial<NavigationItem>) =>
      apiClient.post(`/admin/navigation/${selectedMenu?.id}/items`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["navigation", SITE_ID] });
      setShowAddItem(false);
    },
  });

  const updateItemMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<NavigationItem> }) =>
      apiClient.put(`/admin/navigation/items/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["navigation", SITE_ID] });
      setEditingItem(null);
    },
  });

  const deleteItemMutation = useMutation({
    mutationFn: (id: string) => apiClient.delete(`/admin/navigation/items/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["navigation", SITE_ID] }),
  });

  const currentMenu = menus.find((m) => m.id === selectedMenu?.id);
  const items = currentMenu?.items || [];

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Navigation</h1>
          <p className="text-gray-600 mt-1">Manage your site navigation menus</p>
        </div>
        {selectedMenu && (
          <button
            onClick={() => setShowAddItem(true)}
            className="px-4 py-2 bg-primary text-white rounded-lg font-semibold text-sm hover:bg-primary/90 transition-colors"
          >
            + Add Menu Item
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="text-gray-500">Loading...</div>
      ) : (
        <div className="flex gap-6">
          {/* Menu selector */}
          <div className="w-48 flex-shrink-0">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">
              Menus
            </h3>
            <div className="space-y-1">
              {menus.map((menu) => (
                <button
                  key={menu.id}
                  onClick={() => setSelectedMenu(menu)}
                  className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-colors ${
                    selectedMenu?.id === menu.id
                      ? "bg-primary/10 text-primary font-medium"
                      : "text-gray-600 hover:bg-gray-100"
                  }`}
                >
                  {menu.name}
                  <span className="block text-xs text-gray-400">{menu.identifier}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Menu items */}
          <div className="flex-1">
            {selectedMenu ? (
              <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-200 bg-gray-50">
                  <h3 className="font-semibold text-gray-900">{selectedMenu.name}</h3>
                  <p className="text-sm text-gray-500">identifier: {selectedMenu.identifier}</p>
                </div>

                {items.length === 0 ? (
                  <div className="p-8 text-center text-gray-500">
                    No items in this menu.{" "}
                    <button
                      onClick={() => setShowAddItem(true)}
                      className="text-primary hover:underline"
                    >
                      Add first item
                    </button>
                  </div>
                ) : (
                  <div className="divide-y divide-gray-100">
                    {renderNavItems(items, (item) => {
                      setEditingItem(item);
                    }, (id) => {
                      if (confirm("Delete this menu item?")) {
                        deleteItemMutation.mutate(id);
                      }
                    })}
                  </div>
                )}
              </div>
            ) : (
              <div className="text-center py-12 text-gray-500">
                Select a menu to manage its items
              </div>
            )}
          </div>
        </div>
      )}

      {/* Add/Edit Item Modal */}
      {(showAddItem || editingItem) && (
        <NavItemFormModal
          item={editingItem}
          onSubmit={(data) => {
            if (editingItem) {
              updateItemMutation.mutate({ id: editingItem.id, data });
            } else {
              createItemMutation.mutate(data);
            }
          }}
          onClose={() => { setShowAddItem(false); setEditingItem(null); }}
          isLoading={createItemMutation.isPending || updateItemMutation.isPending}
        />
      )}
    </div>
  );
}

function renderNavItems(
  items: NavigationItem[],
  onEdit: (item: NavigationItem) => void,
  onDelete: (id: string) => void,
  depth = 0
): React.ReactNode {
  return items.map((item) => (
    <div key={item.id}>
      <div
        className="flex items-center justify-between px-6 py-3 hover:bg-gray-50"
        style={{ paddingLeft: `${24 + depth * 24}px` }}
      >
        <div className="flex items-center gap-3">
          {depth > 0 && <span className="text-gray-300">└</span>}
          <div>
            <span className="font-medium text-gray-900 text-sm">{item.label}</span>
            {item.url && (
              <span className="ml-2 text-xs text-gray-400">{item.url}</span>
            )}
          </div>
          {!item.is_active && (
            <span className="text-xs bg-gray-100 text-gray-500 px-2 py-0.5 rounded-full">Hidden</span>
          )}
          {item.target === "_blank" && (
            <span className="text-xs bg-blue-100 text-blue-600 px-2 py-0.5 rounded-full">External</span>
          )}
        </div>
        <div className="flex gap-2">
          <button onClick={() => onEdit(item)} className="text-xs text-primary hover:underline">
            Edit
          </button>
          <button onClick={() => onDelete(item.id)} className="text-xs text-red-600 hover:underline">
            Delete
          </button>
        </div>
      </div>
      {item.children && item.children.length > 0 &&
        renderNavItems(item.children, onEdit, onDelete, depth + 1)}
    </div>
  ));
}

function NavItemFormModal({
  item,
  onSubmit,
  onClose,
  isLoading,
}: {
  item: NavigationItem | null;
  onSubmit: (data: Partial<NavigationItem>) => void;
  onClose: () => void;
  isLoading: boolean;
}) {
  const [form, setForm] = useState({
    label: item?.label || "",
    url: item?.url || "",
    target: item?.target || "_self",
    icon: item?.icon || "",
    is_active: item?.is_active ?? true,
    sort_order: item?.sort_order || 0,
  });

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl">
        <h3 className="text-lg font-bold text-gray-900 mb-4">
          {item ? "Edit Menu Item" : "Add Menu Item"}
        </h3>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Label *</label>
            <input
              type="text"
              value={form.label}
              onChange={(e) => setForm({ ...form, label: e.target.value })}
              placeholder="Home, About, Contact"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">URL</label>
            <input
              type="text"
              value={form.url}
              onChange={(e) => setForm({ ...form, url: e.target.value })}
              placeholder="/, /about, /#features, https://..."
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Target</label>
              <select
                value={form.target}
                onChange={(e) => setForm({ ...form, target: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              >
                <option value="_self">Same tab</option>
                <option value="_blank">New tab</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Sort Order</label>
              <input
                type="number"
                value={form.sort_order}
                onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={form.is_active}
              onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
              className="w-4 h-4 text-primary rounded"
            />
            <span className="text-sm text-gray-700">Visible</span>
          </label>
        </div>

        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50">
            Cancel
          </button>
          <button
            onClick={() => onSubmit(form)}
            disabled={isLoading || !form.label}
            className="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 disabled:opacity-50"
          >
            {isLoading ? "Saving..." : item ? "Update" : "Add"}
          </button>
        </div>
      </div>
    </div>
  );
}
