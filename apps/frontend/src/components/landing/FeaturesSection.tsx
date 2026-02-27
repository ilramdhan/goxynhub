import type { PageSection, ContentMap, Feature } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface FeaturesSectionProps {
  section: PageSection;
  contents: ContentMap;
  features?: Feature[];
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
  heart: "❤️",
  bolt: "⚡",
  cloud: "☁️",
  database: "🗄️",
  settings: "⚙️",
  chart: "📈",
};

const defaultFeatures: Feature[] = [
  {
    id: "1",
    site_id: "",
    section_id: null,
    title: "Lightning Fast",
    description: "Built for performance. Load times under 100ms with our optimized infrastructure.",
    icon: "zap",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 1,
    created_at: "",
    updated_at: "",
  },
  {
    id: "2",
    site_id: "",
    section_id: null,
    title: "Secure by Default",
    description: "Enterprise-grade security with end-to-end encryption and SOC 2 compliance.",
    icon: "shield",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 2,
    created_at: "",
    updated_at: "",
  },
  {
    id: "3",
    site_id: "",
    section_id: null,
    title: "Easy Integration",
    description: "Connect with 100+ tools you already use. REST API and webhooks included.",
    icon: "plug",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 3,
    created_at: "",
    updated_at: "",
  },
  {
    id: "4",
    site_id: "",
    section_id: null,
    title: "Real-time Analytics",
    description: "Track every interaction with detailed analytics and custom dashboards.",
    icon: "bar-chart",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 4,
    created_at: "",
    updated_at: "",
  },
  {
    id: "5",
    site_id: "",
    section_id: null,
    title: "Team Collaboration",
    description: "Work together seamlessly with roles, permissions, and real-time updates.",
    icon: "users",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 5,
    created_at: "",
    updated_at: "",
  },
  {
    id: "6",
    site_id: "",
    section_id: null,
    title: "24/7 Support",
    description: "Our expert team is always available to help you succeed.",
    icon: "headphones",
    icon_color: null,
    image_url: null,
    image_alt: null,
    link_url: null,
    link_text: null,
    is_active: true,
    sort_order: 6,
    created_at: "",
    updated_at: "",
  },
];

export function FeaturesSection({ section, contents, features }: FeaturesSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Everything you need to succeed");
  const subtitle = getContentValue(contents, "subtitle");

  // Use real data from API if available, otherwise use defaults
  const displayFeatures = features && features.length > 0 ? features : defaultFeatures;

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
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">{title}</h2>
          {subtitle && <p className="text-xl text-gray-600">{subtitle}</p>}
        </div>

        {/* Features grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {displayFeatures.map((feature) => (
            <div
              key={feature.id}
              className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
            >
              {feature.image_url ? (
                <img
                  src={feature.image_url}
                  alt={feature.image_alt || feature.title}
                  className="w-12 h-12 object-cover rounded-lg mb-4"
                />
              ) : (
                <div
                  className="text-4xl mb-4"
                  style={{ color: feature.icon_color || undefined }}
                >
                  {(feature.icon && iconMap[feature.icon]) || "✨"}
                </div>
              )}
              <h3 className="text-xl font-semibold text-gray-900 mb-2">{feature.title}</h3>
              {feature.description && (
                <p className="text-gray-600">{feature.description}</p>
              )}
              {feature.link_url && feature.link_text && (
                <a
                  href={feature.link_url}
                  className="inline-flex items-center mt-3 text-primary text-sm font-medium hover:underline"
                >
                  {feature.link_text} →
                </a>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
