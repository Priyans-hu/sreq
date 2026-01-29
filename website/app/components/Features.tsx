"use client";

import { motion } from "framer-motion";
import { Shield, Globe, Layers, Terminal, History, Lock } from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Auto Credentials",
    description:
      "Fetches secrets from Consul, AWS Secrets Manager, env vars, and dotenv files. No more manual token juggling.",
  },
  {
    icon: Globe,
    title: "Service Discovery",
    description:
      "Resolves service URLs from providers automatically. Just specify the service name and environment.",
  },
  {
    icon: Layers,
    title: "Multi-Environment",
    description:
      "Switch between dev, staging, and prod with a single flag. Each environment gets its own credential chain.",
  },
  {
    icon: Terminal,
    title: "Interactive TUI",
    description:
      "Browse services, build requests, and view history in a terminal UI. Or use the familiar CLI syntax.",
  },
  {
    icon: History,
    title: "Request History",
    description:
      "Track, replay, and export previous requests as curl or HTTPie commands. Never lose a request again.",
  },
  {
    icon: Lock,
    title: "Secure Caching",
    description:
      "AES-256 encrypted local credential cache with TTL. Supports offline mode for working without network access.",
  },
];

const containerVariants = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const cardVariants = {
  hidden: { opacity: 0, y: 30 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, ease: "easeOut" as const },
  },
};

export default function Features() {
  return (
    <section className="py-24 px-6" id="features">
      <div className="max-w-6xl mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.6 }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            Everything you need
          </h2>
          <p className="text-neutral-400 max-w-xl mx-auto">
            Built for teams that manage microservices across multiple
            environments.
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-100px" }}
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          {features.map((feature) => (
            <motion.div
              key={feature.title}
              variants={cardVariants}
              className="group p-6 rounded-xl bg-neutral-900/50 border border-neutral-800 hover:border-emerald-500/30 transition-colors"
            >
              <div className="w-10 h-10 rounded-lg bg-emerald-500/10 flex items-center justify-center mb-4 group-hover:bg-emerald-500/20 transition-colors">
                <feature.icon className="w-5 h-5 text-emerald-400" />
              </div>
              <h3 className="text-lg font-semibold mb-2">{feature.title}</h3>
              <p className="text-neutral-400 text-sm leading-relaxed">
                {feature.description}
              </p>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
