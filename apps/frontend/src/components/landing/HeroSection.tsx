import React from "react";
import Image from "next/image";
import Link from "next/link";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface HeroSectionProps {
  section: PageSection;
  contents: ContentMap;
}

export function HeroSection({ section, contents }: HeroSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Build Something Amazing");
  const titleHighlight = getContentValue(contents, "title_highlight");
  const subtitle = getContentValue(
    contents,
    "subtitle",
    "The most powerful platform to create, manage, and scale your digital presence."
  );
  const ctaPrimaryText = getContentValue(contents, "cta_primary_text", "Get Started Free");
  const ctaPrimaryLink = getContentValue(contents, "cta_primary_link", "#");
  const ctaSecondaryText = getContentValue(contents, "cta_secondary_text");
  const ctaSecondaryLink = getContentValue(contents, "cta_secondary_link", "#");
  const heroImage = getContentValue(contents, "hero_image");
  const heroImageAlt = getContentValue(contents, "hero_image_alt", "Product screenshot");
  const socialProofText = getContentValue(contents, "social_proof_text");

  // Build title with highlight
  let titleParts: { text: string; highlight: boolean }[] = [];
  if (titleHighlight && title.includes(titleHighlight)) {
    const parts = title.split(titleHighlight);
    titleParts = [
      { text: parts[0], highlight: false },
      { text: titleHighlight, highlight: true },
      { text: parts[1] || "", highlight: false },
    ];
  } else {
    titleParts = [{ text: title, highlight: false }];
  }

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;
  if (section.bg_image) {
    sectionStyle.backgroundImage = `url(${section.bg_image})`;
    sectionStyle.backgroundSize = "cover";
    sectionStyle.backgroundPosition = "center";
  }

  return (
    <section
      id={section.identifier || "hero"}
      className={`relative overflow-hidden py-20 lg:py-32 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      {/* Background overlay */}
      {section.bg_overlay && (
        <div
          className="absolute inset-0"
          style={{
            backgroundColor: section.bg_overlay_color || "rgba(0,0,0,0.5)",
            opacity: section.bg_overlay_opacity || 0.5,
          }}
        />
      )}

      <div className="relative container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-4xl mx-auto">
          {/* Badge */}
          {badgeText && (
            <div className="inline-flex items-center px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-6 animate-fade-in">
              {badgeText}
            </div>
          )}

          {/* Title */}
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-gray-900 mb-6 animate-fade-in">
            {titleParts.map((part, i) =>
              part.highlight ? (
                <span key={i} className="text-primary">
                  {part.text}
                </span>
              ) : (
                <span key={i}>{part.text}</span>
              )
            )}
          </h1>

          {/* Subtitle */}
          {subtitle && (
            <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto animate-fade-in">
              {subtitle}
            </p>
          )}

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12 animate-fade-in">
            {ctaPrimaryText && (
              <Link
                href={ctaPrimaryLink}
                className="inline-flex items-center justify-center px-8 py-3 rounded-lg bg-primary text-white font-semibold text-lg hover:bg-primary/90 transition-colors shadow-lg hover:shadow-xl"
              >
                {ctaPrimaryText}
              </Link>
            )}
            {ctaSecondaryText && (
              <Link
                href={ctaSecondaryLink}
                className="inline-flex items-center justify-center px-8 py-3 rounded-lg border-2 border-gray-300 text-gray-700 font-semibold text-lg hover:border-primary hover:text-primary transition-colors"
              >
                {ctaSecondaryText}
              </Link>
            )}
          </div>

          {/* Social proof */}
          {socialProofText && (
            <p className="text-sm text-gray-500 animate-fade-in">
              {socialProofText}
            </p>
          )}
        </div>

        {/* Hero image */}
        {heroImage && (
          <div className="mt-16 relative mx-auto max-w-5xl animate-fade-in">
            <div className="rounded-2xl overflow-hidden shadow-2xl border border-gray-200">
              <Image
                src={heroImage}
                alt={heroImageAlt}
                width={1200}
                height={675}
                className="w-full h-auto"
                priority
              />
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
