import type { Metadata } from "next";
import { parseSidebar, buildSearchIndex } from "../lib/docs";
import DocsShell from "./components/DocsShell";

export const metadata: Metadata = {
  title: "sreq docs",
  description:
    "Documentation for sreq — service-aware API client with automatic credential resolution.",
};

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const sections = parseSidebar();
  const searchIndex = buildSearchIndex();

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <DocsShell sections={sections} searchIndex={searchIndex}>
        {children}
      </DocsShell>
    </div>
  );
}
