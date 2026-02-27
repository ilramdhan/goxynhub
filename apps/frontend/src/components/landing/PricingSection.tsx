import Link from "next/link";
import type { PageSection, ContentMap, PricingPlan } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface PricingSectionProps {
  section: PageSection;
  contents: ContentMap;
  plans?: PricingPlan[];
}

const defaultPlans: PricingPlan[] = [
  {
    id: "1",
    site_id: "",
    section_id: null,
    name: "Starter",
    description: "Perfect for individuals and small projects",
    price_monthly: 9,
    price_yearly: 90,
    currency: "USD",
    price_label: null,
    is_popular: false,
    is_custom: false,
    badge_text: null,
    cta_text: "Get Started",
    cta_link: "/signup?plan=starter",
    features: ["1 website", "5 pages", "10GB storage", "Basic analytics", "Email support"],
    features_excluded: [],
    is_active: true,
    sort_order: 1,
    created_at: "",
    updated_at: "",
  },
  {
    id: "2",
    site_id: "",
    section_id: null,
    name: "Pro",
    description: "For growing businesses and teams",
    price_monthly: 29,
    price_yearly: 290,
    currency: "USD",
    price_label: null,
    is_popular: true,
    is_custom: false,
    badge_text: "Most Popular",
    cta_text: "Start Free Trial",
    cta_link: "/signup?plan=pro",
    features: ["5 websites", "Unlimited pages", "50GB storage", "Advanced analytics", "Priority support", "Custom domain", "API access"],
    features_excluded: [],
    is_active: true,
    sort_order: 2,
    created_at: "",
    updated_at: "",
  },
  {
    id: "3",
    site_id: "",
    section_id: null,
    name: "Enterprise",
    description: "For large organizations with custom needs",
    price_monthly: null,
    price_yearly: null,
    currency: "USD",
    price_label: null,
    is_popular: false,
    is_custom: true,
    badge_text: null,
    cta_text: "Contact Sales",
    cta_link: "/contact",
    features: ["Unlimited websites", "Unlimited pages", "500GB storage", "Custom analytics", "24/7 dedicated support", "SLA guarantee"],
    features_excluded: [],
    is_active: true,
    sort_order: 3,
    created_at: "",
    updated_at: "",
  },
];

export function PricingSection({ section, contents, plans }: PricingSectionProps) {
  const badgeText = getContentValue(contents, "badge_text", "Pricing");
  const title = getContentValue(contents, "title", "Simple, transparent pricing");
  const subtitle = getContentValue(contents, "subtitle", "Choose the plan that works best for you.");

  const displayPlans = plans && plans.length > 0 ? plans : defaultPlans;

  return (
    <section
      id={section.identifier || "pricing"}
      className={`py-20 lg:py-32 bg-gray-50 ${section.css_class || ""}`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          {badgeText && (
            <div className="inline-flex items-center px-4 py-2 rounded-full bg-primary/10 text-primary text-sm font-medium mb-4">
              {badgeText}
            </div>
          )}
          <h2 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">{title}</h2>
          {subtitle && <p className="text-xl text-gray-600">{subtitle}</p>}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {displayPlans.map((plan) => {
            const planFeatures = Array.isArray(plan.features)
              ? plan.features.map((f) => String(f))
              : [];

            return (
              <div
                key={plan.id}
                className={`relative bg-white rounded-2xl p-8 shadow-sm border ${
                  plan.is_popular
                    ? "border-primary shadow-lg scale-105"
                    : "border-gray-100"
                }`}
              >
                {(plan.is_popular || plan.badge_text) && (
                  <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                    <span className="bg-primary text-white text-sm font-semibold px-4 py-1 rounded-full">
                      {plan.badge_text || "Most Popular"}
                    </span>
                  </div>
                )}

                <div className="mb-6">
                  <h3 className="text-xl font-bold text-gray-900 mb-2">{plan.name}</h3>
                  {plan.description && (
                    <p className="text-gray-600 text-sm">{plan.description}</p>
                  )}
                </div>

                <div className="mb-6">
                  {plan.is_custom ? (
                    <div className="text-3xl font-bold text-gray-900">Custom</div>
                  ) : (
                    <div>
                      <span className="text-4xl font-bold text-gray-900">
                        ${plan.price_monthly}
                      </span>
                      <span className="text-gray-500">
                        /{plan.price_label || "month"}
                      </span>
                      {plan.price_yearly && (
                        <p className="text-sm text-green-600 mt-1">
                          ${plan.price_yearly}/year (save {Math.round((1 - plan.price_yearly / (plan.price_monthly! * 12)) * 100)}%)
                        </p>
                      )}
                    </div>
                  )}
                </div>

                <ul className="space-y-3 mb-8">
                  {planFeatures.map((feature, i) => (
                    <li key={i} className="flex items-center gap-2 text-sm text-gray-700">
                      <span className="text-green-500 font-bold flex-shrink-0">✓</span>
                      {feature}
                    </li>
                  ))}
                </ul>

                <Link
                  href={plan.cta_link || "#"}
                  className={`block w-full text-center py-3 px-6 rounded-lg font-semibold transition-colors ${
                    plan.is_popular
                      ? "bg-primary text-white hover:bg-primary/90"
                      : "border-2 border-gray-300 text-gray-700 hover:border-primary hover:text-primary"
                  }`}
                >
                  {plan.cta_text}
                </Link>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
