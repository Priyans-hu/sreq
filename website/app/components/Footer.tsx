import { Github } from "lucide-react";

const links = [
  { label: "GitHub", href: "https://github.com/Priyans-hu/sreq" },
  { label: "Releases", href: "https://github.com/Priyans-hu/sreq/releases" },
  { label: "Issues", href: "https://github.com/Priyans-hu/sreq/issues" },
  {
    label: "Documentation",
    href: "https://github.com/Priyans-hu/sreq#usage",
  },
];

export default function Footer() {
  return (
    <footer className="border-t border-neutral-800 py-12 px-6">
      <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-3">
          <span className="text-lg font-bold tracking-tight">sreq</span>
          <span className="text-neutral-600">|</span>
          <span className="text-sm text-neutral-500">MIT License</span>
        </div>

        <nav className="flex items-center gap-6">
          {links.map((link) => (
            <a
              key={link.label}
              href={link.href}
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm text-neutral-500 hover:text-neutral-300 transition-colors"
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-2 text-sm text-neutral-500">
          Built by{" "}
          <a
            href="https://github.com/Priyans-hu"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 text-neutral-300 hover:text-emerald-400 transition-colors"
          >
            <Github className="w-3.5 h-3.5" />
            Priyanshu
          </a>
        </div>
      </div>
    </footer>
  );
}
