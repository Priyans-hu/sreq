"use client";

import { motion } from "framer-motion";

const lines = [
  { text: "$ sreq POST /api/v1/users -s auth-service -e dev -d '{\"name\":\"test\"}'", color: "text-emerald-400" },
  { text: "", color: "" },
  { text: "\u2192 Resolving auth-service in dev...", color: "text-neutral-500" },
  { text: "\u2192 Found: https://auth.dev.internal:8443", color: "text-cyan-400" },
  { text: "\u2192 Fetching credentials from AWS Secrets Manager...", color: "text-neutral-500" },
  { text: "\u2192 Using token: sk-...redacted", color: "text-neutral-500" },
  { text: "", color: "" },
  { text: "HTTP/1.1 201 Created", color: "text-yellow-400" },
  { text: "Content-Type: application/json", color: "text-neutral-500" },
  { text: "", color: "" },
  { text: "{", color: "text-neutral-300" },
  { text: '  "id": "usr_2kF9xL",', color: "text-neutral-300" },
  { text: '  "name": "test",', color: "text-neutral-300" },
  { text: '  "created_at": "2026-01-29T10:00:00Z"', color: "text-neutral-300" },
  { text: "}", color: "text-neutral-300" },
];

export default function Demo() {
  return (
    <section className="py-24 px-6" id="demo">
      <div className="max-w-4xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6 }}
          className="text-center mb-12"
        >
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            See it in action
          </h2>
          <p className="text-neutral-400 max-w-xl mx-auto">
            One command. Automatic service discovery and credential injection.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6, ease: "easeOut" }}
        >
          <div className="rounded-xl border border-neutral-800 bg-neutral-950 overflow-hidden shadow-2xl shadow-emerald-500/5">
            {/* Title bar */}
            <div className="flex items-center gap-2 px-4 py-3 bg-neutral-900 border-b border-neutral-800">
              <div className="w-3 h-3 rounded-full bg-red-500/80" />
              <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
              <div className="w-3 h-3 rounded-full bg-green-500/80" />
              <span className="ml-3 text-xs text-neutral-500 font-mono">
                terminal
              </span>
            </div>

            {/* Terminal body */}
            <div className="p-6 font-mono text-sm leading-7 overflow-x-auto">
              {lines.map((line, i) => (
                <div key={i} className={line.color || "h-5"}>
                  {line.text}
                </div>
              ))}
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
