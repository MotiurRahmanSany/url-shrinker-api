import ShrinkForm from "./components/ShrinkForm";
import { Zap, Shield, BarChart3, Globe } from "lucide-react";

export default function Home() {
  return (
    <main className="flex-1">
      {/* Hero Section */}
      <section className="py-20 px-8 text-center space-y-12 max-w-5xl mx-auto">
        <div className="space-y-6">
          <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight">
            Simplify Your Links, <br />
            <span className="text-gradient">Amplify Your Impact.</span>
          </h1>
          <p className="text-xl text-slate-400 max-w-2xl mx-auto leading-relaxed">
            A production-grade URL shortener built for speed, transparency, and deep analytics. 
            Shrink your long URLs and track every click in real-time.
          </p>
        </div>

        <ShrinkForm />

        {/* Feature Highlights */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 pt-20">
          <FeatureCard 
            icon={<Zap className="w-6 h-6 text-indigo-400" />}
            title="Lightning Fast"
            description="Built with Go and Redis for sub-millisecond redirect performance."
          />
          <FeatureCard 
            icon={<BarChart3 className="w-6 h-6 text-cyan-400" />}
            title="Live Analytics"
            description="Track clicks, referrers, and audience engagement metrics."
          />
          <FeatureCard 
            icon={<Shield className="w-6 h-6 text-emerald-400" />}
            title="Secure & Reliable"
            description="JWT-based authentication and rate limiting to keep your data safe."
          />
        </div>
      </section>

      {/* Social Proof / Stats */}
      <section className="border-t border-slate-800/50 bg-slate-900/20 py-16 px-8">
        <div className="max-w-6xl mx-auto flex flex-wrap justify-center gap-12 md:gap-24 opacity-50 grayscale hover:grayscale-0 transition-all duration-700">
           <div className="flex items-center gap-2 font-bold text-2xl">
             <Globe className="w-8 h-8" /> <span>Global Scale</span>
           </div>
           <div className="flex items-center gap-2 font-bold text-2xl text-indigo-400">
             <span>99.9% Uptime</span>
           </div>
           <div className="flex items-center gap-2 font-bold text-2xl">
             <span>Zero Latency</span>
           </div>
        </div>
      </section>
    </main>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="premium-card text-left group">
      <div className="inline-flex p-3 rounded-xl bg-slate-800/50 mb-4 group-hover:scale-110 transition-transform">
        {icon}
      </div>
      <h3 className="text-xl font-bold mb-2">{title}</h3>
      <p className="text-slate-400 leading-relaxed text-sm">
        {description}
      </p>
    </div>
  );
}
