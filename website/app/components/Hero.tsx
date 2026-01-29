"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { Copy, Check, ExternalLink, BookOpen } from "lucide-react";

export default function Hero() {
  const [copied, setCopied] = useState(false);
  const installCmd = "brew install Priyans-hu/tap/sreq";

  const handleCopy = async () => {
    await navigator.clipboard.writeText(installCmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden px-6">
      {/* Background glow */}
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] bg-emerald-500/10 rounded-full blur-[120px]" />
        <div className="absolute top-1/3 left-1/3 w-[400px] h-[400px] bg-cyan-500/8 rounded-full blur-[100px]" />
      </div>

      <div className="relative z-10 max-w-4xl mx-auto text-center">
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, ease: "easeOut" }}
        >
          <div className="inline-flex items-center gap-2 px-4 py-1.5 mb-8 rounded-full border border-emerald-500/20 bg-emerald-500/5 text-emerald-400 text-sm">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            Open source CLI tool
          </div>

          <h1 className="text-5xl sm:text-6xl md:text-7xl font-bold tracking-tight leading-[1.1] mb-6">
            API requests without{" "}
            <span className="bg-gradient-to-r from-emerald-400 to-cyan-400 bg-clip-text text-transparent">
              the credential hassle
            </span>
          </h1>

          <p className="text-lg sm:text-xl text-neutral-400 max-w-2xl mx-auto mb-10 leading-relaxed">
            sreq automatically resolves service URLs and credentials from Consul,
            AWS Secrets Manager, and more — so you can focus on building, not
            configuring.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.2, ease: "easeOut" }}
          className="mb-10"
        >
          <div className="inline-flex items-center gap-3 bg-neutral-900 border border-neutral-800 rounded-lg px-5 py-3 font-mono text-sm">
            <span className="text-emerald-400">$</span>
            <span className="text-neutral-300">{installCmd}</span>
            <button
              onClick={handleCopy}
              className="ml-2 p-1.5 rounded-md hover:bg-neutral-800 transition-colors text-neutral-500 hover:text-neutral-300"
              aria-label="Copy install command"
            >
              {copied ? (
                <Check className="w-4 h-4 text-emerald-400" />
              ) : (
                <Copy className="w-4 h-4" />
              )}
            </button>
          </div>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.4, ease: "easeOut" }}
          className="flex flex-col sm:flex-row items-center justify-center gap-4"
        >
          <a
            href="https://github.com/Priyans-hu/sreq"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-black font-semibold transition-colors"
          >
            <ExternalLink className="w-4 h-4" />
            View on GitHub
          </a>
          <a
            href="/sreq/docs/"
            className="inline-flex items-center gap-2 px-6 py-3 rounded-lg border border-neutral-700 hover:border-neutral-500 text-neutral-300 hover:text-white font-semibold transition-colors"
          >
            <BookOpen className="w-4 h-4" />
            Documentation
          </a>
        </motion.div>
      </div>
    </section>
  );
}
