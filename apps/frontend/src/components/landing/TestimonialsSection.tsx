import Image from "next/image";
import type { PageSection, ContentMap } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface TestimonialsSectionProps {
  section: PageSection;
  contents: ContentMap;
}

const defaultTestimonials = [
  {
    author_name: "Sarah Johnson",
    author_title: "CEO at TechCorp",
    author_avatar: null,
    content: "This platform transformed how we manage our digital presence. The CMS is incredibly intuitive.",
    rating: 5,
  },
  {
    author_name: "Michael Chen",
    author_title: "CTO at StartupXYZ",
    author_avatar: null,
    content: "We migrated from a legacy system and the difference is night and day. Setup took less than an hour.",
    rating: 5,
  },
  {
    author_name: "Emily Rodriguez",
    author_title: "Marketing Director at GrowthCo",
    author_avatar: null,
    content: "The ability to update landing page content without touching code has been a game-changer.",
    rating: 5,
  },
];

export function TestimonialsSection({ section, contents }: TestimonialsSectionProps) {
  const badgeText = getContentValue(contents, "badge_text", "Testimonials");
  const title = getContentValue(contents, "title", "Loved by thousands of teams");
  const subtitle = getContentValue(contents, "subtitle");

  return (
    <section
      id={section.identifier || "testimonials"}
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

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {defaultTestimonials.map((testimonial, index) => (
            <div
              key={index}
              className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100"
            >
              {/* Stars */}
              <div className="flex gap-1 mb-4">
                {Array.from({ length: testimonial.rating }).map((_, i) => (
                  <span key={i} className="text-yellow-400 text-lg">★</span>
                ))}
              </div>

              {/* Content */}
              <p className="text-gray-700 mb-6 leading-relaxed">
                &ldquo;{testimonial.content}&rdquo;
              </p>

              {/* Author */}
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold">
                  {testimonial.author_name.charAt(0)}
                </div>
                <div>
                  <div className="font-semibold text-gray-900">{testimonial.author_name}</div>
                  <div className="text-sm text-gray-500">{testimonial.author_title}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
