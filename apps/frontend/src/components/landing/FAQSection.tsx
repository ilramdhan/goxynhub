"use client";

import { useState } from "react";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface FAQSectionProps {
  section: PageSection;
  contents: ContentMap;
}

const defaultFAQs = [
  {
    question: "How do I get started?",
    answer: "Simply sign up for a free account, choose your plan, and you can have your landing page live in minutes. No technical knowledge required.",
  },
  {
    question: "Can I change my plan later?",
    answer: "Yes, you can upgrade or downgrade your plan at any time. Changes take effect immediately and billing is prorated.",
  },
  {
    question: "Is there a free trial?",
    answer: "Yes! All plans come with a 14-day free trial. No credit card required to start.",
  },
  {
    question: "What payment methods do you accept?",
    answer: "We accept all major credit cards (Visa, MasterCard, American Express), PayPal, and bank transfers for annual plans.",
  },
  {
    question: "Can I use my own domain?",
    answer: "Yes, you can connect your custom domain on Pro and Enterprise plans. We provide free SSL certificates for all custom domains.",
  },
  {
    question: "How secure is my data?",
    answer: "We take security seriously. All data is encrypted at rest and in transit. We are SOC 2 Type II certified and GDPR compliant.",
  },
];

export function FAQSection({ section, contents }: FAQSectionProps) {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const badgeText = getContentValue(contents, "badge_text", "FAQ");
  const title = getContentValue(contents, "title", "Frequently asked questions");
  const subtitle = getContentValue(contents, "subtitle");

  return (
    <section
      id={section.identifier || "faq"}
      className={`py-20 lg:py-32 ${section.css_class || ""}`}
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

        <div className="max-w-3xl mx-auto space-y-4">
          {defaultFAQs.map((faq, index) => (
            <div
              key={index}
              className="bg-white rounded-xl border border-gray-200 overflow-hidden"
            >
              <button
                className="w-full text-left px-6 py-5 flex items-center justify-between gap-4 hover:bg-gray-50 transition-colors"
                onClick={() => setOpenIndex(openIndex === index ? null : index)}
              >
                <span className="font-semibold text-gray-900">{faq.question}</span>
                <span className="text-gray-400 text-xl flex-shrink-0">
                  {openIndex === index ? "−" : "+"}
                </span>
              </button>
              {openIndex === index && (
                <div className="px-6 pb-5 text-gray-600 leading-relaxed">
                  {faq.answer}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
