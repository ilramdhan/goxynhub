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
  const title = getContentValue(contents, "title", "Scale Your Business with Confidence");
  const titleHighlight = getContentValue(contents, "title_highlight");
  const subtitle = getContentValue(
    contents,
    "subtitle",
    "The all-in-one platform that helps startups and enterprises grow faster with powerful analytics, automation, and team collaboration tools."
  );
  const ctaPrimaryText = getContentValue(contents, "cta_primary_text", "Start Free Trial");
  const ctaPrimaryLink = getContentValue(contents, "cta_primary_link", "#pricing");
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
      className={`relative overflow-hidden py-24 lg:py-32 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      {/* Background gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-indigo-50 via-white to-purple-50 -z-10" />
      
      {/* Decorative blobs */}
      <div className="absolute top-0 left-1/4 w-96 h-96 bg-indigo-100 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob -z-10" />
      <div className="absolute top-0 right-1/4 w-96 h-96 bg-purple-100 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob animation-delay-2000 -z-10" />
      <div className="absolute bottom-0 left-1/2 w-96 h-96 bg-pink-100 rounded-full mix-blend-multiply filter blur-3xl opacity-40 animate-blob animation-delay-4000 -z-10" />

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
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium mb-8 animate-fade-in">
              <span className="w-2 h-2 rounded-full bg-indigo-500 animate-pulse" />
              {badgeText}
            </div>
          )}

          {/* Title */}
          <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold tracking-tight text-gray-900 mb-6 animate-fade-in leading-[1.1]">
            {titleParts.map((part, i) =>
              part.highlight ? (
                <span
                  key={i}
                  className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600"
                >
                  {part.text}
                </span>
              ) : (
                <span key={i}>{part.text}</span>
              )
            )}
          </h1>

          {/* Subtitle */}
          {subtitle && (
            <p className="text-xl text-gray-600 mb-10 max-w-2xl mx-auto animate-fade-in leading-relaxed">
              {subtitle}
            </p>
          )}

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12 animate-fade-in">
            {ctaPrimaryText && (
              <Link
                href={ctaPrimaryLink}
                className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-xl bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-lg hover:from-indigo-700 hover:to-purple-700 transition-all shadow-lg shadow-indigo-200 hover:shadow-xl hover:shadow-indigo-300 hover:-translate-y-0.5"
              >
                {ctaPrimaryText}
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              </Link>
            )}
            {ctaSecondaryText && (
              <Link
                href={ctaSecondaryLink}
                className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-xl border-2 border-gray-200 text-gray-700 font-semibold text-lg hover:border-indigo-300 hover:text-indigo-700 hover:bg-indigo-50 transition-all"
              >
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {ctaSecondaryText}
              </Link>
            )}
          </div>

          {/* Social proof */}
          {socialProofText && (
            <p className="text-sm text-gray-500 animate-fade-in flex items-center justify-center gap-2">
              <span className="flex -space-x-2">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="w-7 h-7 rounded-full bg-gradient-to-br from-indigo-400 to-purple-400 border-2 border-white flex items-center justify-center text-white text-xs font-bold"
                  >
                    {String.fromCharCode(64 + i)}
                  </div>
                ))}
              </span>
              {socialProofText}
            </p>
          )}
        </div>

        {/* Hero image */}
        {heroImage && (
          <div className="mt-16 relative mx-auto max-w-5xl animate-fade-in">
            <div className="rounded-2xl overflow-hidden shadow-2xl border border-gray-200/50 ring-1 ring-gray-900/5">
              {/* Browser chrome */}
              <div className="bg-gray-100 px-4 py-3 flex items-center gap-2 border-b border-gray-200">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-red-400" />
                  <div className="w-3 h-3 rounded-full bg-yellow-400" />
                  <div className="w-3 h-3 rounded-full bg-green-400" />
                </div>
                <div className="flex-1 mx-4">
                  <div className="bg-white rounded-md px-3 py-1 text-xs text-gray-400 text-center">
                    app.growthOS.io/dashboard
                  </div>
                </div>
              </div>
              <Image
                src={heroImage}
                alt={heroImageAlt}
                width={1200}
                height={675}
                className="w-full h-auto"
                priority
              />
            </div>
            {/* Floating cards */}
            <div className="absolute -left-8 top-1/4 bg-white rounded-xl shadow-lg p-3 border border-gray-100 hidden lg:flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
                <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
              <div>
                <p className="text-xs text-gray-500">Revenue</p>
                <p className="text-sm font-bold text-gray-900">+24.5%</p>
              </div>
            </div>
            <div className="absolute -right-8 bottom-1/4 bg-white rounded-xl shadow-lg p-3 border border-gray-100 hidden lg:flex items-center gap-3">
              <div className="w-10 h-10 rounded-lg bg-indigo-100 flex items-center justify-center">
                <svg className="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <div>
                <p className="text-xs text-gray-500">Active Users</p>
                <p className="text-sm font-bold text-gray-900">10,247</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
