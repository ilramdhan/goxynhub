import React from "react";
import Link from "next/link";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface CTASectionProps {
  section: PageSection;
  contents: ContentMap;
}

export function CTASection({ section, contents }: CTASectionProps) {
  const title = getContentValue(contents, "title", "Ready to scale your business?");
  const subtitle = getContentValue(
    contents,
    "subtitle",
    "Join 10,000+ teams already growing with GrowthOS. Start your free trial today — no credit card required."
  );
  const ctaPrimaryText = getContentValue(contents, "cta_primary_text", "Start Free Trial");
  const ctaPrimaryLink = getContentValue(contents, "cta_primary_link", "#pricing");
  const ctaSecondaryText = getContentValue(contents, "cta_secondary_text");
  const ctaSecondaryLink = getContentValue(contents, "cta_secondary_link", "/contact");

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  return (
    <section
      id={section.identifier || "cta"}
      className={`py-24 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative rounded-3xl overflow-hidden bg-gradient-to-br from-indigo-600 via-purple-600 to-indigo-800 px-8 py-16 lg:px-16 text-center">
          {/* Background pattern */}
          <div className="absolute inset-0 opacity-10">
            <div className="absolute inset-0" style={{
              backgroundImage: `radial-gradient(circle at 1px 1px, white 1px, transparent 0)`,
              backgroundSize: "40px 40px",
            }} />
          </div>
          
          {/* Decorative circles */}
          <div className="absolute -top-20 -right-20 w-80 h-80 rounded-full bg-white/5" />
          <div className="absolute -bottom-20 -left-20 w-80 h-80 rounded-full bg-white/5" />
          
          <div className="relative max-w-3xl mx-auto">
            <h2 className="text-4xl sm:text-5xl font-bold text-white mb-4 tracking-tight">
              {title}
            </h2>
            {subtitle && (
              <p className="text-xl text-indigo-200 mb-10 leading-relaxed">
                {subtitle}
              </p>
            )}
            
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              {ctaPrimaryText && (
                <Link
                  href={ctaPrimaryLink}
                  className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-xl bg-white text-indigo-700 font-semibold text-lg hover:bg-indigo-50 transition-all shadow-lg hover:shadow-xl hover:-translate-y-0.5"
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
                  className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-xl border-2 border-white/30 text-white font-semibold text-lg hover:border-white/60 hover:bg-white/10 transition-all"
                >
                  {ctaSecondaryText}
                </Link>
              )}
            </div>
            
            {/* Trust indicators */}
            <div className="mt-10 flex flex-wrap items-center justify-center gap-6 text-indigo-200 text-sm">
              <span className="flex items-center gap-1.5">
                <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                No credit card required
              </span>
              <span className="flex items-center gap-1.5">
                <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                14-day free trial
              </span>
              <span className="flex items-center gap-1.5">
                <svg className="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                Cancel anytime
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
