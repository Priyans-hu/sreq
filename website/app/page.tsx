import Hero from "./components/Hero";
import Features from "./components/Features";
import Demo from "./components/Demo";
import HowItWorks from "./components/HowItWorks";
import Installation from "./components/Installation";
import Footer from "./components/Footer";

export default function Home() {
  return (
    <main>
      <Hero />
      <Features />
      <Demo />
      <HowItWorks />
      <Installation />
      <Footer />
    </main>
  );
}
