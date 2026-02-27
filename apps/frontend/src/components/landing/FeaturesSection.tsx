import React from "react";
import type { PageSection, ContentMap, Feature } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface FeaturesSectionProps {
  section: PageSection;
  contents: ContentMap;
  features?: Feature[];
}

const ICON_MAP: Record<string, React.ReactNode> = {
  "bar-chart": (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
    </svg>
  ),
  zap: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  ),
  users: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
  ),
  shield: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
    </svg>
  ),
  plug: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  ),
  headphones: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
    </svg>
  ),
  globe: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  ),
  star: (
    <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
    </svg>
  ),
};

const GRADIENT_COLORS = [
  "from-indigo-500 to-purple-600",
  "from-purple-500 to-pink-600",
  "from-cyan-500 to-blue-600",
  "from-emerald-500 to-teal-600",
  "from-amber-500 to-orange-600",
  "from-rose-500 to-red-600",
  "from-blue-500 to-indigo-600",
  "from-violet-500 to-purple-600",
];

const BG_COLORS = [
  "bg-indigo-50",
  "bg-purple-50",
  "bg-cyan-50",
  "bg-emerald-50",
  "bg-amber-50",
  "bg-rose-50",
  "bg-blue-50",
  "bg-violet-50",
];

export function FeaturesSection({ section, contents, features = [] }: FeaturesSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Everything you need to grow");
  const subtitle = getContentValue(contents, "subtitle");

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  return (
    <section
      id={section.identifier || "features"}
      className={`py-24 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          {badgeText && (
            <div className="inline-flex items-center px-4 py-1.5 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium mb-4">
              {badgeText}
            </div>
          )}
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">
            {title}
          </h2>
          {subtitle && (
            <p className="text-xl text-gray-600 leading-relaxed">{subtitle}</p>
          )}
        </div>

        {/* Bento Grid */}
        {features.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 auto-rows-fr">
            {features.map((feature, index) => {
              const isLarge = index === 0 || index === 3;
              const gradient = GRADIENT_COLORS[index % GRADIENT_COLORS.length];
              const bgColor = BG_COLORS[index % BG_COLORS.length];
              const icon = feature.icon ? ICON_MAP[feature.icon] : null;

              return (
                <div
                  key={feature.id}
                  className={`group relative rounded-2xl p-6 border border-gray-100 bg-white hover:shadow-lg hover:-translate-y-1 transition-all duration-300 ${
                    isLarge ? "md:col-span-2" : ""
                  }`}
                >
                  {/* Icon */}
                  <div
                    className={`inline-flex w-12 h-12 rounded-xl ${bgColor} items-center justify-center mb-4`}
                    style={{ color: feature.icon_color || undefined }}
                  >
                    {icon || (
                      <div className={`w-6 h-6 rounded-lg bg-gradient-to-br ${gradient}`} />
                    )}
                  </div>

                  {/* Content */}
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    {feature.title}
                  </h3>
                  {feature.description && (
                    <p className="text-gray-600 text-sm leading-relaxed">
                      {feature.description}
                    </p>
                  )}

                  {/* Link */}
                  {feature.link_url && feature.link_text && (
                    <a
                      href={feature.link_url}
                      className="inline-flex items-center gap-1 mt-4 text-sm font-medium text-indigo-600 hover:text-indigo-700"
                    >
                      {feature.link_text}
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                      </svg>
                    </a>
                  )}

                  {/* Hover gradient overlay */}
                  <div className={`absolute inset-0 rounded-2xl bg-gradient-to-br ${gradient} opacity-0 group-hover:opacity-5 transition-opacity duration-300`} />
                </div>
              );
            })}
          </div>
        ) : (
          // Default features grid when no data
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
              { icon: "bar-chart", title: "Real-time Analytics", desc: "Get instant insights with live dashboards and custom reports.", color: "bg-indigo-50", textColor: "text-indigo-600" },
              { icon: "zap", title: "Smart Automation", desc: "Automate workflows to save time and reduce human error.", color: "bg-purple-50", textColor: "text-purple-600" },
              { icon: "users", title: "Team Collaboration", desc: "Work together with real-time collaboration tools.", color: "bg-cyan-50", textColor: "text-cyan-600" },
              { icon: "shield", title: "Enterprise Security", desc: "Bank-grade security with SOC 2 compliance and SSO.", color: "bg-emerald-50", textColor: "text-emerald-600" },
              { icon: "plug", title: "200+ Integrations", desc: "Connect with tools you already use via REST API.", color: "bg-amber-50", textColor: "text-amber-600" },
              { icon: "headphones", title: "24/7 Support", desc: "Expert support team available whenever you need help.", color: "bg-rose-50", textColor: "text-rose-600" },
            ].map((f, i) => (
              <div key={i} className="group rounded-2xl p-6 border border-gray-100 bg-white hover:shadow-lg hover:-translate-y-1 transition-all duration-300">
                <div className={`inline-flex w-12 h-12 rounded-xl ${f.color} ${f.textColor} items-center justify-center mb-4`}>
                  {ICON_MAP[f.icon]}
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">{f.title}</h3>
                <p className="text-gray-600 text-sm leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
