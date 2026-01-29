"use client";

import { motion } from "framer-motion";
import { Terminal, FileText, Zap } from "lucide-react";

const steps = [
  {
    icon: Terminal,
    title: "Install",
    description: "One command via Homebrew, curl, or go install.",
  },
  {
    icon: FileText,
    title: "Configure",
    description:
      "Point sreq at your service registry. Or just use environment variables.",
  },
  {
    icon: Zap,
    title: "Request",
    description:
      "Make requests with automatic credential injection. No tokens to copy-paste.",
  },
];

const containerVariants = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.15,
    },
  },
};

const stepVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, ease: "easeOut" as const },
  },
};

export default function HowItWorks() {
  return (
    <section className="py-24 px-6" id="how-it-works">
      <div className="max-w-5xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6 }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            Get started in seconds
          </h2>
          <p className="text-neutral-400 max-w-xl mx-auto">
            Three simple steps to never copy-paste credentials again.
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-100px" }}
          className="grid grid-cols-1 md:grid-cols-3 gap-8 relative"
        >
          {/* Dashed connector lines (desktop only) */}
          <div className="hidden md:block absolute top-12 left-1/3 w-1/3 border-t-2 border-dashed border-neutral-800" />
          <div className="hidden md:block absolute top-12 right-1/3 w-1/3 border-t-2 border-dashed border-neutral-800" />

          {steps.map((step, index) => (
            <motion.div
              key={step.title}
              variants={stepVariants}
              className="text-center"
            >
              <div className="relative inline-flex items-center justify-center w-24 h-24 rounded-full bg-neutral-900 border border-neutral-800 mb-6">
                <step.icon className="w-8 h-8 text-emerald-400" />
                <span className="absolute -top-2 -right-2 w-7 h-7 rounded-full bg-emerald-500 text-black text-xs font-bold flex items-center justify-center">
                  {index + 1}
                </span>
              </div>
              <h3 className="text-xl font-semibold mb-2">{step.title}</h3>
              <p className="text-neutral-400 text-sm max-w-xs mx-auto">
                {step.description}
              </p>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
