"use client";

import Link from "next/link";
import { useState } from "react";
import type {
  Page,
  PageSection,
  Feature,
  Testimonial,
  PricingPlan,
  FAQ,
  NavigationMenu,
} from "@/types/api.types";

interface LandingPageProps {
  page: Page | null;
  sections: PageSection[];
  sectionContents: Record<string, Record<string, string>>;
  features: Feature[];
  testimonials: Testimonial[];
  pricingPlans: PricingPlan[];
  faqs: FAQ[];
  navigation: NavigationMenu | null;
  footerNavigation: NavigationMenu | null;
  settings: Record<string, string>;
}

// ─── Icon Map ─────────────────────────────────────────────────────────────────
const iconMap: Record<string, string> = {
  zap: "⚡", shield: "🛡️", plug: "🔌", "bar-chart": "📊", users: "👥",
  headphones: "🎧", star: "⭐", lock: "🔒", globe: "🌐", code: "💻",
  rocket: "🚀", heart: "❤️", cloud: "☁️", database: "🗄️", settings: "⚙️",
  chart: "📈", check: "✅", bolt: "⚡",
};

// ─── Default Data ─────────────────────────────────────────────────────────────
const defaultFeatures: Feature[] = [
  { id: "1", site_id: "", section_id: null, title: "Lightning Fast", description: "Sub-100ms load times with edge deployment.", icon: "zap", icon_color: "#6366f1", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 1, created_at: "", updated_at: "" },
  { id: "2", site_id: "", section_id: null, title: "SEO Optimized", description: "Built-in SEO tools, meta tags, and structured data.", icon: "bar-chart", icon_color: "#10b981", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 2, created_at: "", updated_at: "" },
  { id: "3", site_id: "", section_id: null, title: "A/B Testing", description: "Test variants automatically and publish winners.", icon: "star", icon_color: "#f59e0b", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 3, created_at: "", updated_at: "" },
  { id: "4", site_id: "", section_id: null, title: "Analytics", description: "Track conversions and revenue in real-time.", icon: "chart", icon_color: "#3b82f6", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 4, created_at: "", updated_at: "" },
  { id: "5", site_id: "", section_id: null, title: "Custom Domains", description: "Connect your domain with free SSL included.", icon: "globe", icon_color: "#8b5cf6", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 5, created_at: "", updated_at: "" },
  { id: "6", site_id: "", section_id: null, title: "Team Collaboration", description: "Role-based permissions for your entire team.", icon: "users", icon_color: "#ec4899", image_url: null, image_alt: null, link_url: null, link_text: null, is_active: true, sort_order: 6, created_at: "", updated_at: "" },
];

const defaultTestimonials: Testimonial[] = [
  { id: "1", site_id: "", section_id: null, author_name: "Sarah Johnson", author_title: "Head of Marketing", author_company: "TechCorp", author_avatar: null, content: "LandingCMS transformed our workflow. We went from weeks to hours. Conversion rate up 47% in the first month.", rating: 5, source: null, source_url: null, is_featured: true, is_active: true, sort_order: 1, created_at: "", updated_at: "" },
  { id: "2", site_id: "", section_id: null, author_name: "Michael Chen", author_title: "CTO", author_company: "StartupXYZ", author_avatar: null, content: "As a developer, I was skeptical. But the API and customization options won me over. Perfect balance of simplicity and power.", rating: 5, source: null, source_url: null, is_featured: true, is_active: true, sort_order: 2, created_at: "", updated_at: "" },
  { id: "3", site_id: "", section_id: null, author_name: "Emily Rodriguez", author_title: "Growth Lead", author_company: "ScaleUp Co.", author_avatar: null, content: "The A/B testing feature alone paid for our subscription 10x over. We optimize continuously without any developer involvement.", rating: 5, source: null, source_url: null, is_featured: true, is_active: true, sort_order: 3, created_at: "", updated_at: "" },
];

