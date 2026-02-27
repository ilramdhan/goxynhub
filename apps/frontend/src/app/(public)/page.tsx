import type { Metadata } from "next";
import { getPublicPage } from "@/lib/api/pages.api";
import { contentsToMap, getContentValue } from "@/types/api.types";
import { HeroSection } from "@/components/landing/HeroSection";
import { FeaturesSection } from "@/components/landing/FeaturesSection";
import { StatsSection } from "@/components/landing/StatsSection";
import { TestimonialsSection } from "@/components/landing/TestimonialsSection";
import { PricingSection } from "@/components/landing/PricingSection";
import { FAQSection } from "@/components/landing/FAQSection";
import { CTASection } from "@/components/landing/CTASection";
import { SiteHeader } from "@/components/landing/SiteHeader";
import { SiteFooter } from "@/components/landing/SiteFooter";
import type { PageSection } from "@/types/api.types";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

// Generate dynamic metadata from CMS
export async function generateMetadata(): Promise<Metadata> {
  try {
    const page = await getPublicPage(SITE_ID);
    return {
      title: page.seo_title || page.title,
      description: page.seo_description || page.description || undefined,
      keywords: page.seo_keywords || undefined,
      openGraph: {
        title: page.og_title || page.title,
        description: page.og_description || page.description || undefined,
        images: page.og_image ? [page.og_image] : undefined,
        type: (page.og_type as "website") || "website",
      },
      twitter: {
        card: (page.twitter_card as "summary_large_image") || "summary_large_image",
        title: page.twitter_title || page.og_title || page.title,
        description: page.twitter_description || page.og_description || undefined,
        images: page.twitter_image ? [page.twitter_image] : undefined,
      },
      alternates: {
        canonical: page.canonical_url || undefined,
      },
      robots: page.robots_meta || "index, follow",
    };
  } catch {
    return {
      title: "Landing Page",
      description: "Welcome to our landing page",
    };
  }
}

// Render section based on type
function renderSection(section: PageSection) {
  if (!section.is_visible) return null;

  const contentMap = contentsToMap(section.contents || []);

  switch (section.type) {
    case "hero":
      return <HeroSection key={section.id} section={section} contents={contentMap} />;
    case "features":
      return <FeaturesSection key={section.id} section={section} contents={contentMap} />;
    case "stats":
      return <StatsSection key={section.id} section={section} contents={contentMap} />;
    case "testimonials":
      return <TestimonialsSection key={section.id} section={section} contents={contentMap} />;
    case "pricing":
      return <PricingSection key={section.id} section={section} contents={contentMap} />;
    case "faq":
      return <FAQSection key={section.id} section={section} contents={contentMap} />;
    case "cta":
      return <CTASection key={section.id} section={section} contents={contentMap} />;
    default:
      return null;
  }
}

export default async function HomePage() {
  let page = null;
  let error = null;

  try {
    page = await getPublicPage(SITE_ID);
  } catch (err) {
    error = err;
    console.error("Failed to load page:", err);
  }

  if (error || !page) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Welcome
          </h1>
          <p className="text-gray-600">
            This page is being set up. Please check back soon.
          </p>
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Custom head tags from CMS */}
      {page.custom_head && (
        <div dangerouslySetInnerHTML={{ __html: page.custom_head }} />
      )}

      <SiteHeader siteId={SITE_ID} />

      <main>
        {page.sections?.map((section) => renderSection(section))}
      </main>

      <SiteFooter siteId={SITE_ID} />
    </>
  );
}
