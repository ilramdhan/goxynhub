import type { Metadata } from "next";
import {
  getPublicPage,
  getPublicFeatures,
  getPublicTestimonials,
  getPublicPricingPlans,
  getPublicFAQs,
  getPublicNavigation,
} from "@/lib/api/public.api";
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
import type {
  PageSection,
  ContentMap,
  Feature,
  Testimonial,
  PricingPlan,
  FAQ,
  NavigationMenu,
} from "@/types/api.types";

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
        images: page.og_image ? [{ url: page.og_image }] : undefined,
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

// Props passed to each section
interface SectionProps {
  section: PageSection;
  contents: ContentMap;
  features?: Feature[];
  testimonials?: Testimonial[];
  pricingPlans?: PricingPlan[];
  faqs?: FAQ[];
}

// Render section based on type
function renderSection(props: SectionProps) {
  const { section, contents, features, testimonials, pricingPlans, faqs } = props;
  if (!section.is_visible) return null;

  switch (section.type) {
    case "hero":
      return <HeroSection key={section.id} section={section} contents={contents} />;
    case "features":
      return (
        <FeaturesSection
          key={section.id}
          section={section}
          contents={contents}
          features={features}
        />
      );
    case "stats":
      return <StatsSection key={section.id} section={section} contents={contents} />;
    case "testimonials":
      return (
        <TestimonialsSection
          key={section.id}
          section={section}
          contents={contents}
          testimonials={testimonials}
        />
      );
    case "pricing":
      return (
        <PricingSection
          key={section.id}
          section={section}
          contents={contents}
          plans={pricingPlans}
        />
      );
    case "faq":
      return (
        <FAQSection
          key={section.id}
          section={section}
          contents={contents}
          faqs={faqs}
        />
      );
    case "cta":
      return <CTASection key={section.id} section={section} contents={contents} />;
    default:
      return null;
  }
}

export default async function HomePage() {
  // Fetch all data in parallel
  const [page, features, testimonials, pricingPlans, faqs, navigation] = await Promise.allSettled([
    getPublicPage(SITE_ID),
    getPublicFeatures(SITE_ID),
    getPublicTestimonials(SITE_ID),
    getPublicPricingPlans(SITE_ID),
    getPublicFAQs(SITE_ID),
    getPublicNavigation(SITE_ID, "header"),
  ]);

  const pageData = page.status === "fulfilled" ? page.value : null;
  const featuresData = features.status === "fulfilled" ? features.value : [];
  const testimonialsData = testimonials.status === "fulfilled" ? testimonials.value : [];
  const pricingData = pricingPlans.status === "fulfilled" ? pricingPlans.value : [];
  const faqsData = faqs.status === "fulfilled" ? faqs.value : [];
  const navData = navigation.status === "fulfilled" ? navigation.value : null;

  if (!pageData) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center max-w-md mx-auto px-4">
          <div className="text-6xl mb-4">🚀</div>
          <h1 className="text-3xl font-bold text-gray-900 mb-4">Welcome</h1>
          <p className="text-gray-600 mb-6">
            This landing page is being set up. Please configure your site in the admin panel.
          </p>
          <a
            href="/admin/login"
            className="inline-flex items-center px-6 py-3 bg-primary text-white rounded-lg font-semibold hover:bg-primary/90 transition-colors"
          >
            Go to Admin Panel
          </a>
        </div>
      </div>
    );
  }

  return (
    <>
      {/* Custom head tags from CMS */}
      {pageData.custom_head && (
        <div dangerouslySetInnerHTML={{ __html: pageData.custom_head }} />
      )}

      <SiteHeader navigation={navData} />

      <main>
        {pageData.sections?.map((section) => {
          const contentMap = contentsToMap(section.contents || []);
          return renderSection({
            section,
            contents: contentMap,
            features: featuresData,
            testimonials: testimonialsData,
            pricingPlans: pricingData,
            faqs: faqsData,
          });
        })}
      </main>

      <SiteFooter siteId={SITE_ID} />
    </>
  );
}
