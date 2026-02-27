import Link from "next/link";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface CTASectionProps {
  section: PageSection;
  contents: ContentMap;
}

export function CTASection({ section, contents }: CTASectionProps) {
  const title = getContentValue(contents, "title", "Ready to get started?");
  const subtitle = getContentValue(
    contents,
    "subtitle",
    "Join thousands of teams already using our platform. No credit card required."
  );
  const ctaPrimaryText = getContentValue(contents, "cta_primary_text", "Start Free Trial");
  const ctaPrimaryLink = getContentValue(contents, "cta_primary_link", "/signup");
  const ctaSecondaryText = getContentValue(contents, "cta_secondary_text", "Contact Sales");
  const ctaSecondaryLink = getContentValue(contents, "cta_secondary_link", "/contact");

  return (
    <section
      id={section.identifier || "cta"}
      className={`py-20 lg:py-32 bg-primary ${section.css_class || ""}`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">{title}</h2>
        {subtitle && (
          <p className="text-xl text-primary-foreground/80 mb-8 max-w-2xl mx-auto">
            {subtitle}
          </p>
        )}
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          {ctaPrimaryText && (
            <Link
              href={ctaPrimaryLink}
              className="inline-flex items-center justify-center px-8 py-3 rounded-lg bg-white text-primary font-semibold text-lg hover:bg-gray-100 transition-colors shadow-lg"
            >
              {ctaPrimaryText}
            </Link>
          )}
          {ctaSecondaryText && (
            <Link
              href={ctaSecondaryLink}
              className="inline-flex items-center justify-center px-8 py-3 rounded-lg border-2 border-white/50 text-white font-semibold text-lg hover:border-white hover:bg-white/10 transition-colors"
            >
              {ctaSecondaryText}
            </Link>
          )}
        </div>
      </div>
    </section>
  );
}
