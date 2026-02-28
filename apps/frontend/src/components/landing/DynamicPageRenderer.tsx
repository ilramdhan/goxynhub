"use client";

import type { ReactNode } from "react";
import Link from "next/link";
import Image from "next/image";
import type {
  Page,
  PageSection,
  SectionContent,
  Feature,
  Testimonial,
  PricingPlan,
  FAQ,
  NavigationMenu,
} from "@/types/api.types";

interface DynamicPageRendererProps {
  page: Page;
  sections: PageSection[];
  features: Feature[];
  testimonials: Testimonial[];
  pricingPlans: PricingPlan[];
  faqs: FAQ[];
  navigation: NavigationMenu | null;
  footerNavigation: NavigationMenu | null;
  settings: Record<string, string>;
}

// Helper to get content value from section contents
function getContent(contents: SectionContent[] | undefined, key: string, defaultValue = ""): string {
  if (!contents) return defaultValue;
  const item = contents.find((c) => c.key === key);
  return item?.value || defaultValue;
}

// Render a single section based on its type
function renderSection(section: PageSection, _settings: Record<string, string>): ReactNode {
  if (!section.is_visible) return null;

  const contents = section.contents || [];

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;
  if (section.bg_image) {
    sectionStyle.backgroundImage = `url(${section.bg_image})`;
    sectionStyle.backgroundSize = "cover";
    sectionStyle.backgroundPosition = "center";
  }
  if (section.padding_top) sectionStyle.paddingTop = section.padding_top;
  if (section.padding_bottom) sectionStyle.paddingBottom = section.padding_bottom;

  switch (section.type) {
    case "hero":
      return (
        <section
          key={section.id}
          id={section.identifier || "hero"}
          className={`relative py-20 lg:py-32 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          {!section.bg_color && !section.bg_image && (
            <div className="absolute inset-0 bg-gradient-to-br from-indigo-50 via-white to-purple-50 -z-10" />
          )}
          {section.bg_overlay && (
            <div
              className="absolute inset-0"
              style={{
                backgroundColor: section.bg_overlay_color || "rgba(0,0,0,0.5)",
                opacity: section.bg_overlay_opacity || 0.5,
              }}
            />
          )}
          <div className="relative container mx-auto px-4 sm:px-6 lg:px-8 text-center max-w-4xl">
            {getContent(contents, "badge_text") && (
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium mb-8">
                {getContent(contents, "badge_text")}
              </div>
            )}
            <h1 className="text-5xl sm:text-6xl font-bold tracking-tight text-gray-900 mb-6">
              {getContent(contents, "title", "Welcome")}
            </h1>
            {getContent(contents, "subtitle") && (
              <p className="text-xl text-gray-600 mb-10 max-w-2xl mx-auto">
                {getContent(contents, "subtitle")}
              </p>
            )}
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              {getContent(contents, "cta_primary_text") && (
                <Link
                  href={getContent(contents, "cta_primary_link", "#")}
                  className="inline-flex items-center justify-center px-8 py-4 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-lg hover:from-indigo-700 hover:to-purple-700 transition-all shadow-lg"
                >
                  {getContent(contents, "cta_primary_text")}
                </Link>
              )}
              {getContent(contents, "cta_secondary_text") && (
                <Link
                  href={getContent(contents, "cta_secondary_link", "#")}
                  className="inline-flex items-center justify-center px-8 py-4 rounded-xl border-2 border-gray-200 text-gray-700 font-semibold text-lg hover:border-indigo-300 hover:text-indigo-700 transition-all"
                >
                  {getContent(contents, "cta_secondary_text")}
                </Link>
              )}
            </div>
            {getContent(contents, "hero_image") && (
              <div className="mt-16 relative mx-auto max-w-5xl">
                <Image
                  src={getContent(contents, "hero_image")}
                  alt={getContent(contents, "hero_image_alt", "Hero image")}
                  width={1200}
                  height={675}
                  className="w-full h-auto rounded-2xl shadow-2xl"
                />
              </div>
            )}
          </div>
        </section>
      );

    case "html":
    case "custom":
      return (
        <section
          key={section.id}
          id={section.identifier || section.id}
          className={`py-16 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            {getContent(contents, "html") && (
              <div
                className="prose prose-lg max-w-none"
                dangerouslySetInnerHTML={{ __html: getContent(contents, "html") }}
              />
            )}
            {getContent(contents, "content") && (
              <div
                className="prose prose-lg max-w-none"
                dangerouslySetInnerHTML={{ __html: getContent(contents, "content") }}
              />
            )}
          </div>
        </section>
      );

    case "about":
      return (
        <section
          key={section.id}
          id={section.identifier || "about"}
          className={`py-20 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto">
              {getContent(contents, "badge_text") && (
                <div className="inline-flex items-center px-4 py-1.5 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium mb-6">
                  {getContent(contents, "badge_text")}
                </div>
              )}
              <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-6">
                {getContent(contents, "title", "About Us")}
              </h2>
              {getContent(contents, "subtitle") && (
                <p className="text-xl text-gray-600 mb-8">
                  {getContent(contents, "subtitle")}
                </p>
              )}
              {getContent(contents, "content") && (
                <div
                  className="prose prose-lg max-w-none text-gray-700"
                  dangerouslySetInnerHTML={{ __html: getContent(contents, "content") }}
                />
              )}
              {getContent(contents, "body") && (
                <div className="text-gray-700 leading-relaxed whitespace-pre-wrap">
                  {getContent(contents, "body")}
                </div>
              )}
              {getContent(contents, "image") && (
                <div className="mt-10">
                  <Image
                    src={getContent(contents, "image")}
                    alt={getContent(contents, "image_alt", "About image")}
                    width={800}
                    height={500}
                    className="w-full h-auto rounded-2xl shadow-lg"
                  />
                </div>
              )}
            </div>
          </div>
        </section>
      );

    case "contact":
      return (
        <section
          key={section.id}
          id={section.identifier || "contact"}
          className={`py-20 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-2xl mx-auto text-center">
              <h2 className="text-4xl font-bold text-gray-900 mb-4">
                {getContent(contents, "title", "Contact Us")}
              </h2>
              {getContent(contents, "subtitle") && (
                <p className="text-xl text-gray-600 mb-8">
                  {getContent(contents, "subtitle")}
                </p>
              )}
              {getContent(contents, "email") && (
                <a
                  href={`mailto:${getContent(contents, "email")}`}
                  className="text-indigo-600 hover:text-indigo-700 text-lg font-medium"
                >
                  {getContent(contents, "email")}
                </a>
              )}
              {getContent(contents, "content") && (
                <div
                  className="prose prose-lg max-w-none mt-8 text-left"
                  dangerouslySetInnerHTML={{ __html: getContent(contents, "content") }}
                />
              )}
            </div>
          </div>
        </section>
      );

    case "cta":
      return (
        <section
          key={section.id}
          id={section.identifier || "cta"}
          className={`py-20 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="bg-gray-900 rounded-3xl p-12 text-center">
              <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
                {getContent(contents, "title", "Ready to get started?")}
              </h2>
              {getContent(contents, "subtitle") && (
                <p className="text-gray-400 text-lg mb-8">
                  {getContent(contents, "subtitle")}
                </p>
              )}
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                {getContent(contents, "cta_primary_text") && (
                  <Link
                    href={getContent(contents, "cta_primary_link", "#")}
                    className="inline-flex items-center justify-center px-8 py-3.5 bg-white text-gray-900 font-semibold rounded-xl hover:bg-gray-100 transition-all"
                  >
                    {getContent(contents, "cta_primary_text")}
                  </Link>
                )}
                {getContent(contents, "cta_secondary_text") && (
                  <Link
                    href={getContent(contents, "cta_secondary_link", "#")}
                    className="inline-flex items-center justify-center px-8 py-3.5 border border-gray-700 text-gray-300 font-semibold rounded-xl hover:border-gray-500 hover:text-white transition-all"
                  >
                    {getContent(contents, "cta_secondary_text")}
                  </Link>
                )}
              </div>
            </div>
          </div>
        </section>
      );

    case "stats": {
      const statsData = [
        { value: getContent(contents, "stat_1_value"), label: getContent(contents, "stat_1_label") },
        { value: getContent(contents, "stat_2_value"), label: getContent(contents, "stat_2_label") },
        { value: getContent(contents, "stat_3_value"), label: getContent(contents, "stat_3_label") },
        { value: getContent(contents, "stat_4_value"), label: getContent(contents, "stat_4_label") },
      ].filter((s) => s.value);

      return (
        <section
          key={section.id}
          id={section.identifier || "stats"}
          className={`py-16 bg-gray-900 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
              {statsData.map((stat, i) => (
                <div key={i} className="text-center">
                  <div className="text-4xl font-bold text-white mb-2">{stat.value}</div>
                  <div className="text-gray-400 text-sm">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </section>
      );
    }

    case "video":
      return (
        <section
          key={section.id}
          id={section.identifier || "video"}
          className={`py-20 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            <div className="max-w-4xl mx-auto text-center">
              {getContent(contents, "title") && (
                <h2 className="text-4xl font-bold text-gray-900 mb-4">
                  {getContent(contents, "title")}
                </h2>
              )}
              {getContent(contents, "subtitle") && (
                <p className="text-xl text-gray-600 mb-8">
                  {getContent(contents, "subtitle")}
                </p>
              )}
              {getContent(contents, "video_url") && (
                <div className="relative aspect-video rounded-2xl overflow-hidden shadow-2xl">
                  <iframe
                    src={getContent(contents, "video_url")}
                    className="w-full h-full"
                    allowFullScreen
                    title={getContent(contents, "title", "Video")}
                  />
                </div>
              )}
            </div>
          </div>
        </section>
      );

    default:
      // Generic section renderer for any other type
      return (
        <section
          key={section.id}
          id={section.identifier || section.id}
          className={`py-16 ${section.css_class || ""}`}
          style={sectionStyle}
        >
          <div className="container mx-auto px-4 sm:px-6 lg:px-8">
            {getContent(contents, "title") && (
              <h2 className="text-3xl font-bold text-gray-900 mb-4">
                {getContent(contents, "title")}
              </h2>
            )}
            {getContent(contents, "subtitle") && (
              <p className="text-xl text-gray-600 mb-8">
                {getContent(contents, "subtitle")}
              </p>
            )}
            {getContent(contents, "content") && (
              <div
                className="prose prose-lg max-w-none"
                dangerouslySetInnerHTML={{ __html: getContent(contents, "content") }}
              />
            )}
            {getContent(contents, "body") && (
              <div className="text-gray-700 leading-relaxed whitespace-pre-wrap">
                {getContent(contents, "body")}
              </div>
            )}
          </div>
        </section>
      );
  }
}

export function DynamicPageRenderer({
  page,
  sections,
  features: _features,
  testimonials: _testimonials,
  pricingPlans: _pricingPlans,
  faqs: _faqs,
  navigation,
  footerNavigation,
  settings,
}: DynamicPageRendererProps) {
  const siteName = settings.site_title || "LandingCMS";
  const logoUrl = settings.logo_url || null;
  const logoLetter = siteName.charAt(0).toUpperCase();
  const currentYear = new Date().getFullYear();
  const footerCopyright = settings.footer_copyright || `© ${currentYear} ${siteName}. All rights reserved.`;

  type NavItem = { id: string; label: string; url: string | null; target: string };
  const navItems: NavItem[] = navigation?.items?.filter((i) => i.is_active && i.depth === 0).map((i) => ({
    id: i.id,
    label: i.label,
    url: i.url,
    target: i.target,
  })) || [];

  type FooterLink = { id: string; label: string; url: string | null };
  const footerLinks: FooterLink[] = footerNavigation?.items?.filter((i) => i.is_active).map((i) => ({
    id: i.id,
    label: i.label,
    url: i.url,
  })) || [
    { id: "1", label: "Privacy Policy", url: settings.privacy_policy_url || "/privacy" },
    { id: "2", label: "Terms of Service", url: settings.terms_url || "/terms" },
  ];

  const socialLinks = [
    settings.social_twitter && { label: "Twitter", href: settings.social_twitter, icon: "𝕏" },
    settings.social_linkedin && { label: "LinkedIn", href: settings.social_linkedin, icon: "in" },
    settings.social_github && { label: "GitHub", href: settings.social_github, icon: "⌥" },
  ].filter(Boolean) as { label: string; href: string; icon: string }[];

  const visibleSections = sections.filter((s) => s.is_visible);

  return (
    <div className="min-h-screen bg-white flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-50 bg-white/90 backdrop-blur-xl border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link href="/" className="flex items-center gap-2.5">
              {logoUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={logoUrl} alt={siteName} className="w-8 h-8 object-contain rounded-lg" />
              ) : (
                <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center shadow-sm">
                  <span className="text-white font-bold text-sm">{logoLetter}</span>
                </div>
              )}
              <span className="font-bold text-gray-900 text-lg tracking-tight">{siteName}</span>
            </Link>

            {/* Navigation */}
            {navItems.length > 0 && (
              <nav className="hidden md:flex items-center gap-1">
                {navItems.map((item) => (
                  <Link
                    key={item.id}
                    href={item.url || "#"}
                    target={item.target === "_blank" ? "_blank" : undefined}
                    rel={item.target === "_blank" ? "noopener noreferrer" : undefined}
                    className="px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-all"
                  >
                    {item.label}
                  </Link>
                ))}
              </nav>
            )}

            {/* Back to home */}
            <Link
              href="/"
              className="text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors"
            >
              ← Home
            </Link>
          </div>
        </div>
      </header>

      {/* Page title breadcrumb */}
      <div className="bg-gray-50 border-b border-gray-100 py-4">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-2 text-sm text-gray-500">
            <Link href="/" className="hover:text-gray-700">Home</Link>
            <span>/</span>
            <span className="text-gray-900 font-medium">{page.title}</span>
          </div>
        </div>
      </div>

      {/* Main content */}
      <main className="flex-1">
        {visibleSections.length > 0 ? (
          visibleSections.map((section) => renderSection(section, settings))
        ) : (
          // Fallback: show page title and description if no sections
          <div className="py-20">
            <div className="container mx-auto px-4 sm:px-6 lg:px-8 max-w-4xl">
              <h1 className="text-4xl font-bold text-gray-900 mb-6">{page.title}</h1>
              {page.description && (
                <p className="text-xl text-gray-600">{page.description}</p>
              )}
            </div>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex flex-col md:flex-row items-center justify-between gap-6">
            {/* Logo */}
            <Link href="/" className="flex items-center gap-2.5">
              {logoUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={logoUrl} alt={siteName} className="w-6 h-6 object-contain rounded" />
              ) : (
                <div className="w-6 h-6 bg-gradient-to-br from-indigo-500 to-purple-600 rounded flex items-center justify-center">
                  <span className="text-white font-bold text-xs">{logoLetter}</span>
                </div>
              )}
              <span className="font-bold text-white text-sm">{siteName}</span>
            </Link>

            {/* Footer links */}
            <nav className="flex flex-wrap justify-center gap-x-6 gap-y-2">
              {footerLinks.map((link) => (
                <Link
                  key={link.id}
                  href={link.url || "#"}
                  className="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  {link.label}
                </Link>
              ))}
            </nav>

            {/* Social links */}
            {socialLinks.length > 0 && (
              <div className="flex items-center gap-3">
                {socialLinks.map((link) => (
                  <a
                    key={link.label}
                    href={link.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="w-8 h-8 bg-gray-800 hover:bg-gray-700 rounded-lg flex items-center justify-center text-sm transition-colors"
                    aria-label={link.label}
                  >
                    {link.icon}
                  </a>
                ))}
              </div>
            )}
          </div>

          <div className="mt-6 pt-6 border-t border-gray-800 text-center">
            <p className="text-sm text-gray-500">{footerCopyright}</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
