"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  listPricingPlans,
  createPricingPlan,
  updatePricingPlan,
  deletePricingPlan,
} from "@/lib/api/components.api";
import type { PricingPlan } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

export default function PricingPage() {
  const queryClient = useQueryClient();
  const [editingItem, setEditingItem] = useState<PricingPlan | null>(null);
  const [showForm, setShowForm] = useState(false);

  const { data: plans = [], isLoading } = useQuery({
    queryKey: ["pricing", SITE_ID],
    queryFn: () => listPricingPlans(SITE_ID),
  });

  const createMutation = useMutation({
    mutationFn: createPricingPlan,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pricing", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<PricingPlan> }) =>
      updatePricingPlan(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pricing", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deletePricingPlan,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["pricing", SITE_ID] }),
  });

  const handleSubmit = (data: Partial<PricingPlan>) => {
    if (editingItem) {
      updateMutation.mutate({ id: editingItem.id, data });
    } else {
      createMutation.mutate({ ...data, site_id: SITE_ID });
    }
  };

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Pricing Plans</h1>
          <p className="text-gray-600 mt-1">Manage your pricing tiers</p>
        </div>
        <button
          onClick={() => { setEditingItem(null); setShowForm(true); }}
          className="px-4 py-2 bg-primary text-white rounded-lg font-semibold text-sm hover:bg-primary/90 transition-colors"
        >
          + Add Plan
        </button>
      </div>

      {isLoading ? (
        <div className="text-gray-500">Loading...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {plans.map((plan) => {
            const planFeatures = Array.isArray(plan.features)
              ? plan.features.map((f) => String(f))
              : [];
            return (
              <div
                key={plan.id}
                className={`bg-white rounded-xl border p-6 ${plan.is_popular ? "border-primary" : "border-gray-200"}`}
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h3 className="font-bold text-gray-900">{plan.name}</h3>
                    {plan.badge_text && (
                      <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded-full">
                        {plan.badge_text}
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => { setEditingItem(plan); setShowForm(true); }}
                      className="text-xs text-primary hover:underline"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => confirm("Delete this plan?") && deleteMutation.mutate(plan.id)}
                      className="text-xs text-red-600 hover:underline"
                    >
                      Delete
                    </button>
                  </div>
                </div>

                <div className="mb-4">
                  {plan.is_custom ? (
                    <span className="text-2xl font-bold text-gray-900">Custom</span>
                  ) : (
                    <div>
                      <span className="text-2xl font-bold text-gray-900">
                        ${plan.price_monthly}
                      </span>
                      <span className="text-gray-500 text-sm">/mo</span>
                    </div>
                  )}
                </div>

                <ul className="space-y-1 mb-4">
                  {planFeatures.slice(0, 4).map((f, i) => (
                    <li key={i} className="text-xs text-gray-600 flex items-center gap-1">
                      <span className="text-green-500">✓</span> {f}
                    </li>
                  ))}
                  {planFeatures.length > 4 && (
                    <li className="text-xs text-gray-400">+{planFeatures.length - 4} more</li>
                  )}
                </ul>

                <div className="flex gap-2">
                  {plan.is_popular && (
                    <span className="text-xs bg-primary/10 text-primary px-2 py-0.5 rounded-full">Popular</span>
                  )}
                  <span className={`text-xs px-2 py-0.5 rounded-full ${plan.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-600"}`}>
                    {plan.is_active ? "Active" : "Inactive"}
                  </span>
                </div>
              </div>
            );
          })}
          {plans.length === 0 && (
            <div className="col-span-3 text-center py-12 text-gray-500">
              No pricing plans yet. Add your first one!
            </div>
          )}
        </div>
      )}

      {showForm && (
        <PricingFormModal
          item={editingItem}
          onSubmit={handleSubmit}
          onClose={() => { setShowForm(false); setEditingItem(null); }}
          isLoading={createMutation.isPending || updateMutation.isPending}
        />
      )}
    </div>
  );
}