const defaultPricingPlans: PricingPlan[] = [
  { id: "1", site_id: "", section_id: null, name: "Starter", description: "For individuals", price_monthly: 9, price_yearly: 90, currency: "USD", price_label: null, is_popular: false, is_custom: false, badge_text: null, cta_text: "Get Started", cta_link: "/signup?plan=starter", features: ["1 landing page", "5K visitors/mo", "Basic analytics", "Custom domain", "SSL"], features_excluded: [], is_active: true, sort_order: 1, created_at: "", updated_at: "" },
  { id: "2", site_id: "", section_id: null, name: "Pro", description: "For growing teams", price_monthly: 29, price_yearly: 290, currency: "USD", price_label: null, is_popular: true, is_custom: false, badge_text: "Most Popular", cta_text: "Start Free Trial", cta_link: "/signup?plan=pro", features: ["Unlimited pages", "50K visitors/mo", "Advanced analytics", "A/B testing", "3 team members", "API access"], features_excluded: [], is_active: true, sort_order: 2, created_at: "", updated_at: "" },
  { id: "3", site_id: "", section_id: null, name: "Enterprise", description: "For large organizations", price_monthly: null, price_yearly: null, currency: "USD", price_label: null, is_popular: false, is_custom: true, badge_text: null, cta_text: "Contact Sales", cta_link: "/contact", features: ["Unlimited everything", "Custom analytics", "Dedicated support", "SLA guarantee", "SSO/SAML"], features_excluded: [], is_active: true, sort_order: 3, created_at: "", updated_at: "" },
];

const defaultFAQs: FAQ[] = [
  { id: "1", site_id: "", section_id: null, question: "How do I get started?", answer: "Sign up for a free account and launch your first page in under 5 minutes. No credit card required.", category: null, is_active: true, sort_order: 1, created_at: "", updated_at: "" },
  { id: "2", site_id: "", section_id: null, question: "Do I need coding skills?", answer: "Not at all! Our drag-and-drop editor is designed for non-technical users. Developers can also use our API for advanced customizations.", category: null, is_active: true, sort_order: 2, created_at: "", updated_at: "" },
  { id: "3", site_id: "", section_id: null, question: "Can I use my own domain?", answer: "Yes! Connect your custom domain on all paid plans. Free SSL certificates included. Setup takes less than 5 minutes.", category: null, is_active: true, sort_order: 3, created_at: "", updated_at: "" },
  { id: "4", site_id: "", section_id: null, question: "Is there a free trial?", answer: "Yes! All plans include a 14-day free trial with full access to all features. No credit card required.", category: null, is_active: true, sort_order: 4, created_at: "", updated_at: "" },
  { id: "5", site_id: "", section_id: null, question: "How does A/B testing work?", answer: "Create multiple page variants, split traffic automatically, and publish the winner with one click. We track all conversions for you.", category: null, is_active: true, sort_order: 5, created_at: "", updated_at: "" },
  { id: "6", site_id: "", section_id: null, question: "Is my data secure?", answer: "Absolutely. SOC 2 Type II certified, GDPR compliant, AES-256 encryption at rest and in transit.", category: null, is_active: true, sort_order: 6, created_at: "", updated_at: "" },
];

// ─── Main Component ───────────────────────────────────────────────────────────

export function LandingPage({
  page: _page,
  sections: _sections,
  sectionContents,
  features,
  testimonials,
  pricingPlans,
  faqs,
  navigation,
  footerNavigation,
  settings,
}: LandingPageProps) {
  const hero = sectionContents["hero"] || {};
  const featuresContent = sectionContents["features"] || {};
  const statsContent = sectionContents["stats"] || {};
  const testimonialsContent = sectionContents["testimonials"] || {};
  const pricingContent = sectionContents["pricing"] || {};
  const faqContent = sectionContents["faq"] || {};
  const ctaContent = sectionContents["cta"] || {};

  const displayFeatures = features.length > 0 ? features : defaultFeatures;
  const displayTestimonials = testimonials.length > 0 ? testimonials : defaultTestimonials;
  const displayPlans = pricingPlans.length > 0 ? pricingPlans : defaultPricingPlans;
  const displayFAQs = faqs.length > 0 ? faqs : defaultFAQs;

  const siteName = settings.site_title || "LandingCMS";
  const logoUrl = settings.logo_url || null;
  const logoLetter = siteName.charAt(0).toUpperCase();

  return (
    <div className="min-h-screen bg-white">
      <Header navigation={navigation} siteName={siteName} logoUrl={logoUrl} logoLetter={logoLetter} settings={settings} />
      <main>
        <HeroSection content={hero} />
        <LogosSection settings={settings} />
        <FeaturesSection content={featuresContent} features={displayFeatures} />
        <StatsSection content={statsContent} />
        <TestimonialsSection content={testimonialsContent} testimonials={displayTestimonials} />
        <PricingSection content={pricingContent} plans={displayPlans} />
        <FAQSection content={faqContent} faqs={displayFAQs} />
        <CTASection content={ctaContent} />
      </main>
      <Footer navigation={footerNavigation} siteName={siteName} settings={settings} logoUrl={logoUrl} logoLetter={logoLetter} />
    </div>
  );
}

