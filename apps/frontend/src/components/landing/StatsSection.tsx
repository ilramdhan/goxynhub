import React from "react";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface StatsSectionProps {
  section: PageSection;
  contents: ContentMap;
}

export function StatsSection({ section, contents }: StatsSectionProps) {
  const stats = [
    {
      value: getContentValue(contents, "stat_1_value", "10,000+"),
      label: getContentValue(contents, "stat_1_label", "Active Teams"),
    },
    {
      value: getContentValue(contents, "stat_2_value", "99.9%"),
      label: getContentValue(contents, "stat_2_label", "Uptime SLA"),
    },
    {
      value: getContentValue(contents, "stat_3_value", "50+"),
      label: getContentValue(contents, "stat_3_label", "Countries"),
    },
    {
      value: getContentValue(contents, "stat_4_value", "24/7"),
      label: getContentValue(contents, "stat_4_label", "Support"),
    },
  ].filter((s) => s.value && s.label);

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  return (
    <section
      id={section.identifier || "stats"}
      className={`py-20 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative rounded-3xl overflow-hidden bg-gradient-to-br from-indigo-600 via-purple-600 to-indigo-800 p-12 lg:p-16">
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
          
          <div className="relative grid grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-12">
            {stats.map((stat, index) => (
              <div key={index} className="text-center">
                <div className="text-4xl sm:text-5xl font-bold text-white mb-2 tracking-tight">
                  {stat.value}
                </div>
                <div className="text-indigo-200 text-sm sm:text-base font-medium">
                  {stat.label}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