function PricingFormModal({
  item,
  onSubmit,
  onClose,
  isLoading,
}: {
  item: PricingPlan | null;
  onSubmit: (data: Partial<PricingPlan>) => void;
  onClose: () => void;
  isLoading: boolean;
}) {
  const existingFeatures = Array.isArray(item?.features)
    ? item.features.map((f) => String(f)).join("\n")
    : "";

  const [form, setForm] = useState({
    name: item?.name || "",
    description: item?.description || "",
    price_monthly: item?.price_monthly?.toString() || "",
    price_yearly: item?.price_yearly?.toString() || "",
    currency: item?.currency || "USD",
    cta_text: item?.cta_text || "Get Started",
    cta_link: item?.cta_link || "",
    badge_text: item?.badge_text || "",
    is_popular: item?.is_popular ?? false,
    is_custom: item?.is_custom ?? false,
    is_active: item?.is_active ?? true,
    sort_order: item?.sort_order || 0,
    features_text: existingFeatures,
  });

  const handleSubmit = () => {
    const features = form.features_text
      .split("\n")
      .map((f) => f.trim())
      .filter(Boolean);

    onSubmit({
      name: form.name,
      description: form.description || undefined,
      price_monthly: form.price_monthly ? parseFloat(form.price_monthly) : undefined,
      price_yearly: form.price_yearly ? parseFloat(form.price_yearly) : undefined,
      currency: form.currency,
      cta_text: form.cta_text,
      cta_link: form.cta_link || undefined,
      badge_text: form.badge_text || undefined,
      is_popular: form.is_popular,
      is_custom: form.is_custom,
      is_active: form.is_active,
      sort_order: form.sort_order,
      features: features as any,
    });
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl max-h-[90vh] overflow-y-auto">
        <h3 className="text-lg font-bold text-gray-900 mb-4">
          {item ? "Edit Pricing Plan" : "Add Pricing Plan"}
        </h3>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Plan Name *</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              placeholder="Starter, Pro, Enterprise"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <input
              type="text"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Monthly Price ($)</label>
              <input
                type="number"
                value={form.price_monthly}
                onChange={(e) => setForm({ ...form, price_monthly: e.target.value })}
                placeholder="29"
                disabled={form.is_custom}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary disabled:bg-gray-50"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Yearly Price ($)</label>
              <input
                type="number"
                value={form.price_yearly}
                onChange={(e) => setForm({ ...form, price_yearly: e.target.value })}
                placeholder="290"
                disabled={form.is_custom}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary disabled:bg-gray-50"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">CTA Text</label>
              <input
                type="text"
                value={form.cta_text}
                onChange={(e) => setForm({ ...form, cta_text: e.target.value })}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">CTA Link</label>
              <input
                type="text"
                value={form.cta_link}
                onChange={(e) => setForm({ ...form, cta_link: e.target.value })}
                placeholder="/signup?plan=pro"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Badge Text</label>
            <input
              type="text"
              value={form.badge_text}
              onChange={(e) => setForm({ ...form, badge_text: e.target.value })}
              placeholder="Most Popular, Best Value"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Features (one per line)
            </label>
            <textarea
              value={form.features_text}
              onChange={(e) => setForm({ ...form, features_text: e.target.value })}
              rows={6}
              placeholder="Unlimited pages&#10;Custom domain&#10;Priority support"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div className="flex flex-wrap gap-4">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={form.is_popular}
                onChange={(e) => setForm({ ...form, is_popular: e.target.checked })}
                className="w-4 h-4 text-primary rounded"
              />
              <span className="text-sm text-gray-700">Most Popular</span>
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={form.is_custom}
                onChange={(e) => setForm({ ...form, is_custom: e.target.checked })}
                className="w-4 h-4 text-primary rounded"
              />
              <span className="text-sm text-gray-700">Custom Pricing</span>
            </label>
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

        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50">
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={isLoading || !form.name}
            className="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 disabled:opacity-50"
          >
            {isLoading ? "Saving..." : item ? "Update" : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
