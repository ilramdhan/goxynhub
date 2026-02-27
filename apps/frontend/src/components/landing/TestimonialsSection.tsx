import React from "react";
import Image from "next/image";
import type { PageSection, ContentMap, Testimonial } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface TestimonialsSectionProps {
  section: PageSection;
  contents: ContentMap;
  testimonials?: Testimonial[];
}

function StarRating({ rating }: { rating: number }) {
  return (
    <div className="flex gap-0.5">
      {[1, 2, 3, 4, 5].map((star) => (
        <svg
          key={star}
          className={`w-4 h-4 ${star <= rating ? "text-amber-400" : "text-gray-200"}`}
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
        </svg>
      ))}
    </div>
  );
}

export function TestimonialsSection({ section, contents, testimonials = [] }: TestimonialsSectionProps) {
  const badgeText = getContentValue(contents, "badge_text");
  const title = getContentValue(contents, "title", "Loved by thousands of teams");
  const subtitle = getContentValue(contents, "subtitle");

  const sectionStyle: React.CSSProperties = {};
  if (section.bg_color) sectionStyle.backgroundColor = section.bg_color;

  const displayTestimonials = testimonials.filter((t) => t.is_active);

  return (
    <section
      id={section.identifier || "testimonials"}
      className={`py-24 ${section.css_class || ""}`}
      style={sectionStyle}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16">
          {badgeText && (
            <div className="inline-flex items-center px-4 py-1.5 rounded-full bg-amber-50 border border-amber-100 text-amber-700 text-sm font-medium mb-4">
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

        {/* Testimonials Grid */}
        {displayTestimonials.length > 0 ? (
          <div className="columns-1 md:columns-2 lg:columns-3 gap-6 space-y-6">
            {displayTestimonials.map((testimonial) => (
              <div
                key={testimonial.id}
                className="break-inside-avoid bg-white rounded-2xl p-6 border border-gray-100 shadow-sm hover:shadow-md transition-shadow"
              >
                {/* Rating */}
                {testimonial.rating && (
                  <StarRating rating={testimonial.rating} />
                )}

                {/* Content */}
                <blockquote className="mt-3 text-gray-700 text-sm leading-relaxed">
                  &ldquo;{testimonial.content}&rdquo;
                </blockquote>

                {/* Author */}
                <div className="mt-4 flex items-center gap-3">
                  {testimonial.author_avatar ? (
                    <Image
                      src={testimonial.author_avatar}
                      alt={testimonial.author_name}
                      width={40}
                      height={40}
                      className="w-10 h-10 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-400 to-purple-400 flex items-center justify-center text-white font-semibold text-sm">
                      {testimonial.author_name.charAt(0)}
                    </div>
                  )}
                  <div>
                    <div className="font-semibold text-gray-900 text-sm">
                      {testimonial.author_name}
                    </div>
                    <div className="text-gray-500 text-xs">
                      {[testimonial.author_title, testimonial.author_company]
                        .filter(Boolean)
                        .join(" at ")}
                    </div>
                  </div>
                  {testimonial.is_featured && (
                    <div className="ml-auto">
                      <span className="text-xs bg-amber-50 text-amber-700 px-2 py-0.5 rounded-full border border-amber-100">
                        Featured
                      </span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          // Default testimonials
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[
              {
                name: "Sarah Chen",
                title: "Head of Growth",
                company: "Stripe",
                content: "GrowthOS transformed how we track and optimize our growth metrics. The real-time analytics are incredibly powerful.",
                rating: 5,
              },
              {
                name: "Marcus Johnson",
                title: "CTO",
                company: "Notion",
                content: "We evaluated 10+ platforms before choosing GrowthOS. The API integrations are seamless and the support team is exceptional.",
                rating: 5,
              },
              {
                name: "Emily Rodriguez",
                title: "VP of Marketing",
                company: "Linear",
                content: "The AI-powered insights have been a game-changer. We've seen a 40% improvement in conversion rates since switching.",
                rating: 5,
              },
            ].map((t, i) => (
              <div key={i} className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                <StarRating rating={t.rating} />
                <blockquote className="mt-3 text-gray-700 text-sm leading-relaxed">
                  &ldquo;{t.content}&rdquo;
                </blockquote>
                <div className="mt-4 flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-400 to-purple-400 flex items-center justify-center text-white font-semibold text-sm">
                    {t.name.charAt(0)}
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 text-sm">{t.name}</div>
                    <div className="text-gray-500 text-xs">{t.title} at {t.company}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}