// ─── Header ───────────────────────────────────────────────────────────────────

function Header({
  navigation,
  siteName,
  logoUrl,
  logoLetter,
  settings,
}: {
  navigation: NavigationMenu | null;
  siteName: string;
  logoUrl: string | null;
  logoLetter: string;
  settings: Record<string, string>;
}) {
  const [mobileOpen, setMobileOpen] = useState(false);

  type NavItem = { id: string; label: string; url: string | null; target: string };
  const navItems: NavItem[] = navigation?.items?.filter((i) => i.is_active && i.depth === 0).map((i) => ({
    id: i.id,
    label: i.label,
    url: i.url,
    target: i.target,
  })) || [
    { id: "1", label: "Features", url: "/#features", target: "_self" },
    { id: "2", label: "Pricing", url: "/#pricing", target: "_self" },
    { id: "3", label: "Testimonials", url: "/#testimonials", target: "_self" },
    { id: "4", label: "FAQ", url: "/#faq", target: "_self" },
  ];

  // CTA button text and link from settings
  const ctaText = settings.header_cta_text || "Get started";
  const ctaLink = settings.header_cta_link || "/#pricing";
  const signInText = settings.header_signin_text || "Sign in";

  return (
    <header className="sticky top-0 z-50 bg-white/80 backdrop-blur-xl border-b border-gray-100">
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

          {/* Desktop Nav */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => (
              <Link
                key={item.id}
                href={item.url || "#"}
                className="px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg transition-all"
              >
                {item.label}
              </Link>
            ))}
          </nav>

          {/* CTA */}
          <div className="flex items-center gap-3">
            <Link href="/admin/login" className="hidden sm:block text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors">
              {signInText}
            </Link>
            <Link
              href={ctaLink}
              className="px-4 py-2 bg-gray-900 text-white text-sm font-semibold rounded-lg hover:bg-gray-700 transition-colors shadow-sm"
            >
              {ctaText}
            </Link>
            <button
              className="md:hidden p-2 text-gray-600 hover:text-gray-900 rounded-lg hover:bg-gray-50"
              onClick={() => setMobileOpen(!mobileOpen)}
            >
              {mobileOpen ? "✕" : "☰"}
            </button>
          </div>
        </div>

        {/* Mobile Nav */}
        {mobileOpen && (
          <div className="md:hidden border-t border-gray-100 py-3 space-y-1">
            {navItems.map((item) => (
              <Link
                key={item.id}
                href={item.url || "#"}
                className="block px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50 rounded-lg"
                onClick={() => setMobileOpen(false)}
              >
                {item.label}
              </Link>
            ))}
          </div>
        )}
      </div>
    </header>
  );
}

// ─── Hero Section ─────────────────────────────────────────────────────────────

