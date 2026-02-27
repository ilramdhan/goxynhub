"use client";

import { useState } from "react";
import Link from "next/link";
import type { NavigationMenu, NavigationItem } from "@/types/api.types";

interface SiteHeaderProps {
  navigation?: NavigationMenu | null;
}

const defaultNavItems = [
  { label: "Features", href: "/#features" },
  { label: "Pricing", href: "/#pricing" },
  { label: "FAQ", href: "/#faq" },
];

export function SiteHeader({ navigation }: SiteHeaderProps) {
  const [mobileOpen, setMobileOpen] = useState(false);

  // Use navigation from CMS if available, otherwise use defaults
  const navItems =
    navigation?.items && navigation.items.length > 0
      ? navigation.items
          .filter((item) => item.is_active && item.depth === 0)
          .map((item) => ({
            label: item.label,
            href: item.url || "#",
            target: item.target,
          }))
      : defaultNavItems;

  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-gray-200">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">L</span>
            </div>
            <span className="font-bold text-gray-900 text-lg">LandingCMS</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-8">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                target={"target" in item ? item.target : "_self"}
                className="text-gray-600 hover:text-gray-900 font-medium transition-colors"
              >
                {item.label}
              </Link>
            ))}
          </nav>

          {/* CTA */}
          <div className="flex items-center gap-4">
            <Link
              href="/admin/login"
              className="hidden sm:inline-flex text-gray-600 hover:text-gray-900 font-medium transition-colors text-sm"
            >
              Sign In
            </Link>
            <Link
              href="/#pricing"
              className="inline-flex items-center justify-center px-4 py-2 rounded-lg bg-primary text-white font-semibold text-sm hover:bg-primary/90 transition-colors"
            >
              Get Started
            </Link>

            {/* Mobile menu button */}
            <button
              className="md:hidden p-2 text-gray-600 hover:text-gray-900"
              onClick={() => setMobileOpen(!mobileOpen)}
            >
              {mobileOpen ? "✕" : "☰"}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileOpen && (
          <div className="md:hidden border-t border-gray-200 py-4">
            <nav className="flex flex-col gap-3">
              {navItems.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="text-gray-600 hover:text-gray-900 font-medium transition-colors py-1"
                  onClick={() => setMobileOpen(false)}
                >
                  {item.label}
                </Link>
              ))}
              <Link
                href="/admin/login"
                className="text-gray-600 hover:text-gray-900 font-medium transition-colors py-1"
                onClick={() => setMobileOpen(false)}
              >
                Sign In
              </Link>
            </nav>
          </div>
        )}
      </div>
    </header>
  );
}
