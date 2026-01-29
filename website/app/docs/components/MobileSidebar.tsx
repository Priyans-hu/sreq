"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X, ChevronDown, Search } from "lucide-react";
import type { SidebarSection } from "../../lib/docs";

interface MobileSidebarProps {
  sections: SidebarSection[];
  onSearchOpen?: () => void;
}

export default function MobileSidebar({
  sections,
  onSearchOpen,
}: MobileSidebarProps) {
  const [open, setOpen] = useState(false);
  const pathname = usePathname();

  return (
    <div className="lg:hidden">
      <div className="fixed top-0 left-0 right-0 z-30 flex items-center gap-2 px-4 py-3 bg-[#0a0a0a]/90 backdrop-blur-sm border-b border-neutral-800">
        <button
          onClick={() => setOpen(true)}
          className="p-1.5 rounded-md text-neutral-400 hover:text-white transition-colors"
          aria-label="Open menu"
        >
          <Menu className="w-5 h-5" />
        </button>
        <Link href="/docs" className="text-sm font-semibold text-white">
          sreq docs
        </Link>
        {onSearchOpen && (
          <button
            onClick={onSearchOpen}
            className="ml-auto p-1.5 rounded-md text-neutral-400 hover:text-white transition-colors"
            aria-label="Search"
          >
            <Search className="w-4 h-4" />
          </button>
        )}
      </div>

      {open && (
        <>
          <div
            className="fixed inset-0 bg-black/60 z-40"
            onClick={() => setOpen(false)}
          />
          <div className="fixed inset-y-0 left-0 w-[280px] bg-[#0a0a0a] border-r border-neutral-800 z-50 overflow-y-auto">
            <div className="flex items-center justify-between p-4 border-b border-neutral-800">
              <Link
                href="/"
                className="text-lg font-bold text-white"
                onClick={() => setOpen(false)}
              >
                sreq
              </Link>
              <button
                onClick={() => setOpen(false)}
                className="p-1.5 rounded-md text-neutral-400 hover:text-white hover:bg-neutral-800 transition-colors"
                aria-label="Close menu"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <nav className="px-3 py-4">
              {sections.map((section) => (
                <MobileGroup
                  key={section.heading || section.links[0]?.href}
                  section={section}
                  pathname={pathname}
                  onNavigate={() => setOpen(false)}
                />
              ))}
            </nav>
          </div>
        </>
      )}
    </div>
  );
}

function MobileGroup({
  section,
  pathname,
  onNavigate,
}: {
  section: SidebarSection;
  pathname: string;
  onNavigate: () => void;
}) {
  const [expanded, setExpanded] = useState(true);
  const normalizedPath = pathname.replace(/\/$/, "");

  return (
    <div className="mb-1">
      {section.heading ? (
        <button
          onClick={() => setExpanded(!expanded)}
          className="flex items-center justify-between w-full px-3 py-2 mt-3 text-[11px] font-semibold uppercase tracking-widest text-neutral-500"
        >
          {section.heading}
          <ChevronDown
            className={`w-3 h-3 transition-transform duration-200 ${expanded ? "" : "-rotate-90"}`}
          />
        </button>
      ) : null}

      {expanded && (
        <ul className="space-y-0.5">
          {section.links.map((link) => {
            const linkPath = link.href.replace(/\/$/, "");
            const isActive = normalizedPath.endsWith(linkPath);

            return (
              <li key={link.href}>
                <Link
                  href={link.href}
                  onClick={onNavigate}
                  className={`block px-3 py-2 text-sm rounded-md transition-colors ${
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
