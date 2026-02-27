"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import type { NavigationMenu } from "@/types/api.types";

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
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

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
    <header
      className={`sticky top-0 z-50 transition-all duration-300 ${
        scrolled
          ? "bg-white/95 backdrop-blur-md shadow-sm border-b border-gray-200/50"
          : "bg-white/80 backdrop-blur-md border-b border-gray-200/30"
      }`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2.5 group">
            <div className="w-8 h-8 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-lg flex items-center justify-center shadow-sm group-hover:shadow-md transition-shadow">
              <span className="text-white font-bold text-sm">G</span>
            </div>
            <span className="font-bold text-gray-900 text-lg tracking-tight">GrowthOS</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                target={"target" in item ? item.target : "_self"}
                className="px-4 py-2 text-gray-600 hover:text-gray-900 font-medium text-sm transition-colors rounded-lg hover:bg-gray-100"
              >
                {item.label}
              </Link>
            ))}
          </nav>

          {/* CTA */}
          <div className="flex items-center gap-3">
            <Link
              href="/admin/login"
              className="hidden sm:inline-flex text-gray-600 hover:text-gray-900 font-medium transition-colors text-sm px-3 py-2 rounded-lg hover:bg-gray-100"
            >
              Sign In
            </Link>
            <Link
              href="/#pricing"
              className="inline-flex items-center justify-center px-4 py-2 rounded-lg bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-sm hover:from-indigo-700 hover:to-purple-700 transition-all shadow-sm hover:shadow-md"
            >
              Get Started
            </Link>

            {/* Mobile menu button */}
            <button
              className="md:hidden p-2 text-gray-600 hover:text-gray-900 rounded-lg hover:bg-gray-100 transition-colors"
              onClick={() => setMobileOpen(!mobileOpen)}
              aria-label="Toggle menu"
            >
              {mobileOpen ? (
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {mobileOpen && (
          <div className="md:hidden border-t border-gray-200 py-4 animate-fade-in">
            <nav className="flex flex-col gap-1">
              {navItems.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="px-4 py-2.5 text-gray-600 hover:text-gray-900 font-medium transition-colors rounded-lg hover:bg-gray-100 text-sm"
                  onClick={() => setMobileOpen(false)}
                >
                  {item.label}
                </Link>
              ))}
              <Link
                href="/admin/login"
                className="px-4 py-2.5 text-gray-600 hover:text-gray-900 font-medium transition-colors rounded-lg hover:bg-gray-100 text-sm"
                onClick={() => setMobileOpen(false)}
              >
                Sign In
              </Link>
              <Link
                href="/#pricing"
                className="mt-2 px-4 py-2.5 bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-sm rounded-lg text-center"
                onClick={() => setMobileOpen(false)}
              >
                Get Started Free
              </Link>
            </nav>
          </div>
        )}
      </div>
    </header>
  );
}
