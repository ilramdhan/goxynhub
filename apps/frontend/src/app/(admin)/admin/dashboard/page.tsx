"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { listPages } from "@/lib/api/pages.api";
import { useAuthStore } from "@/stores/auth.store";

export default function DashboardPage() {
  const { user } = useAuthStore();

  const { data: pagesData } = useQuery({
    queryKey: ["pages"],
    queryFn: () => listPages({ per_page: 5 }),
  });

  const stats = [
    {
      label: "Total Pages",
      value: pagesData?.meta.total ?? 0,
      icon: "📄",
      href: "/admin/pages",
      color: "bg-blue-50 text-blue-600",
    },
    {
      label: "Published Pages",
      value: "-",
      icon: "✅",
      href: "/admin/pages?status=published",
      color: "bg-green-50 text-green-600",
    },
    {
      label: "Draft Pages",
      value: "-",
      icon: "📝",
      href: "/admin/pages?status=draft",
      color: "bg-yellow-50 text-yellow-600",
    },
    {
      label: "Media Files",
      value: "-",
      icon: "🖼️",
      href: "/admin/media",
      color: "bg-purple-50 text-purple-600",
    },
  ];

  const quickActions = [
    { label: "Create New Page", href: "/admin/pages/new", icon: "➕" },
    { label: "Upload Media", href: "/admin/media", icon: "📤" },
    { label: "Edit Navigation", href: "/admin/navigation", icon: "🧭" },
    { label: "Site Settings", href: "/admin/settings", icon: "⚙️" },
  ];

  return (
    <div className="p-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900">
          Welcome back, {user?.full_name?.split(" ")[0] || "Admin"}! 👋
        </h1>
        <p className="text-gray-600 mt-1">
          Here&apos;s an overview of your landing page CMS.
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => (
          <Link
            key={stat.label}
            href={stat.href}
            className="bg-white rounded-xl p-6 border border-gray-200 hover:shadow-md transition-shadow"
          >
            <div className={`inline-flex p-3 rounded-lg ${stat.color} mb-4`}>
              <span className="text-2xl">{stat.icon}</span>
            </div>
            <div className="text-2xl font-bold text-gray-900 mb-1">
              {stat.value}
            </div>
            <div className="text-sm text-gray-600">{stat.label}</div>
          </Link>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Quick Actions */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Quick Actions
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {quickActions.map((action) => (
              <Link
                key={action.label}
                href={action.href}
                className="flex items-center gap-3 p-3 rounded-lg border border-gray-200 hover:border-primary hover:bg-primary/5 transition-colors"
              >
                <span className="text-xl">{action.icon}</span>
                <span className="text-sm font-medium text-gray-700">
                  {action.label}
                </span>
              </Link>
            ))}
          </div>
        </div>

        {/* Recent Pages */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">
              Recent Pages
            </h2>
            <Link
              href="/admin/pages"
              className="text-sm text-primary hover:underline"
            >
              View all
            </Link>
          </div>
          <div className="space-y-3">
            {pagesData?.data.slice(0, 5).map((page) => (
              <Link
                key={page.id}
                href={`/admin/pages/${page.id}`}
                className="flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div>
                  <p className="text-sm font-medium text-gray-900">
                    {page.title}
                  </p>
                  <p className="text-xs text-gray-500">/{page.slug}</p>
                </div>
                <span
                  className={`text-xs px-2 py-1 rounded-full font-medium ${
                    page.status === "published"
                      ? "bg-green-100 text-green-700"
                      : page.status === "draft"
                      ? "bg-yellow-100 text-yellow-700"
                      : "bg-gray-100 text-gray-700"
                  }`}
                >
                  {page.status}
                </span>
              </Link>
            ))}
            {(!pagesData?.data || pagesData.data.length === 0) && (
              <p className="text-sm text-gray-500 text-center py-4">
                No pages yet.{" "}
                <Link href="/admin/pages/new" className="text-primary hover:underline">
                  Create your first page
                </Link>
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
