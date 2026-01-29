import Link from "next/link";
import { ChevronLeft, ChevronRight } from "lucide-react";
import type { SidebarLink } from "../../lib/docs";

interface PrevNextLinksProps {
  prev: SidebarLink | null;
  next: SidebarLink | null;
}

export default function PrevNextLinks({ prev, next }: PrevNextLinksProps) {
  if (!prev && !next) return null;

  return (
    <div className="flex items-center justify-between mt-12 pt-6 border-t border-neutral-800">
      {prev ? (
        <Link
          href={prev.href}
          className="group flex items-center gap-2 text-sm text-neutral-400 hover:text-emerald-400 transition-colors"
        >
          <ChevronLeft className="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" />
          {prev.label}
        </Link>
      ) : (
        <div />
      )}
      {next ? (
        <Link
          href={next.href}
          className="group flex items-center gap-2 text-sm text-neutral-400 hover:text-emerald-400 transition-colors"
        >
          {next.label}
          <ChevronRight className="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
        </Link>
      ) : (
        <div />
      )}
    </div>
  );
}
