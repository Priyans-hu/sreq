"use client";

import { useState, useEffect, useCallback } from "react";
import type { SidebarSection, SearchEntry } from "../../lib/docs";
import Sidebar from "./Sidebar";
import MobileSidebar from "./MobileSidebar";
import SearchModal from "./SearchModal";

interface DocsShellProps {
  sections: SidebarSection[];
  searchIndex: SearchEntry[];
  children: React.ReactNode;
}

export default function DocsShell({
  sections,
  searchIndex,
  children,
}: DocsShellProps) {
  const [searchOpen, setSearchOpen] = useState(false);

  const openSearch = useCallback(() => setSearchOpen(true), []);
  const closeSearch = useCallback(() => setSearchOpen(false), []);

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        setSearchOpen((prev) => !prev);
      }
    }
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <>
      <Sidebar sections={sections} onSearchOpen={openSearch} />
      <MobileSidebar sections={sections} onSearchOpen={openSearch} />
      <main className="lg:ml-[260px] xl:mr-[200px] min-h-screen">
        <div className="max-w-[750px] mx-auto px-6 py-16 lg:py-12">
          {children}
        </div>
      </main>
      <SearchModal
        entries={searchIndex}
        open={searchOpen}
        onClose={closeSearch}
      />
    </>
  );
}
