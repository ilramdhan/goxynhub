import type { Metadata } from "next";
import { notFound, redirect } from "next/navigation";
import {
  getPublicPage,
  getPublicFeatures,
  getPublicTestimonials,
  getPublicPricingPlans,
  getPublicFAQs,
  getPublicNavigation,
  getPublicSite,
} from "@/lib/api/public.api";
import { contentsToMap } from "@/types/api.types";
import { LandingPage } from "@/components/landing/LandingPage";

const SITE_ID = process.env.NEXT_PUBLIC_DEFAULT_SITE_ID || "";

interface PageProps {
  params: { slug: string };
}

// Generate dynamic metadata from CMS
export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = params;

  // Redirect home slug to root
  if (slug === "home") {
    return { title: "Home" };
  }

  try {
    const [page, site] = await Promise.allSettled([
      getPublicPage(SITE_ID, slug),
      getPublicSite(SITE_ID),
    ]);

    const pageData = page.status === "fulfilled" ? page.value : null;
    const siteData = site.status === "fulfilled" ? site.value : null;

    const settingsMap: Record<string, string> = {};
    if (siteData?.settings) {
      siteData.settings.forEach((s) => {
        if (s.value) settingsMap[s.key] = s.value;
      });
    }

    if (!pageData) {
      return { title: "Page Not Found" };
    }

    return {
      title: pageData.seo_title || pageData.title || settingsMap.site_title || "LandingCMS",
      description: pageData.seo_description || settingsMap.seo_description || undefined,
      keywords: pageData.seo_keywords || settingsMap.seo_keywords || undefined,
      openGraph: {
        title: pageData.og_title || pageData.title || settingsMap.site_title || "LandingCMS",
        description: pageData.og_description || settingsMap.og_description || undefined,
        images: pageData.og_image ? [{ url: pageData.og_image }] : undefined,
        type: "website",
      },
      twitter: {
        card: "summary_large_image",
        title: pageData.twitter_title || pageData.og_title || pageData.title || "LandingCMS",
        description: pageData.twitter_description || pageData.og_description || undefined,
        images: pageData.twitter_image ? [pageData.twitter_image] : undefined,
      },
      alternates: {
        canonical: pageData.canonical_url || undefined,
      },
      robots: pageData.robots_meta || "index, follow",
    };
  } catch {
    return {
      title: "LandingCMS",
      description: "Build dynamic landing pages",
    };
  }
}

export default async function DynamicPage({ params }: PageProps) {
  const { slug } = params;

  // Redirect /home to root /
  if (slug === "home") {
    redirect("/");
  }

  // Fetch all data in parallel
  const [page, features, testimonials, pricingPlans, faqs, headerNav, footerNav, site] =
    await Promise.allSettled([
      getPublicPage(SITE_ID, slug),
      getPublicFeatures(SITE_ID),
      getPublicTestimonials(SITE_ID),
      getPublicPricingPlans(SITE_ID),
      getPublicFAQs(SITE_ID),
      getPublicNavigation(SITE_ID, "header"),
      getPublicNavigation(SITE_ID, "footer"),
      getPublicSite(SITE_ID),
    ]);

  const pageData = page.status === "fulfilled" ? page.value : null;

  // If page not found, return 404
  if (!pageData) {
    notFound();
  }

  const featuresData = features.status === "fulfilled" ? features.value : [];
  const testimonialsData = testimonials.status === "fulfilled" ? testimonials.value : [];
  const pricingData = pricingPlans.status === "fulfilled" ? pricingPlans.value : [];
  const faqsData = faqs.status === "fulfilled" ? faqs.value : [];
  const navData = headerNav.status === "fulfilled" ? headerNav.value : null;
  const footerNavData = footerNav.status === "fulfilled" ? footerNav.value : null;
  const siteData = site.status === "fulfilled" ? site.value : null;

  // Convert site settings to map
  const settingsMap: Record<string, string> = {};
  if (siteData?.settings) {
    siteData.settings.forEach((s) => {
      if (s.value) settingsMap[s.key] = s.value;
    });
  }
  // Also include site-level fields in settings map
  if (siteData?.logo_url) settingsMap.logo_url = siteData.logo_url;
  if (siteData?.favicon_url) settingsMap.favicon_url = siteData.favicon_url;
  if (siteData?.name && !settingsMap.site_title) settingsMap.site_title = siteData.name;
  if (siteData?.description && !settingsMap.site_description) settingsMap.site_description = siteData.description;

  // Build content maps for each section
  const sectionContents: Record<string, Record<string, string>> = {};
  if (pageData?.sections) {
    for (const section of pageData.sections) {
      if (section.identifier && section.contents) {
        const contentMap = contentsToMap(section.contents);
        sectionContents[section.identifier] = Object.fromEntries(
          Object.entries(contentMap).map(([k, v]) => [k, v.value || ""])
        );
      }
    }
  }

  return (
    <LandingPage
      page={pageData}
      sections={pageData?.sections || []}
      sectionContents={sectionContents}
      features={featuresData}
      testimonials={testimonialsData}
      pricingPlans={pricingData}
      faqs={faqsData}
      navigation={navData}
      footerNavigation={footerNavData}
      settings={settingsMap}
    />
  );
}