function HeroSection({ content }: { content: Record<string, string> }) {
  const title = content.title || "Build Landing Pages That Convert";
  const highlight = content.title_highlight || "Convert";
  const subtitle = content.subtitle || "Create, manage, and optimize your landing pages with our powerful CMS. No coding required. Launch in minutes.";
  const badge = content.badge_text || "🚀 Now with AI-powered suggestions";
  const ctaPrimary = content.cta_primary_text || "Start Free Trial";
  const ctaPrimaryLink = content.cta_primary_link || "#pricing";
  const ctaSecondary = content.cta_secondary_text || "See Demo";
  const ctaSecondaryLink = content.cta_secondary_link || "#demo";
  const socialProof = content.social_proof_text || "Trusted by 10,000+ teams worldwide";

  const titleParts = highlight && title.includes(highlight)
    ? title.split(highlight)
    : null;

  return (
    <section className="relative overflow-hidden bg-white pt-16 pb-24 sm:pt-24 sm:pb-32">
      {/* Background gradient */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[600px] bg-gradient-to-b from-indigo-50 via-purple-50/50 to-transparent rounded-full blur-3xl opacity-60" />
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        {/* Badge */}
        {badge && (
          <div className="inline-flex items-center gap-2 px-4 py-1.5 bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium rounded-full mb-8">
            {badge}
          </div>
        )}

        {/* Title */}
        <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold tracking-tight text-gray-900 mb-6 leading-[1.1]">
          {titleParts ? (
            <>
              {titleParts[0]}
              <span className="bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                {highlight}
              </span>
              {titleParts[1]}
            </>
          ) : (
            title
          )}
        </h1>

        {/* Subtitle */}
        <p className="text-xl text-gray-500 max-w-2xl mx-auto mb-10 leading-relaxed">
          {subtitle}
        </p>

        {/* CTAs */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
          <Link
            href={ctaPrimaryLink}
            className="inline-flex items-center justify-center px-8 py-3.5 bg-gray-900 text-white font-semibold rounded-xl hover:bg-gray-700 transition-all shadow-lg hover:shadow-xl text-base"
          >
            {ctaPrimary}
            <span className="ml-2">→</span>
          </Link>
          {ctaSecondary && (
            <Link
              href={ctaSecondaryLink}
              className="inline-flex items-center justify-center px-8 py-3.5 bg-white text-gray-700 font-semibold rounded-xl border border-gray-200 hover:border-gray-300 hover:bg-gray-50 transition-all text-base"
            >
              {ctaSecondary}
            </Link>
          )}
        </div>

        {/* Social proof */}
        {socialProof && (
          <p className="text-sm text-gray-400 font-medium">{socialProof}</p>
        )}

        {/* Hero visual - Bento grid preview */}
        <div className="mt-16 relative mx-auto max-w-5xl">
          <div className="bg-gray-900 rounded-2xl p-1 shadow-2xl">
            <div className="bg-gray-800 rounded-xl p-4">
              {/* Mock dashboard */}
              <div className="grid grid-cols-3 gap-3">
                <div className="col-span-2 bg-gray-700 rounded-lg p-4 h-32 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-3xl font-bold text-white">2.3x</div>
                    <div className="text-gray-400 text-sm mt-1">Conversion Rate</div>
                  </div>
                </div>
                <div className="bg-indigo-600 rounded-lg p-4 h-32 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-3xl font-bold text-white">47%</div>
                    <div className="text-indigo-200 text-sm mt-1">Revenue Lift</div>
                  </div>
                </div>
                <div className="bg-gray-700 rounded-lg p-4 h-20 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-xl font-bold text-white">10K+</div>
                    <div className="text-gray-400 text-xs mt-0.5">Users</div>
                  </div>
                </div>
                <div className="bg-purple-600 rounded-lg p-4 h-20 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-xl font-bold text-white">99.9%</div>
                    <div className="text-purple-200 text-xs mt-0.5">Uptime</div>
                  </div>
                </div>
                <div className="bg-gray-700 rounded-lg p-4 h-20 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-xl font-bold text-white">&lt;5min</div>
                    <div className="text-gray-400 text-xs mt-0.5">Launch</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          {/* Glow effect */}
          <div className="absolute -inset-4 bg-gradient-to-r from-indigo-500/20 to-purple-500/20 rounded-3xl blur-2xl -z-10" />
        </div>
      </div>
    </section>
  );
}

// ─── Logos Section ────────────────────────────────────────────────────────────

function LogosSection({ settings }: { settings: Record<string, string> }) {
  // Parse logos from settings (comma-separated) or use defaults
  const logosStr = settings.logos_text || "Stripe,Vercel,Notion,Linear,Figma,Loom,Intercom,Segment";
  const logos = logosStr.split(",").map((l) => l.trim()).filter(Boolean);
  const logosTitle = settings.logos_title || "Trusted by teams at";

  return (
    <section className="py-12 border-y border-gray-100 bg-gray-50/50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <p className="text-center text-sm font-medium text-gray-400 mb-8 uppercase tracking-wider">
          {logosTitle}
        </p>
        <div className="flex flex-wrap justify-center items-center gap-8 sm:gap-12">
          {logos.map((logo) => (
            <span key={logo} className="text-gray-300 font-bold text-lg tracking-tight hover:text-gray-400 transition-colors cursor-default">
              {logo}
            </span>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Features Section (Bento Grid) ───────────────────────────────────────────

function FeaturesSection({
  content,
  features,
}: {
  content: Record<string, string>;
  features: Feature[];
}) {
  const badge = content.badge_text || "Features";
  const title = content.title || "Everything you need to succeed";
  const subtitle = content.subtitle || "Powerful tools designed for marketers and developers.";

  return (
    <section id="features" className="py-24 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center max-w-2xl mx-auto mb-16">
          <span className="inline-block px-3 py-1 bg-indigo-50 text-indigo-700 text-sm font-medium rounded-full mb-4">
            {badge}
          </span>
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">{title}</h2>
          <p className="text-lg text-gray-500">{subtitle}</p>
        </div>

        {/* Bento Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {features.slice(0, 9).map((feature, index) => {
            // Make first and last items span 2 columns on large screens
            const isWide = index === 0 || index === features.length - 1;
            return (
              <div
                key={feature.id}
                className={`group relative bg-gray-50 hover:bg-white border border-gray-100 hover:border-gray-200 rounded-2xl p-6 transition-all hover:shadow-lg ${
                  isWide && features.length > 3 ? "lg:col-span-1" : ""
                }`}
              >
                {/* Icon */}
                <div
                  className="w-12 h-12 rounded-xl flex items-center justify-center text-2xl mb-4 transition-transform group-hover:scale-110"
                  style={{ backgroundColor: `${feature.icon_color || "#6366f1"}15` }}
                >
                  {(feature.icon && iconMap[feature.icon]) || "✨"}
                </div>

                {/* Content */}
                <h3 className="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
                {feature.description && (
                  <p className="text-gray-500 text-sm leading-relaxed">{feature.description}</p>
                )}

                {/* Link */}
                {feature.link_url && feature.link_text && (
                  <Link
                    href={feature.link_url}
                    className="inline-flex items-center mt-4 text-sm font-medium text-indigo-600 hover:text-indigo-700"
                  >
                    {feature.link_text}
                    <span className="ml-1">→</span>
                  </Link>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}

// ─── Stats Section ────────────────────────────────────────────────────────────

function StatsSection({ content }: { content: Record<string, string> }) {
  const stats = [
    { value: content.stat_1_value || "10,000+", label: content.stat_1_label || "Active Users" },
    { value: content.stat_2_value || "99.9%", label: content.stat_2_label || "Uptime SLA" },
    { value: content.stat_3_value || "2.3x", label: content.stat_3_label || "Avg. Conversion Lift" },
    { value: content.stat_4_value || "< 5min", label: content.stat_4_label || "Time to Launch" },
  ];

  return (
    <section className="py-16 bg-gray-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
          {stats.map((stat, i) => (
            <div key={i} className="text-center">
              <div className="text-4xl lg:text-5xl font-bold text-white mb-2 tracking-tight">
                {stat.value}
              </div>
              <div className="text-gray-400 text-sm font-medium">{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Testimonials Section ─────────────────────────────────────────────────────

function TestimonialsSection({
  content,
  testimonials,
}: {
  content: Record<string, string>;
  testimonials: Testimonial[];
}) {
  const badge = content.badge_text || "Testimonials";
  const title = content.title || "Loved by 10,000+ teams";
  const subtitle = content.subtitle || "See what our customers are saying.";

  return (
    <section id="testimonials" className="py-24 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-2xl mx-auto mb-16">
          <span className="inline-block px-3 py-1 bg-indigo-50 text-indigo-700 text-sm font-medium rounded-full mb-4">
            {badge}
          </span>
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">{title}</h2>
          <p className="text-lg text-gray-500">{subtitle}</p>
        </div>

        {/* Masonry-style grid */}
        <div className="columns-1 md:columns-2 lg:columns-3 gap-4 space-y-4">
          {testimonials.map((t) => (
            <div
              key={t.id}
              className="break-inside-avoid bg-white rounded-2xl p-6 border border-gray-100 shadow-sm hover:shadow-md transition-shadow"
            >
              {/* Stars */}
              {t.rating && (
                <div className="flex gap-0.5 mb-3">
                  {Array.from({ length: t.rating }).map((_, i) => (
                    <span key={i} className="text-amber-400 text-sm">★</span>
                  ))}
                </div>
              )}

              {/* Content */}
              <p className="text-gray-700 text-sm leading-relaxed mb-4">
                &ldquo;{t.content}&rdquo;
              </p>

              {/* Author */}
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-full bg-gradient-to-br from-indigo-400 to-purple-500 flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
                  {t.author_name.charAt(0)}
                </div>
                <div>
                  <div className="font-semibold text-gray-900 text-sm">{t.author_name}</div>
                  <div className="text-gray-400 text-xs">
                    {[t.author_title, t.author_company].filter(Boolean).join(" · ")}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── Pricing Section ──────────────────────────────────────────────────────────

function PricingSection({
  content,
  plans,
}: {
  content: Record<string, string>;
  plans: PricingPlan[];
}) {
  const [isYearly, setIsYearly] = useState(false);
  const badge = content.badge_text || "Pricing";
  const title = content.title || "Simple, transparent pricing";
  const subtitle = content.subtitle || "Start free, scale as you grow.";

  return (
    <section id="pricing" className="py-24 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-2xl mx-auto mb-12">
          <span className="inline-block px-3 py-1 bg-indigo-50 text-indigo-700 text-sm font-medium rounded-full mb-4">
            {badge}
          </span>
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">{title}</h2>
          <p className="text-lg text-gray-500 mb-8">{subtitle}</p>

          {/* Toggle */}
          <div className="inline-flex items-center gap-3 bg-gray-100 rounded-xl p-1">
            <button
              onClick={() => setIsYearly(false)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                !isYearly ? "bg-white text-gray-900 shadow-sm" : "text-gray-500"
              }`}
            >
              Monthly
            </button>
            <button
              onClick={() => setIsYearly(true)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                isYearly ? "bg-white text-gray-900 shadow-sm" : "text-gray-500"
              }`}
            >
              Yearly
              <span className="ml-1.5 text-xs text-green-600 font-semibold">-20%</span>
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {plans.map((plan) => {
            const price = isYearly ? plan.price_yearly : plan.price_monthly;
            const planFeatures = Array.isArray(plan.features)
              ? plan.features.map((f) => String(f))
              : [];

            return (
              <div
                key={plan.id}
                className={`relative rounded-2xl p-8 transition-all ${
                  plan.is_popular
                    ? "bg-gray-900 text-white shadow-2xl scale-105"
                    : "bg-gray-50 border border-gray-100 hover:border-gray-200 hover:shadow-lg"
                }`}
              >
                {plan.badge_text && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                    <span className="bg-indigo-600 text-white text-xs font-bold px-3 py-1 rounded-full">
                      {plan.badge_text}
                    </span>
                  </div>
                )}

                <div className="mb-6">
                  <h3 className={`text-xl font-bold mb-1 ${plan.is_popular ? "text-white" : "text-gray-900"}`}>
                    {plan.name}
                  </h3>
                  {plan.description && (
                    <p className={`text-sm ${plan.is_popular ? "text-gray-400" : "text-gray-500"}`}>
                      {plan.description}
                    </p>
                  )}
                </div>

                <div className="mb-6">
                  {plan.is_custom ? (
                    <div className={`text-3xl font-bold ${plan.is_popular ? "text-white" : "text-gray-900"}`}>
                      Custom
                    </div>
                  ) : (
                    <div>
                      <span className={`text-4xl font-bold ${plan.is_popular ? "text-white" : "text-gray-900"}`}>
                        ${price}
                      </span>
                      <span className={`text-sm ml-1 ${plan.is_popular ? "text-gray-400" : "text-gray-500"}`}>
                        /{isYearly ? "year" : "month"}
                      </span>
                    </div>
                  )}
                </div>

                <ul className="space-y-3 mb-8">
                  {planFeatures.map((feature, i) => (
                    <li key={i} className="flex items-start gap-2.5 text-sm">
                      <span className={`mt-0.5 flex-shrink-0 ${plan.is_popular ? "text-indigo-400" : "text-indigo-600"}`}>
                        ✓
                      </span>
                      <span className={plan.is_popular ? "text-gray-300" : "text-gray-600"}>
                        {feature}
                      </span>
                    </li>
                  ))}
                </ul>

                <Link
                  href={plan.cta_link || "#"}
                  className={`block w-full text-center py-3 px-6 rounded-xl font-semibold text-sm transition-all ${
                    plan.is_popular
                      ? "bg-white text-gray-900 hover:bg-gray-100"
                      : "bg-gray-900 text-white hover:bg-gray-700"
                  }`}
                >
                  {plan.cta_text}
                </Link>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}

// ─── FAQ Section ──────────────────────────────────────────────────────────────

function FAQSection({
  content,
  faqs,
}: {
  content: Record<string, string>;
  faqs: FAQ[];
}) {
  const [openId, setOpenId] = useState<string | null>(null);
  const badge = content.badge_text || "FAQ";
  const title = content.title || "Frequently asked questions";
  const subtitle = content.subtitle || "Everything you need to know.";

  return (
    <section id="faq" className="py-24 bg-gray-50">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-16">
          <span className="inline-block px-3 py-1 bg-indigo-50 text-indigo-700 text-sm font-medium rounded-full mb-4">
            {badge}
          </span>
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">{title}</h2>
          <p className="text-lg text-gray-500">{subtitle}</p>
        </div>

        <div className="space-y-3">
          {faqs.map((faq) => (
            <div
              key={faq.id}
              className="bg-white rounded-2xl border border-gray-100 overflow-hidden"
            >
              <button
                className="w-full text-left px-6 py-5 flex items-center justify-between gap-4 hover:bg-gray-50 transition-colors"
                onClick={() => setOpenId(openId === faq.id ? null : faq.id)}
              >
                <span className="font-semibold text-gray-900 text-sm sm:text-base">{faq.question}</span>
                <span className={`text-gray-400 text-xl flex-shrink-0 transition-transform ${openId === faq.id ? "rotate-45" : ""}`}>
                  +
                </span>
              </button>
              {openId === faq.id && (
                <div className="px-6 pb-5 text-gray-600 text-sm leading-relaxed border-t border-gray-50 pt-4">
                  {faq.answer}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ─── CTA Section ──────────────────────────────────────────────────────────────

function CTASection({ content }: { content: Record<string, string> }) {
  const title = content.title || "Ready to build your perfect landing page?";
  const subtitle = content.subtitle || "Join 10,000+ teams. Start your free trial today — no credit card required.";
  const ctaPrimary = content.cta_primary_text || "Start Free Trial";
  const ctaPrimaryLink = content.cta_primary_link || "/signup";
  const ctaSecondary = content.cta_secondary_text || "Schedule a Demo";
  const ctaSecondaryLink = content.cta_secondary_link || "/demo";

  return (
    <section className="py-24 bg-white">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative bg-gray-900 rounded-3xl p-12 sm:p-16 text-center overflow-hidden">
          {/* Background decoration */}
          <div className="absolute inset-0 -z-0">
            <div className="absolute top-0 left-1/4 w-64 h-64 bg-indigo-600/20 rounded-full blur-3xl" />
            <div className="absolute bottom-0 right-1/4 w-64 h-64 bg-purple-600/20 rounded-full blur-3xl" />
          </div>

          <div className="relative z-10">
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-white mb-4 tracking-tight">
              {title}
            </h2>
            <p className="text-gray-400 text-lg mb-10 max-w-2xl mx-auto">{subtitle}</p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Link
                href={ctaPrimaryLink}
                className="inline-flex items-center justify-center px-8 py-3.5 bg-white text-gray-900 font-semibold rounded-xl hover:bg-gray-100 transition-all shadow-lg text-base"
              >
                {ctaPrimary}
                <span className="ml-2">→</span>
              </Link>
              {ctaSecondary && (
                <Link
                  href={ctaSecondaryLink}
                  className="inline-flex items-center justify-center px-8 py-3.5 border border-gray-700 text-gray-300 font-semibold rounded-xl hover:border-gray-500 hover:text-white transition-all text-base"
                >
                  {ctaSecondary}
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ─── Footer ───────────────────────────────────────────────────────────────────

function Footer({
  navigation,
  siteName,
  settings,
  logoUrl,
  logoLetter,
}: {
  navigation: NavigationMenu | null;
  siteName: string;
  settings: Record<string, string>;
  logoUrl: string | null;
  logoLetter: string;
}) {
  const currentYear = new Date().getFullYear();
  const description = settings.site_description || "The most powerful CMS for creating dynamic landing pages.";
  const footerCopyright = settings.footer_copyright || `© ${currentYear} ${siteName}. All rights reserved.`;
  const privacyPolicyUrl = settings.privacy_policy_url || "/privacy";
  const termsUrl = settings.terms_url || "/terms";

  type FooterLink = { id: string; label: string; url: string | null };
  const footerLinks: FooterLink[] = navigation?.items?.filter((i) => i.is_active).map((i) => ({
    id: i.id,
    label: i.label,
    url: i.url,
  })) || [
    { id: "1", label: "Privacy Policy", url: privacyPolicyUrl },
    { id: "2", label: "Terms of Service", url: termsUrl },
    { id: "3", label: "Cookie Policy", url: settings.cookie_policy_url || "/cookies" },
  ];

  const socialLinks = [
    settings.social_twitter && { label: "Twitter", href: settings.social_twitter, icon: "𝕏" },
    settings.social_linkedin && { label: "LinkedIn", href: settings.social_linkedin, icon: "in" },
    settings.social_github && { label: "GitHub", href: settings.social_github, icon: "⌥" },
    settings.social_instagram && { label: "Instagram", href: settings.social_instagram, icon: "📷" },
    settings.social_youtube && { label: "YouTube", href: settings.social_youtube, icon: "▶" },
  ].filter(Boolean) as { label: string; href: string; icon: string }[];

  return (
    <footer className="bg-gray-900 text-gray-400">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-2.5 mb-4">
              {logoUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img src={logoUrl} alt={siteName} className="w-8 h-8 object-contain rounded-lg" />
              ) : (
                <div className="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-sm">{logoLetter}</span>
                </div>
              )}
              <span className="font-bold text-white text-lg">{siteName}</span>
            </div>
            <p className="text-sm leading-relaxed max-w-xs">{description}</p>
          </div>

          {/* Product */}
          <div>
            <h4 className="text-white font-semibold text-sm mb-4">Product</h4>
            <ul className="space-y-2 text-sm">
              <li><Link href="/#features" className="hover:text-white transition-colors">Features</Link></li>
              <li><Link href="/#pricing" className="hover:text-white transition-colors">Pricing</Link></li>
              <li><Link href="/#testimonials" className="hover:text-white transition-colors">Testimonials</Link></li>
              <li><Link href="/#faq" className="hover:text-white transition-colors">FAQ</Link></li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h4 className="text-white font-semibold text-sm mb-4">Legal</h4>
            <ul className="space-y-2 text-sm">
              {footerLinks.map((link) => (
                <li key={link.id}>
                  <Link href={link.url || "#"} className="hover:text-white transition-colors">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 pt-8 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-sm">{footerCopyright}</p>
          <div className="flex items-center gap-4">
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
        </div>
      </div>
    </footer>
  );
}
