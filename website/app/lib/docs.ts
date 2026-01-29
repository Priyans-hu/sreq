import fs from "fs";
import path from "path";
import matter from "gray-matter";
import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkGfm from "remark-gfm";
import remarkRehype from "remark-rehype";
import rehypeStringify from "rehype-stringify";
import rehypeRaw from "rehype-raw";
import { getHighlighter } from "./shiki";

// ---------- paths ----------

const DOCS_DIR = path.join(process.cwd(), "..", "docs", "content");

function resolveDocPath(slug: string[]): string | null {
  // Try slug as-is (e.g. commands/run → commands/run.md)
  const direct = path.join(DOCS_DIR, ...slug) + ".md";
  if (fs.existsSync(direct)) return direct;

  // Try as directory index (e.g. commands → commands/README.md)
  const index = path.join(DOCS_DIR, ...slug, "README.md");
  if (fs.existsSync(index)) return index;

  return null;
}

// ---------- types ----------

export interface TocItem {
  id: string;
  text: string;
  level: number;
}

export interface DocPage {
  html: string;
  toc: TocItem[];
  meta: Record<string, unknown>;
  title: string;
}

export interface SidebarLink {
  label: string;
  href: string;
}

export interface SidebarSection {
  heading: string;
  links: SidebarLink[];
}

// ---------- markdown pipeline ----------

