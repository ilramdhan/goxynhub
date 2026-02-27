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
      label: getContentValue(contents, "stat_1_label", "Active Users"),
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
  ];

  return (
    <section
      id={section.identifier || "stats"}
      className={`py-16 bg-primary ${section.css_class || ""}`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
          {stats.map((stat, index) => (
            <div key={index} className="text-center">
              <div className="text-4xl lg:text-5xl font-bold text-white mb-2">
                {stat.value}
              </div>
              <div className="text-primary-foreground/80 text-sm lg:text-base">
                {stat.label}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
