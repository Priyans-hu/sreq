"use client";

import { motion } from "framer-motion";
import { Shield, Globe, Layers, Terminal, Zap, Lock } from "lucide-react";

const features = [
  {
    icon: Shield,
    title: "Auto Credentials",
    description:
      "Fetches secrets from AWS Secrets Manager, Consul, or HashiCorp Vault. No more manual token juggling.",
  },
  {
    icon: Globe,
    title: "Service Discovery",
    description:
      "Resolves service URLs from Consul automatically. Just specify the service name and environment.",
  },
  {
    icon: Layers,
    title: "Multi-Environment",
    description:
      "Switch between dev, staging, and prod with a single flag. Each environment gets its own credential chain.",
  },
  {
    icon: Terminal,
    title: "CLI Native",
    description:
      "Familiar curl-like syntax. Drop-in replacement for your API testing workflow.",
  },
  {
    icon: Zap,
    title: "Zero Config",
    description:
      "Works out of the box with sensible defaults. Optional .sreq.yaml for advanced setups.",
  },
  {
    icon: Lock,
    title: "Secure by Default",
    description:
      "Credentials are fetched at request time, never stored locally. Supports mTLS and token rotation.",
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
