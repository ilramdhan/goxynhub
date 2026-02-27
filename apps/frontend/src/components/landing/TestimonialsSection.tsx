import type { PageSection, ContentMap, Testimonial } from "@/types/api.types";
import { getContentValue } from "@/types/api.types";

interface TestimonialsSectionProps {
  section: PageSection;
  contents: ContentMap;
  testimonials?: Testimonial[];
}

const defaultTestimonials: Testimonial[] = [
  {
    id: "1",
    site_id: "",
    section_id: null,
    author_name: "Sarah Johnson",
    author_title: "CEO",
    author_company: "TechCorp",
    author_avatar: null,
    content: "This platform transformed how we manage our digital presence. The CMS is incredibly intuitive and the performance is outstanding.",
    rating: 5,
    source: null,
    source_url: null,
    is_featured: true,
    is_active: true,
    sort_order: 1,
    created_at: "",
    updated_at: "",
  },
  {
    id: "2",
    site_id: "",
    section_id: null,
    author_name: "Michael Chen",
    author_title: "CTO",
    author_company: "StartupXYZ",
    author_avatar: null,
    content: "We migrated from a legacy system and the difference is night and day. Setup took less than an hour and the team was up and running immediately.",
    rating: 5,
    source: null,
    source_url: null,
    is_featured: true,
    is_active: true,
    sort_order: 2,
    created_at: "",
    updated_at: "",
  },
  {
    id: "3",
    site_id: "",
    section_id: null,
    author_name: "Emily Rodriguez",
    author_title: "Marketing Director",
    author_company: "GrowthCo",
    author_avatar: null,
    content: "The ability to update landing page content without touching code has been a game-changer for our marketing team.",
    rating: 5,
    source: null,
    source_url: null,
    is_featured: false,
    is_active: true,
    sort_order: 3,
    created_at: "",
    updated_at: "",
  },
];

export function TestimonialsSection({ section, contents, testimonials }: TestimonialsSectionProps) {
  const badgeText = getContentValue(contents, "badge_text", "Testimonials");
  const title = getContentValue(contents, "title", "Loved by thousands of teams");
  const subtitle = getContentValue(contents, "subtitle");

  const displayTestimonials =
    testimonials && testimonials.length > 0 ? testimonials : defaultTestimonials;

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
          {displayTestimonials.map((testimonial) => (
            <div
              key={testimonial.id}
              className="bg-white rounded-2xl p-8 shadow-sm border border-gray-100"
            >
              {/* Stars */}
              {testimonial.rating && (
                <div className="flex gap-1 mb-4">
                  {Array.from({ length: testimonial.rating }).map((_, i) => (
                    <span key={i} className="text-yellow-400 text-lg">★</span>
                  ))}
                </div>
              )}

              {/* Content */}
              <p className="text-gray-700 mb-6 leading-relaxed">
                &ldquo;{testimonial.content}&rdquo;
              </p>

              {/* Author */}
              <div className="flex items-center gap-3">
                {testimonial.author_avatar ? (
                  <img
                    src={testimonial.author_avatar}
                    alt={testimonial.author_name}
                    className="w-10 h-10 rounded-full object-cover"
                  />
                ) : (
                  <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-semibold">
                    {testimonial.author_name.charAt(0)}
                  </div>
                )}
                <div>
                  <div className="font-semibold text-gray-900">{testimonial.author_name}</div>
                  <div className="text-sm text-gray-500">
                    {[testimonial.author_title, testimonial.author_company]
                      .filter(Boolean)
                      .join(" at ")}
                  </div>
                </div>
              </div>

              {/* Source */}
              {testimonial.source && testimonial.source_url && (
                <a
                  href={testimonial.source_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="mt-3 text-xs text-gray-400 hover:text-gray-600 block"
                >
                  via {testimonial.source}
                </a>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
