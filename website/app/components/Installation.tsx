"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { Copy, Check } from "lucide-react";

const tabs = [
  {
    label: "Homebrew",
    code: "brew install Priyans-hu/tap/sreq",
  },
  {
    label: "Script",
    code: "curl -fsSL https://raw.githubusercontent.com/Priyans-hu/sreq/main/install.sh | bash",
  },
  {
    label: "Go",
    code: "go install github.com/Priyans-hu/sreq/cmd/sreq@latest",
  },
  {
    label: "Binary",
    code: "# Download from GitHub Releases\nhttps://github.com/Priyans-hu/sreq/releases/latest",
  },
];

export default function Installation() {
  const [activeTab, setActiveTab] = useState(0);
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(tabs[activeTab].code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <section className="py-24 px-6" id="installation">
      <div className="max-w-3xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6 }}
          className="text-center mb-12"
        >
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">Installation</h2>
          <p className="text-neutral-400">
            Pick your preferred method and get started.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6, delay: 0.1 }}
        >
          <div className="rounded-xl border border-neutral-800 bg-neutral-950 overflow-hidden">
            {/* Tabs */}
            <div className="flex border-b border-neutral-800">
              {tabs.map((tab, index) => (
                <button
                  key={tab.label}
                  onClick={() => {
                    setActiveTab(index);
                    setCopied(false);
                  }}
                  className={`px-5 py-3 text-sm font-medium transition-colors relative ${
                    activeTab === index
                      ? "text-emerald-400"
                      : "text-neutral-500 hover:text-neutral-300"
                  }`}
                >
                  {tab.label}
                  {activeTab === index && (
                    <motion.div
                      layoutId="tab-underline"
                      className="absolute bottom-0 left-0 right-0 h-0.5 bg-emerald-400"
                    />
                  )}
                </button>
              ))}
            </div>

            {/* Code block */}
            <div className="p-6 flex items-start justify-between gap-4">
              <pre className="font-mono text-sm text-neutral-300 whitespace-pre-wrap break-all flex-1">
                {tabs[activeTab].code}
              </pre>
              <button
                onClick={handleCopy}
                className="p-2 rounded-md hover:bg-neutral-800 transition-colors text-neutral-500 hover:text-neutral-300 shrink-0"
                aria-label="Copy code"
              >
                {copied ? (
                  <Check className="w-4 h-4 text-emerald-400" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
