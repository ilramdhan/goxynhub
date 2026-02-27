import Link from "next/link";
import type { NavigationMenu, SiteSettingsMap } from "@/types/api.types";

interface SiteFooterProps {
  siteId: string;
  footerNavigation?: NavigationMenu | null;
  settings?: SiteSettingsMap;
}

const defaultFooterLinks = [
  { label: "Privacy Policy", href: "/privacy" },
  { label: "Terms of Service", href: "/terms" },
  { label: "Cookie Policy", href: "/cookies" },
];

const defaultSocialLinks = [
  { label: "Twitter", href: "#", icon: "𝕏" },
  { label: "LinkedIn", href: "#", icon: "in" },
  { label: "GitHub", href: "#", icon: "⌥" },
];

export function SiteFooter({ siteId, footerNavigation, settings }: SiteFooterProps) {
  const currentYear = new Date().getFullYear();

  // Use CMS data if available, otherwise use defaults
  const siteName = settings?.site_title || "LandingCMS";
  const siteDescription = settings?.site_description || "The most powerful platform to create, manage, and scale your digital presence.";
  const contactEmail = settings?.contact_email;
  const contactPhone = settings?.contact_phone;

  // Social links from settings
  const socialLinks = [
    settings?.social_twitter && { label: "Twitter", href: settings.social_twitter, icon: "𝕏" },
    settings?.social_linkedin && { label: "LinkedIn", href: settings.social_linkedin, icon: "in" },
    settings?.social_github && { label: "GitHub", href: settings.social_github, icon: "⌥" },
    settings?.social_facebook && { label: "Facebook", href: settings.social_facebook, icon: "f" },
    settings?.social_instagram && { label: "Instagram", href: settings.social_instagram, icon: "◎" },
    settings?.social_youtube && { label: "YouTube", href: settings.social_youtube, icon: "▶" },
  ].filter(Boolean) as { label: string; href: string; icon: string }[];

  const displaySocialLinks = socialLinks.length > 0 ? socialLinks : defaultSocialLinks;

  // Footer navigation from CMS
  const footerLinks =
    footerNavigation?.items && footerNavigation.items.length > 0
      ? footerNavigation.items
          .filter((item) => item.is_active)
          .map((item) => ({ label: item.label, href: item.url || "#" }))
      : defaultFooterLinks;

  return (
    <footer className="bg-gray-900 text-gray-400 py-12">
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          {/* Brand */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center gap-2 mb-4">
              <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">
                  {siteName.charAt(0)}
                </span>
              </div>
              <span className="font-bold text-white text-lg">{siteName}</span>
            </div>
            <p className="text-sm leading-relaxed max-w-xs mb-4">{siteDescription}</p>
            {(contactEmail || contactPhone) && (
              <div className="space-y-1 text-sm">
                {contactEmail && (
                  <a href={`mailto:${contactEmail}`} className="block hover:text-white transition-colors">
                    ✉ {contactEmail}
                  </a>
                )}
                {contactPhone && (
                  <a href={`tel:${contactPhone}`} className="block hover:text-white transition-colors">
                    ☎ {contactPhone}
                  </a>
                )}
              </div>
            )}
          </div>

          {/* Product */}
          <div>
            <h4 className="text-white font-semibold mb-4">Product</h4>
            <ul className="space-y-2 text-sm">
              <li><Link href="/#features" className="hover:text-white transition-colors">Features</Link></li>
              <li><Link href="/#pricing" className="hover:text-white transition-colors">Pricing</Link></li>
              <li><Link href="/#faq" className="hover:text-white transition-colors">FAQ</Link></li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h4 className="text-white font-semibold mb-4">Legal</h4>
            <ul className="space-y-2 text-sm">
              {footerLinks.map((link) => (
                <li key={link.href}>
                  <Link href={link.href} className="hover:text-white transition-colors">
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 pt-8 flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-sm">
            © {currentYear} {siteName}. All rights reserved.
          </p>
          <div className="flex items-center gap-4">
            {displaySocialLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-white transition-colors text-sm"
                aria-label={link.label}
              >
                {link.icon}
              </a>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}
