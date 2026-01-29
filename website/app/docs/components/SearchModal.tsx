"use client";

import { useMemo } from "react";
import type { SearchEntry } from "../../lib/docs";

interface SearchModalProps {
  entries: SearchEntry[];
  open: boolean;
  onClose: () => void;
}

export default function SearchModal({
  entries,
  open,
  onClose,
}: SearchModalProps) {
  if (!open) return null;

  // Key forces remount on each open, resetting all state
  return <SearchModalInner key="open" entries={entries} onClose={onClose} />;
}

import { useState, useEffect, useRef, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Search, X, FileText, ArrowRight } from "lucide-react";

function SearchModalInner({
  entries,
  onClose,
}: {
  entries: SearchEntry[];
  onClose: () => void;
}) {
  const [query, setQuery] = useState("");
  const [activeIndex, setActiveIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  const results = useMemo(
    () => (query.length >= 2 ? doSearch(entries, query) : []),
    [entries, query]
  );

  // Focus input on mount
  useEffect(() => {
    const timer = setTimeout(() => inputRef.current?.focus(), 50);
    return () => clearTimeout(timer);
  }, []);

  const navigate = useCallback(
    (href: string) => {
      router.push(href);
      onClose();
    },
    [router, onClose]
  );

  const handleQueryChange = (value: string) => {
    setQuery(value);
    setActiveIndex(0);
  };

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") {
        onClose();
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveIndex((i) => Math.min(i + 1, results.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveIndex((i) => Math.max(i - 1, 0));
      } else if (e.key === "Enter") {
        // Access results via closure
        setActiveIndex((currentIndex) => {
          const current = results[currentIndex];
          if (current) navigate(current.href);
          return currentIndex;
        });
      }
    }

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [onClose, results, navigate]);

  return (
    <>
      <div
        className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
        onClick={onClose}
      />
      <div className="fixed top-[15%] left-1/2 -translate-x-1/2 w-full max-w-lg z-50">
        <div className="mx-4 bg-neutral-900 border border-neutral-700 rounded-xl shadow-2xl overflow-hidden">
          <div className="flex items-center gap-3 px-4 border-b border-neutral-800">
            <Search className="w-4 h-4 text-neutral-500 shrink-0" />
            <input
              ref={inputRef}
              value={query}
              onChange={(e) => handleQueryChange(e.target.value)}
              placeholder="Search documentation..."
              className="flex-1 py-3.5 bg-transparent text-sm text-white placeholder-neutral-500 outline-none"
            />
            <button
              onClick={onClose}
              className="p-1 rounded text-neutral-500 hover:text-neutral-300 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          </div>

          <div className="max-h-[360px] overflow-y-auto">
            {query.length < 2 ? (
              <div className="px-4 py-8 text-center text-sm text-neutral-500">
                Type at least 2 characters to search
              </div>
            ) : results.length === 0 ? (
              <div className="px-4 py-8 text-center text-sm text-neutral-500">
                No results for &ldquo;{query}&rdquo;
              </div>
            ) : (
              <ul className="py-2">
                {results.map((result, i) => (
                  <li key={result.href}>
                    <button
                      onClick={() => navigate(result.href)}
                      onMouseEnter={() => setActiveIndex(i)}
                      className={`flex items-start gap-3 w-full px-4 py-3 text-left transition-colors ${
                        i === activeIndex
                          ? "bg-emerald-500/10"
                          : "hover:bg-neutral-800/50"
                      }`}
                    >
                      <FileText
                        className={`w-4 h-4 mt-0.5 shrink-0 ${
                          i === activeIndex
                            ? "text-emerald-400"
                            : "text-neutral-500"
                        }`}
                      />
                      <div className="flex-1 min-w-0">
                        <div
                          className={`text-sm font-medium ${
                            i === activeIndex
                              ? "text-emerald-400"
                              : "text-neutral-200"
                          }`}
                        >
                          {result.title}
                        </div>
                        <div className="text-xs text-neutral-500 mt-0.5">
                          {result.section}
                        </div>
                        {result.snippet && (
                          <div className="text-xs text-neutral-500 mt-1 line-clamp-2">
                            {result.snippet}
                          </div>
                        )}
                      </div>
                      {i === activeIndex && (
                        <ArrowRight className="w-3.5 h-3.5 text-emerald-400 mt-0.5 shrink-0" />
                      )}
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <div className="flex items-center gap-4 px-4 py-2.5 border-t border-neutral-800 text-[11px] text-neutral-600">
            <span className="flex items-center gap-1">
              <kbd className="px-1 py-0.5 bg-neutral-800 rounded text-neutral-500">
                ↑↓
              </kbd>
              navigate
            </span>
            <span className="flex items-center gap-1">
              <kbd className="px-1 py-0.5 bg-neutral-800 rounded text-neutral-500">
                ↵
              </kbd>
              open
            </span>
            <span className="flex items-center gap-1">
              <kbd className="px-1 py-0.5 bg-neutral-800 rounded text-neutral-500">
                esc
              </kbd>
              close
            </span>
          </div>
        </div>
      </div>
    </>
  );
}

interface SearchResult {
  title: string;
  href: string;
  section: string;
  snippet: string;
}

function doSearch(entries: SearchEntry[], query: string): SearchResult[] {
  const q = query.toLowerCase();
  const scored: (SearchResult & { score: number })[] = [];

  for (const entry of entries) {
    const titleLower = entry.title.toLowerCase();
    const contentLower = entry.content.toLowerCase();

    let score = 0;

    if (titleLower === q) {
      score += 100;
    } else if (titleLower.startsWith(q)) {
      score += 50;
    } else if (titleLower.includes(q)) {
      score += 30;
    }

    const contentIndex = contentLower.indexOf(q);
    if (contentIndex >= 0) {
      score += 10;
    }

    if (score === 0) continue;

    let snippet = "";
    if (contentIndex >= 0) {
      const start = Math.max(0, contentIndex - 40);
      const end = Math.min(entry.content.length, contentIndex + q.length + 80);
      snippet =
        (start > 0 ? "..." : "") +
        entry.content.slice(start, end).trim() +
        (end < entry.content.length ? "..." : "");
    }

    scored.push({
      title: entry.title,
      href: entry.href,
      section: entry.section,
      snippet,
      score,
    });
  }

  return scored.sort((a, b) => b.score - a.score).slice(0, 10);
}
