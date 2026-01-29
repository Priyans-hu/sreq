"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronDown, ArrowLeft, Search } from "lucide-react";
import { useState } from "react";
import type { SidebarSection } from "../../lib/docs";

interface SidebarProps {
  sections: SidebarSection[];
  onSearchOpen?: () => void;
}

export default function Sidebar({ sections, onSearchOpen }: SidebarProps) {
  const pathname = usePathname();

  return (
    <aside className="hidden lg:flex flex-col fixed left-0 top-0 bottom-0 w-[260px] border-r border-neutral-800 bg-[#0a0a0a] z-20">
      <div className="px-5 pt-5 pb-2">
        <div className="flex items-center justify-between mb-4">
          <Link
            href="/"
            className="text-lg font-bold tracking-tight text-white hover:text-emerald-400 transition-colors"
          >
            sreq
          </Link>
          <Link
            href="/"
            className="flex items-center gap-1 text-xs text-neutral-500 hover:text-neutral-300 transition-colors"
          >
            <ArrowLeft className="w-3 h-3" />
            Home
          </Link>
        </div>

        {onSearchOpen && (
          <button
            onClick={onSearchOpen}
            className="flex items-center gap-2 w-full px-3 py-2 text-sm text-neutral-500 bg-neutral-900 border border-neutral-800 rounded-lg hover:border-neutral-700 hover:text-neutral-400 transition-colors"
          >
            <Search className="w-3.5 h-3.5" />
            <span>Search docs...</span>
            <kbd className="ml-auto text-xs text-neutral-600 bg-neutral-800 px-1.5 py-0.5 rounded">
              ⌘K
            </kbd>
          </button>
        )}
      </div>

      <nav className="flex-1 overflow-y-auto px-3 pb-4 pt-2">
        {sections.map((section) => (
          <SidebarGroup
            key={section.heading || section.links[0]?.href}
            section={section}
            pathname={pathname}
          />
        ))}
      </nav>

      <div className="px-5 py-3 border-t border-neutral-800 text-xs text-neutral-600">
        <a
          href="https://github.com/Priyans-hu/sreq"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-neutral-400 transition-colors"
        >
          GitHub
        </a>
        <span className="mx-2">·</span>
        <a
          href="https://github.com/Priyans-hu/sreq/releases"
          target="_blank"
          rel="noopener noreferrer"
          className="hover:text-neutral-400 transition-colors"
        >
          Releases
        </a>
      </div>
    </aside>
  );
}

function SidebarGroup({
  section,
  pathname,
}: {
  section: SidebarSection;
  pathname: string;
}) {
  // pathname from usePathname() includes basePath, so normalize it
  const normalizedPath = pathname.replace(/\/$/, "");

  const hasActiveChild = section.links.some((link) => {
    const linkPath = link.href.replace(/\/$/, "");
    return normalizedPath.endsWith(linkPath);
  });

  const [open, setOpen] = useState(hasActiveChild || !section.heading);

  return (
    <div className="mb-1">
      {section.heading ? (
        <button
          onClick={() => setOpen(!open)}
          className="flex items-center justify-between w-full px-3 py-2 mt-3 text-[11px] font-semibold uppercase tracking-widest text-neutral-500 hover:text-neutral-300 transition-colors"
        >
          {section.heading}
          <ChevronDown
            className={`w-3 h-3 transition-transform duration-200 ${open ? "" : "-rotate-90"}`}
          />
        </button>
      ) : null}

      {open && (
        <ul className="space-y-0.5">
          {section.links.map((link) => {
            const linkPath = link.href.replace(/\/$/, "");
            const isActive = normalizedPath.endsWith(linkPath);

            return (
              <li key={link.href}>
                <Link
                  href={link.href}
                  className={`block px-3 py-1.5 text-[13px] rounded-md transition-colors ${
                    isActive
                      ? "text-emerald-400 bg-emerald-500/10 font-medium"
                      : "text-neutral-400 hover:text-neutral-200 hover:bg-neutral-800/50"
                  }`}
                >
                  {link.label}
                </Link>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
