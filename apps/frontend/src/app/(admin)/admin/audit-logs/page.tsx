"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { listAuditLogs } from "@/lib/api/components.api";
import type { AuditLog } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

const ACTION_COLORS: Record<string, string> = {
  create: "bg-green-100 text-green-700",
  update: "bg-blue-100 text-blue-700",
  delete: "bg-red-100 text-red-700",
  login: "bg-purple-100 text-purple-700",
  logout: "bg-gray-100 text-gray-600",
  login_failed: "bg-orange-100 text-orange-700",
  publish: "bg-teal-100 text-teal-700",
  unpublish: "bg-yellow-100 text-yellow-700",
  upload: "bg-indigo-100 text-indigo-700",
  password_change: "bg-pink-100 text-pink-700",
};

export default function AuditLogsPage() {
  const [actionFilter, setActionFilter] = useState("");
  const [resourceFilter, setResourceFilter] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["audit-logs", SITE_ID, actionFilter, resourceFilter, page],
    queryFn: () =>
      listAuditLogs(SITE_ID, {
        action: actionFilter || undefined,
        resource_type: resourceFilter || undefined,
        page,
      }),
  });

  const logs = data?.data ?? [];
  const meta = data?.meta;

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">Audit Logs</h1>
        <p className="text-gray-600 mt-1">Track all admin actions and changes</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl border border-gray-200 p-4 mb-6 flex gap-4">
        <select
          value={actionFilter}
          onChange={(e) => { setActionFilter(e.target.value); setPage(1); }}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
        >
          <option value="">All Actions</option>
          <option value="create">Create</option>
          <option value="update">Update</option>
          <option value="delete">Delete</option>
          <option value="login">Login</option>
          <option value="logout">Logout</option>
          <option value="login_failed">Login Failed</option>
          <option value="publish">Publish</option>
          <option value="unpublish">Unpublish</option>
          <option value="upload">Upload</option>
          <option value="password_change">Password Change</option>
        </select>

        <select
          value={resourceFilter}
          onChange={(e) => { setResourceFilter(e.target.value); setPage(1); }}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary"
        >
          <option value="">All Resources</option>
          <option value="page">Pages</option>
          <option value="section">Sections</option>
          <option value="content">Content</option>
          <option value="site">Sites</option>
          <option value="user">Users</option>
          <option value="media">Media</option>
          <option value="feature">Features</option>
          <option value="testimonial">Testimonials</option>
          <option value="pricing_plan">Pricing Plans</option>
          <option value="faq">FAQs</option>
        </select>

        {(actionFilter || resourceFilter) && (
          <button
            onClick={() => { setActionFilter(""); setResourceFilter(""); setPage(1); }}
            className="text-sm text-gray-500 hover:text-gray-700"
          >
            Clear filters
          </button>
        )}
      </div>

      {/* Logs table */}
      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading audit logs...</div>
        ) : logs.length === 0 ? (
          <div className="p-8 text-center text-gray-500">No audit logs found.</div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase">Action</th>
                <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase">Resource</th>
                <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase">User</th>
                <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase">IP</th>
                <th className="text-left px-6 py-3 text-xs font-semibold text-gray-500 uppercase">Time</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {logs.map((log) => (
                <tr key={log.id} className="hover:bg-gray-50">
                  <td className="px-6 py-3">
                    <span
                      className={`text-xs px-2 py-1 rounded-full font-medium ${
                        ACTION_COLORS[log.action] || "bg-gray-100 text-gray-600"
                      }`}
                    >
                      {log.action}
                    </span>
                  </td>
                  <td className="px-6 py-3">
                    <div className="text-sm text-gray-900">{log.resource_type}</div>
                    {log.resource_name && (
                      <div className="text-xs text-gray-500">{log.resource_name}</div>
                    )}
                  </td>
                  <td className="px-6 py-3">
                    <div className="text-sm text-gray-900">{log.user_email || "System"}</div>
                    {log.user_role && (
                      <div className="text-xs text-gray-500">{log.user_role}</div>
                    )}
                  </td>
                  <td className="px-6 py-3 text-sm text-gray-500">
                    {log.ip_address || "—"}
                  </td>
                  <td className="px-6 py-3 text-sm text-gray-500">
                    {new Date(log.created_at).toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Pagination */}
      {meta && meta.total_pages > 1 && (
        <div className="mt-4 flex items-center justify-between text-sm text-gray-600">
          <span>
            Showing {logs.length} of {meta.total} logs
          </span>
          <div className="flex gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50"
            >
              ← Prev
            </button>
            <span className="px-3 py-1">
              Page {page} of {meta.total_pages}
            </span>
            <button
              onClick={() => setPage((p) => Math.min(meta.total_pages, p + 1))}
              disabled={page === meta.total_pages}
              className="px-3 py-1 border border-gray-300 rounded-lg disabled:opacity-50 hover:bg-gray-50"
            >
              Next →
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