export async function getDocBySlug(slug: string[]): Promise<DocPage | null> {
  const filePath = resolveDocPath(slug);
  if (!filePath) return null;

  const raw = fs.readFileSync(filePath, "utf-8");
  const { content, data } = matter(raw);

  const highlighter = await getHighlighter();

  // Extract TOC from markdown headings before processing
  const toc: TocItem[] = [];
  const headingRegex = /^(#{2,3})\s+(.+)$/gm;
  let match;
  while ((match = headingRegex.exec(content)) !== null) {
    const text = match[2].replace(/`([^`]+)`/g, "$1").trim();
    const id = text
      .toLowerCase()
      .replace(/[^\w\s-]/g, "")
      .replace(/\s+/g, "-");
    toc.push({ id, text, level: match[1].length });
  }

  // Process code blocks with shiki before remark
  const processedContent = await highlightCodeBlocks(content, highlighter);

  const result = await unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeRaw)
    .use(rehypeStringify)
    .process(processedContent);

  let html = String(result);

  // Add IDs to h2/h3 for TOC scroll
  html = html.replace(
    /<(h[23])>(.*?)<\/h[23]>/g,
    (_match, tag: string, inner: string) => {
      const text = inner.replace(/<[^>]+>/g, "").trim();
      const id = text
        .toLowerCase()
        .replace(/[^\w\s-]/g, "")
        .replace(/\s+/g, "-");
      return `<${tag} id="${id}">${inner}</${tag}>`;
    }
  );

  // Extract title from first h1 or frontmatter
  const titleMatch = content.match(/^#\s+(.+)$/m);
  const title =
    (data.title as string) ||
    (titleMatch ? titleMatch[1] : slug[slug.length - 1] || "Docs");

  return { html, toc, meta: data, title };
}

// ---------- code highlighting ----------

async function highlightCodeBlocks(
  markdown: string,
  highlighter: Awaited<ReturnType<typeof getHighlighter>>
): Promise<string> {
  const codeBlockRegex = /```(\w+)?\n([\s\S]*?)```/g;
  const replacements: { start: number; end: number; replacement: string }[] =
    [];

  let blockMatch;
  while ((blockMatch = codeBlockRegex.exec(markdown)) !== null) {
    const lang = blockMatch[1] || "plaintext";
    const code = blockMatch[2].trimEnd();

    try {
      const loadedLangs = highlighter.getLoadedLanguages();
      const safeLang = loadedLangs.includes(lang) ? lang : "plaintext";

      const highlighted = highlighter.codeToHtml(code, {
        lang: safeLang,
        theme: "github-dark",
      });

      replacements.push({
        start: blockMatch.index,
        end: blockMatch.index + blockMatch[0].length,
        replacement: highlighted,
      });
    } catch {
      // Fallback: leave as plain code block
    }
  }

  // Apply replacements in reverse order to preserve indices
  let result = markdown;
  for (let i = replacements.length - 1; i >= 0; i--) {
    const { start, end, replacement } = replacements[i];
    result = result.slice(0, start) + replacement + result.slice(end);
  }

  return result;
}

// ---------- sidebar ----------

export function parseSidebar(): SidebarSection[] {
  const sidebarPath = path.join(DOCS_DIR, "_sidebar.md");
  if (!fs.existsSync(sidebarPath)) return [];

  const raw = fs.readFileSync(sidebarPath, "utf-8");
  const sections: SidebarSection[] = [];
  let currentSection: SidebarSection = { heading: "", links: [] };

  for (const line of raw.split("\n")) {
    const trimmed = line.trim();

    // Bold heading: - **Commands**
    const headingMatch = trimmed.match(/^-\s+\*\*(.+?)\*\*$/);
    if (headingMatch) {
      if (currentSection.heading || currentSection.links.length > 0) {
        sections.push(currentSection);
      }
      currentSection = { heading: headingMatch[1], links: [] };
      continue;
    }

    // Link: - [Label](/path) or  - [Label](/path)
    const linkMatch = trimmed.match(/^-\s+\[(.+?)]\((.+?)\)$/);
    if (linkMatch) {
      const href = linkMatch[2];
      // Skip external links
      if (href.startsWith("http")) continue;

      // Transform docsify paths to Next.js paths
      const docPath = href
        .replace(/^\//, "") // remove leading slash
        .replace(/\/$/, ""); // remove trailing slash

      const nextHref = docPath ? `/docs/${docPath}` : "/docs";

      currentSection.links.push({
        label: linkMatch[1],
        href: nextHref,
      });
    }
  }

  if (currentSection.heading || currentSection.links.length > 0) {
    sections.push(currentSection);
  }

  return sections;
}

// ---------- static params ----------

export function getAllDocSlugs(): string[][] {
  const slugs: string[][] = [];

  function walk(dir: string, prefix: string[] = []) {
    const entries = fs.readdirSync(dir, { withFileTypes: true });

    for (const entry of entries) {
      if (entry.name.startsWith("_")) continue;

      if (entry.isDirectory()) {
        walk(path.join(dir, entry.name), [...prefix, entry.name]);
      } else if (entry.name.endsWith(".md")) {
        if (entry.name === "README.md") {
          slugs.push(prefix.length > 0 ? [...prefix] : []);
        } else {
          slugs.push([...prefix, entry.name.replace(/\.md$/, "")]);
        }
      }
    }
  }

  walk(DOCS_DIR);
  return slugs;
}

// ---------- flat list for prev/next ----------

export function getFlatDocList(): SidebarLink[] {
  const sections = parseSidebar();
  return sections.flatMap((s) => s.links);
}

// ---------- search index ----------

export interface SearchEntry {
  title: string;
  href: string;
  section: string;
  content: string; // plain text snippet for searching
}

export function buildSearchIndex(): SearchEntry[] {
  const sections = parseSidebar();
  const entries: SearchEntry[] = [];

  for (const section of sections) {
    for (const link of section.links) {
      // Derive slug from href: "/docs/commands/run" → ["commands", "run"]
      const slug = link.href
        .replace(/^\/docs\/?/, "")
        .split("/")
        .filter(Boolean);

      const filePath = resolveDocPath(slug);
      if (!filePath) continue;

      const raw = fs.readFileSync(filePath, "utf-8");
      const { content } = matter(raw);

      // Strip markdown syntax for plain text search
      const plainText = content
        .replace(/```[\s\S]*?```/g, "") // code blocks
        .replace(/`[^`]+`/g, "") // inline code
        .replace(/#{1,6}\s+/g, "") // headings
        .replace(/\[([^\]]+)]\([^)]+\)/g, "$1") // links
        .replace(/[*_~]+/g, "") // bold/italic/strikethrough
        .replace(/\|[^\n]+/g, "") // table rows
        .replace(/-{3,}/g, "") // horizontal rules
        .replace(/>\s+/g, "") // blockquotes
        .replace(/\n{2,}/g, " ")
        .replace(/\n/g, " ")
        .trim()
        .slice(0, 500); // limit content size

      entries.push({
        title: link.label,
        href: link.href,
        section: section.heading || "Getting Started",
        content: plainText,
      });
    }
  }

  return entries;
}
