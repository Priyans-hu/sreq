import Link from "next/link";
import { ChevronRight } from "lucide-react";

interface BreadcrumbsProps {
  slug: string[];
}

export default function Breadcrumbs({ slug }: BreadcrumbsProps) {
  if (slug.length === 0) return null;

  const crumbs = slug.map((segment, i) => ({
    label: segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, " "),
    href: `/docs/${slug.slice(0, i + 1).join("/")}`,
    isLast: i === slug.length - 1,
  }));

  return (
    <nav className="flex items-center gap-1.5 text-sm text-neutral-500 mb-6">
      <Link
        href="/docs"
        className="hover:text-neutral-300 transition-colors"
      >
        Docs
      </Link>
      {crumbs.map((crumb) => (
        <span key={crumb.href} className="flex items-center gap-1.5">
          <ChevronRight className="w-3.5 h-3.5" />
          {crumb.isLast ? (
            <span className="text-neutral-300">{crumb.label}</span>
          ) : (
            <Link
              href={crumb.href}
              className="hover:text-neutral-300 transition-colors"
            >
              {crumb.label}
            </Link>
          )}
        </span>
      ))}
    </nav>
  );
}
