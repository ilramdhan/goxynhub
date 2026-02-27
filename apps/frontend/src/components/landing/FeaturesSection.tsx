import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface FeaturesSectionProps {
  section: PageSection;
  contents: ContentMap;
}

// Icon mapping for common feature icons
const iconMap: Record<string, string> = {
  zap: "⚡",
  shield: "🛡️",
  plug: "🔌",
  "bar-chart": "📊",
  users: "👥",
  headphones: "🎧",
  star: "⭐",
  check: "✅",
  lock: "🔒",
  globe: "🌐",
  code: "💻",
  rocket: "🚀",
};

export function FeaturesSection({ section, contents }: FeaturesSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Everything you need to succeed");
  const subtitle = getContentValue(contents, "subtitle");

  // Features are loaded from the features table via the section
  // For now, we'll use placeholder data structure
  const features = [
    {
      icon: "zap",
      title: "Lightning Fast",
      description: "Built for performance. Load times under 100ms.",
    },
    {
      icon: "shield",
      title: "Secure by Default",
      description: "Enterprise-grade security with end-to-end encryption.",
    },
    {
      icon: "plug",
      title: "Easy Integration",
      description: "Connect with 100+ tools you already use.",
    },
    {
      icon: "bar-chart",
      title: "Real-time Analytics",
      description: "Track every interaction with detailed analytics.",
    },
    {
      icon: "users",
      title: "Team Collaboration",
      description: "Work together seamlessly with roles and permissions.",
    },
    {
      icon: "headphones",
      title: "24/7 Support",
      description: "Our expert team is always available to help.",
    },
  ];

  return (
    <section
      id={section.identifier || "features"}
      className={`py-20 lg:py-32 bg-gray-50 ${section.css_class || ""}`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          {badgeText && (
            <div className="inline-flex items-center px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-4">
              {badgeText}
            </div>
          )}
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          {subtitle && (
            <p className="text-xl text-gray-600">{subtitle}</p>
          )}
        </div>

        {/* Features grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <div
              key={index}
              className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
            >
              <div className="text-4xl mb-4">
                {iconMap[feature.icon] || "✨"}
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                {feature.title}
              </h3>
              <p className="text-gray-600">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
