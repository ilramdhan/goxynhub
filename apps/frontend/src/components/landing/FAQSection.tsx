"use client";

import React, { useState } from "react";
import type { PageSection, ContentMap, FAQ } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface FAQSectionProps {
  section: PageSection;
  contents: ContentMap;
  faqs?: FAQ[];
}

function FAQItem({ faq }: { faq: FAQ | { question: string; answer: string; category?: string | null } }) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border border-gray-200 rounded-xl overflow-hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex items-center justify-between px-6 py-4 text-left hover:bg-gray-50 transition-colors"
      >
        <span className="font-medium text-gray-900 pr-4">{faq.question}</span>
        <div className={`flex-shrink-0 w-6 h-6 rounded-full border-2 border-gray-300 flex items-center justify-center transition-transform ${isOpen ? "rotate-45 border-indigo-500" : ""}`}>
          <svg className={`w-3 h-3 ${isOpen ? "text-indigo-500" : "text-gray-400"}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 4v16m8-8H4" />
          </svg>
        </div>
      </button>
      {isOpen && (
        <div className="px-6 pb-4 text-gray-600 text-sm leading-relaxed border-t border-gray-100 pt-4">
          {faq.answer}
        </div>
      )}
    </div>
  );
}

export function FAQSection({ section, contents, faqs = [] }: FAQSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Frequently asked questions");
  const subtitle = getContentValue(contents, "subtitle");

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  const activeFAQs = faqs.filter((f) => f.is_active);

  const defaultFAQs = [
    { question: "What is GrowthOS?", answer: "GrowthOS is an all-in-one business growth platform that combines analytics, automation, and team collaboration tools to help you scale faster.", category: "General" },
    { question: "How does the free trial work?", answer: "You can start a 14-day free trial with full access to all Pro features — no credit card required. After the trial, choose to upgrade or continue with our free plan.", category: "General" },
    { question: "Can I change my plan at any time?", answer: "Yes! You can upgrade or downgrade your plan at any time. When you upgrade, you'll be charged the prorated difference.", category: "Billing" },
    { question: "What payment methods do you accept?", answer: "We accept all major credit cards (Visa, Mastercard, American Express), PayPal, and bank transfers for annual plans.", category: "Billing" },
    { question: "Is my data secure?", answer: "Absolutely. GrowthOS is SOC 2 Type II certified and GDPR compliant. We use AES-256 encryption for data at rest and TLS 1.3 for data in transit.", category: "Security" },
    { question: "What kind of support do you offer?", answer: "Starter plan customers get email support. Pro plan customers get priority email and live chat support. Enterprise customers get 24/7 dedicated support.", category: "Support" },
  ];

  const displayFAQs = activeFAQs.length > 0 ? activeFAQs : defaultFAQs;

  // Group FAQs by category
  const categories = [...new Set(displayFAQs.map((f) => f.category).filter(Boolean))];
  const [activeCategory, setActiveCategory] = useState<string | null>(null);

  const filteredFAQs = activeCategory
    ? displayFAQs.filter((f) => f.category === activeCategory)
    : displayFAQs;

  return (
    <section
      id={section.identifier || "faq"}
      className={`py-24 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="max-w-3xl mx-auto">
          {/* Header */}
          <div className="text-center mb-12">
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

          {/* Category filter */}
          {categories.length > 1 && (
            <div className="flex flex-wrap gap-2 justify-center mb-8">
              <button
                onClick={() => setActiveCategory(null)}
                className={`px-4 py-1.5 rounded-full text-sm font-medium transition-all ${
                  !activeCategory
                    ? "bg-indigo-600 text-white"
                    : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                }`}
              >
                All
              </button>
              {categories.map((cat) => (
                <button
                  key={cat}
                  onClick={() => setActiveCategory(cat || null)}
                  className={`px-4 py-1.5 rounded-full text-sm font-medium transition-all ${
                    activeCategory === cat
                      ? "bg-indigo-600 text-white"
                      : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                  }`}
                >
                  {cat}
                </button>
              ))}
            </div>
          )}

          {/* FAQ Items */}
          <div className="space-y-3">
            {filteredFAQs.map((faq, index) => (
              <FAQItem key={"id" in faq ? faq.id : index} faq={faq} />
            ))}
          </div>

          {/* Contact CTA */}
          <div className="mt-12 text-center p-8 bg-gray-50 rounded-2xl border border-gray-200">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              Still have questions?
            </h3>
            <p className="text-gray-600 text-sm mb-4">
              Can&apos;t find the answer you&apos;re looking for? Our support team is here to help.
            </p>
            <a
              href="/contact"
              className="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl bg-indigo-600 text-white font-medium text-sm hover:bg-indigo-700 transition-colors"
            >
              Contact Support
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
              </svg>
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
