"use client";

import React, { useState } from "react";
import Link from "next/link";
import type { PageSection, ContentMap, PricingPlan } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface PricingSectionProps {
  section: PageSection;
  contents: ContentMap;
  plans?: PricingPlan[];
}

export function PricingSection({ section, contents, plans = [] }: PricingSectionProps) {
  const [isYearly, setIsYearly] = useState(false);

  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Simple, transparent pricing");
  const subtitle = getContentValue(contents, "subtitle");

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  const activePlans = plans.filter((p) => p.is_active);

  const defaultPlans = [
    {
      id: "starter",
      name: "Starter",
      description: "Perfect for small teams",
      price_monthly: 29,
      price_yearly: 290,
      currency: "USD",
      is_popular: false,
      is_custom: false,
      cta_text: "Start Free Trial",
      cta_link: "#",
      features: ["Up to 5 team members", "10,000 events/month", "Basic analytics", "Email support", "API access"],
      features_excluded: ["Advanced analytics", "Priority support", "SSO"],
      badge_text: null,
      is_active: true,
      sort_order: 1,
    },
    {
      id: "pro",
      name: "Pro",
      description: "For growing teams",
      price_monthly: 79,
      price_yearly: 790,
      currency: "USD",
      is_popular: true,
      is_custom: false,
      cta_text: "Start Free Trial",
      cta_link: "#",
      features: ["Up to 25 team members", "500K events/month", "Advanced analytics", "Priority support", "Full API access", "Unlimited integrations", "Custom dashboards"],
      features_excluded: ["SSO", "Dedicated support"],
      badge_text: "Most Popular",
      is_active: true,
      sort_order: 2,
    },
    {
      id: "enterprise",
      name: "Enterprise",
      description: "For large organizations",
      price_monthly: null,
      price_yearly: null,
      currency: "USD",
      is_popular: false,
      is_custom: true,
      cta_text: "Contact Sales",
      cta_link: "/contact",
      features: ["Unlimited members", "Unlimited events", "Custom analytics", "24/7 dedicated support", "SSO & SAML", "Custom SLA", "On-premise option"],
      features_excluded: [],
      badge_text: null,
      is_active: true,
      sort_order: 3,
    },
  ];

  const displayPlans = activePlans.length > 0 ? activePlans : defaultPlans;

  const getPrice = (plan: typeof defaultPlans[0] | PricingPlan) => {
    if (plan.is_custom) return null;
    return isYearly ? plan.price_yearly : plan.price_monthly;
  };

  const getFeatures = (plan: typeof defaultPlans[0] | PricingPlan): string[] => {
    if (Array.isArray(plan.features)) {
      return plan.features.map((f) => String(f));
    }
    return [];
  };

  return (
    <section
      id={section.identifier || "pricing"}
      className={`py-24 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-12">
          {badgeText && (
            <div className="inline-flex items-center px-4 py-1.5 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-700 text-sm font-medium mb-4">
              {badgeText}
            </div>
          )}
          <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4 tracking-tight">
            {title}
          </h2>
          {subtitle && (
            <p className="text-xl text-gray-600 leading-relaxed mb-8">{subtitle}</p>
          )}

          {/* Billing toggle */}
          <div className="inline-flex items-center gap-3 bg-gray-100 rounded-xl p-1">
            <button
              onClick={() => setIsYearly(false)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                !isYearly
                  ? "bg-white text-gray-900 shadow-sm"
                  : "text-gray-500 hover:text-gray-700"
              }`}
            >
              Monthly
            </button>
            <button
              onClick={() => setIsYearly(true)}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-all flex items-center gap-2 ${
                isYearly
                  ? "bg-white text-gray-900 shadow-sm"
                  : "text-gray-500 hover:text-gray-700"
              }`}
            >
              Yearly
              <span className="text-xs bg-green-100 text-green-700 px-1.5 py-0.5 rounded-full font-semibold">
                Save 17%
              </span>
            </button>
          </div>
        </div>

        {/* Pricing Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {displayPlans.map((plan) => {
            const price = getPrice(plan);
            const features = getFeatures(plan);
            const isPopular = plan.is_popular;

            return (
              <div
                key={"id" in plan ? plan.id : plan.name}
                className={`relative rounded-2xl p-8 flex flex-col ${
                  isPopular
                    ? "bg-gradient-to-b from-indigo-600 to-purple-700 text-white shadow-2xl shadow-indigo-200 scale-105"
                    : "bg-white border border-gray-200 shadow-sm"
                }`}
              >
                {/* Popular badge */}
                {plan.badge_text && (
                  <div className={`absolute -top-3 left-1/2 -translate-x-1/2 px-4 py-1 rounded-full text-xs font-bold ${
                    isPopular
                      ? "bg-amber-400 text-amber-900"
                      : "bg-indigo-100 text-indigo-700"
                  }`}>
                    {plan.badge_text}
                  </div>
                )}

                {/* Plan name */}
                <div className="mb-6">
                  <h3 className={`text-xl font-bold mb-1 ${isPopular ? "text-white" : "text-gray-900"}`}>
                    {plan.name}
                  </h3>
                  {plan.description && (
                    <p className={`text-sm ${isPopular ? "text-indigo-200" : "text-gray-500"}`}>
                      {plan.description}
                    </p>
                  )}
                </div>

                {/* Price */}
                <div className="mb-6">
                  {plan.is_custom ? (
                    <div>
                      <span className={`text-4xl font-bold ${isPopular ? "text-white" : "text-gray-900"}`}>
                        Custom
                      </span>
                    </div>
                  ) : (
                    <div className="flex items-end gap-1">
                      <span className={`text-sm font-medium ${isPopular ? "text-indigo-200" : "text-gray-500"}`}>
                        $
                      </span>
                      <span className={`text-5xl font-bold tracking-tight ${isPopular ? "text-white" : "text-gray-900"}`}>
                        {price}
                      </span>
                      <span className={`text-sm mb-1 ${isPopular ? "text-indigo-200" : "text-gray-500"}`}>
                        /{isYearly ? "yr" : "mo"}
                      </span>
                    </div>
                  )}
                </div>

                {/* CTA */}
                <Link
                  href={plan.cta_link || "#"}
                  className={`block text-center py-3 px-6 rounded-xl font-semibold text-sm transition-all mb-8 ${
                    isPopular
                      ? "bg-white text-indigo-700 hover:bg-indigo-50"
                      : "bg-indigo-600 text-white hover:bg-indigo-700"
                  }`}
                >
                  {plan.cta_text}
                </Link>

                {/* Features */}
                <ul className="space-y-3 flex-1">
                  {features.map((feature, i) => (
                    <li key={i} className="flex items-start gap-2.5">
                      <svg
                        className={`w-5 h-5 flex-shrink-0 mt-0.5 ${isPopular ? "text-indigo-200" : "text-indigo-500"}`}
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                      <span className={`text-sm ${isPopular ? "text-indigo-100" : "text-gray-600"}`}>
                        {feature}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            );
          })}
        </div>

        {/* Money back guarantee */}
        <p className="text-center text-sm text-gray-500 mt-8 flex items-center justify-center gap-2">
          <svg className="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
          30-day money-back guarantee · No credit card required · Cancel anytime
        </p>
      </div>
    </section>
  );
}
