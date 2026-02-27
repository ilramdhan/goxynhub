"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { listFAQs, createFAQ, updateFAQ, deleteFAQ } from "@/lib/api/components.api";
import type { FAQ } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

export default function FAQsPage() {
  const queryClient = useQueryClient();
  const [editingItem, setEditingItem] = useState<FAQ | null>(null);
  const [showForm, setShowForm] = useState(false);

  const { data: faqs = [], isLoading } = useQuery({
    queryKey: ["faqs", SITE_ID],
    queryFn: () => listFAQs(SITE_ID),
  });

  const createMutation = useMutation({
    mutationFn: createFAQ,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["faqs", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<FAQ> }) => updateFAQ(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["faqs", SITE_ID] });
      setShowForm(false);
      setEditingItem(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteFAQ,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["faqs", SITE_ID] }),
  });

  const handleSubmit = (data: Partial<FAQ>) => {
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
          <h1 className="text-2xl font-bold text-gray-900">FAQs</h1>
          <p className="text-gray-600 mt-1">Manage frequently asked questions</p>
        </div>
        <button
          onClick={() => { setEditingItem(null); setShowForm(true); }}
          className="px-4 py-2 bg-primary text-white rounded-lg font-semibold text-sm hover:bg-primary/90 transition-colors"
        >
          + Add FAQ
        </button>
      </div>

      {isLoading ? (
        <div className="text-gray-500">Loading...</div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
          {faqs.length === 0 ? (
            <div className="p-8 text-center text-gray-500">No FAQs yet. Add your first one!</div>
          ) : (
            <div className="divide-y divide-gray-100">
              {faqs.map((faq, index) => (
                <div key={faq.id} className="p-5 hover:bg-gray-50">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs text-gray-400 font-mono">#{index + 1}</span>
                        {faq.category && (
                          <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">
                            {faq.category}
                          </span>
                        )}
                        <span className={`text-xs px-2 py-0.5 rounded-full ${faq.is_active ? "bg-green-100 text-green-700" : "bg-gray-100 text-gray-600"}`}>
                          {faq.is_active ? "Active" : "Inactive"}
                        </span>
                      </div>
                      <h3 className="font-semibold text-gray-900 mb-1">{faq.question}</h3>
                      <p className="text-sm text-gray-600 line-clamp-2">{faq.answer}</p>
                    </div>
                    <div className="flex gap-2 flex-shrink-0">
                      <button
                        onClick={() => { setEditingItem(faq); setShowForm(true); }}
                        className="text-sm text-primary hover:underline"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => confirm("Delete this FAQ?") && deleteMutation.mutate(faq.id)}
                        className="text-sm text-red-600 hover:underline"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {showForm && (
        <FAQFormModal
          item={editingItem}
          onSubmit={handleSubmit}
          onClose={() => { setShowForm(false); setEditingItem(null); }}
          isLoading={createMutation.isPending || updateMutation.isPending}
        />
      )}
    </div>
  );
}

function FAQFormModal({
  item,
  onSubmit,
  onClose,
  isLoading,
}: {
  item: FAQ | null;
  onSubmit: (data: Partial<FAQ>) => void;
  onClose: () => void;
  isLoading: boolean;
}) {
  const [form, setForm] = useState({
    question: item?.question || "",
    answer: item?.answer || "",
    category: item?.category || "",
    is_active: item?.is_active ?? true,
    sort_order: item?.sort_order || 0,
  });

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl">
        <h3 className="text-lg font-bold text-gray-900 mb-4">
          {item ? "Edit FAQ" : "Add FAQ"}
        </h3>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Question *</label>
            <input
              type="text"
              value={form.question}
              onChange={(e) => setForm({ ...form, question: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Answer *</label>
            <textarea
              value={form.answer}
              onChange={(e) => setForm({ ...form, answer: e.target.value })}
              rows={5}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <input
                type="text"
                value={form.category}
                onChange={(e) => setForm({ ...form, category: e.target.value })}
                placeholder="General, Billing, etc."
                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              />
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
            <span className="text-sm text-gray-700">Active</span>
          </label>
        </div>

        <div className="flex gap-3 mt-6">
          <button onClick={onClose} className="flex-1 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50">
            Cancel
          </button>
          <button
            onClick={() => onSubmit(form)}
            disabled={isLoading || !form.question || !form.answer}
            className="flex-1 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-primary/90 disabled:opacity-50"
          >
            {isLoading ? "Saving..." : item ? "Update" : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
