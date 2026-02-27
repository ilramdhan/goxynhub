import { describe, it, expect } from "vitest";
import {
  cn,
  formatFileSize,
  generateSlug,
  truncate,
  formatDate,
  isExternalUrl,
  sanitizeHtml,
} from "../utils";

describe("cn", () => {
  it("merges class names", () => {
    expect(cn("foo", "bar")).toBe("foo bar");
  });

  it("handles conditional classes", () => {
    expect(cn("foo", false && "bar", "baz")).toBe("foo baz");
  });

  it("deduplicates tailwind classes", () => {
    expect(cn("p-4", "p-8")).toBe("p-8");
  });
});

describe("formatFileSize", () => {
  it("formats bytes", () => {
    expect(formatFileSize(500)).toBe("500 B");
  });

  it("formats kilobytes", () => {
    expect(formatFileSize(1024)).toBe("1.0 KB");
    expect(formatFileSize(1536)).toBe("1.5 KB");
  });

  it("formats megabytes", () => {
    expect(formatFileSize(1024 * 1024)).toBe("1.0 MB");
    expect(formatFileSize(5 * 1024 * 1024)).toBe("5.0 MB");
  });

  it("handles zero", () => {
    expect(formatFileSize(0)).toBe("0 B");
  });
});

describe("generateSlug", () => {
  it("converts to lowercase", () => {
    expect(generateSlug("Hello World")).toBe("hello-world");
  });

  it("replaces spaces with hyphens", () => {
    expect(generateSlug("my page title")).toBe("my-page-title");
  });

  it("removes special characters", () => {
    expect(generateSlug("Hello, World!")).toBe("hello-world");
  });

  it("handles multiple spaces", () => {
    expect(generateSlug("hello   world")).toBe("hello-world");
  });

  it("trims leading/trailing hyphens", () => {
    expect(generateSlug("  hello world  ")).toBe("hello-world");
  });

  it("handles empty string", () => {
    expect(generateSlug("")).toBe("");
  });
});

describe("truncate", () => {
  it("returns string unchanged if within limit", () => {
    expect(truncate("hello", 10)).toBe("hello");
  });

  it("truncates long strings", () => {
    expect(truncate("hello world", 8)).toBe("hello...");
  });

  it("handles exact length", () => {
    expect(truncate("hello", 5)).toBe("hello");
  });
});

describe("isExternalUrl", () => {
  it("identifies http URLs as external", () => {
    expect(isExternalUrl("http://example.com")).toBe(true);
  });

  it("identifies https URLs as external", () => {
    expect(isExternalUrl("https://example.com")).toBe(true);
  });

  it("identifies relative URLs as internal", () => {
    expect(isExternalUrl("/about")).toBe(false);
    expect(isExternalUrl("#section")).toBe(false);
    expect(isExternalUrl("about")).toBe(false);
  });
});

describe("sanitizeHtml", () => {
  it("escapes HTML entities", () => {
    expect(sanitizeHtml("<script>alert('xss')</script>")).toBe(
      "&lt;script&gt;alert(&#039;xss&#039;)&lt;/script&gt;"
    );
  });

  it("escapes ampersands", () => {
    expect(sanitizeHtml("foo & bar")).toBe("foo &amp; bar");
  });

  it("escapes quotes", () => {
    expect(sanitizeHtml('"hello"')).toBe("&quot;hello&quot;");
  });

  it("returns empty string unchanged", () => {
    expect(sanitizeHtml("")).toBe("");
  });
});
