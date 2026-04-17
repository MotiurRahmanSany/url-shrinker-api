"use client";

import { useEffect, useState, use } from "react";
import { useAuth } from "../../context/AuthContext";
import { useRouter } from "next/navigation";
import { apiRequest, BASE_URL } from "@/lib/api";
import { 
  Loader2, 
  ArrowLeft, 
  BarChart2, 
  MousePointer2, 
  Calendar,
  ExternalLink,
  History
} from "lucide-react";
import Link from "next/link";

type ClickStat = {
  day: string;
  total: number;
};

type URLStats = {
  total_clicks: number;
  clicks_today: number;
  daily_timeline: ClickStat[];
};

type URLInfo = {
  id: string;
  original_url: string;
  short_code: string;
  short_url: string;
  created_at: string;
};

export default function Stats({ params }: { params: Promise<{ code: string }> }) {
  const { code } = use(params);
  const { user, loading: authLoading } = useAuth();
  const [stats, setStats] = useState<URLStats | null>(null);
  const [urlInfo, setUrlInfo] = useState<URLInfo | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]);

  const fetchData = async () => {
    setLoading(true);
    
    // Fetch URL info
    const infoRes = await apiRequest<URLInfo>(`/urls/${code}`);
    if (infoRes.success && infoRes.data) {
      setUrlInfo(infoRes.data);
    } else {
      setError("Failed to load URL information");
      setLoading(false);
      return;
    }

    // Fetch Stats
    const statsRes = await apiRequest<URLStats>(`/urls/${code}/stats`);
    if (statsRes.success && statsRes.data) {
      setStats(statsRes.data);
    } else {
      setError("Failed to load analytics data");
    }
    
    setLoading(false);
  };

  useEffect(() => {
    if (user && code) {
      fetchData();
    }
  }, [user, code]);

  if (authLoading || loading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <Loader2 className="w-10 h-10 animate-spin text-indigo-500" />
      </div>
    );
  }

  if (error || !urlInfo || !stats) {
    return (
      <div className="flex-1 flex flex-col items-center justify-center p-6 space-y-4">
        <div className="bg-red-500/10 p-4 rounded-full text-red-500">
           <ArrowLeft className="w-8 h-8" />
        </div>
        <h2 className="text-2xl font-bold">Something went wrong</h2>
        <p className="text-slate-400">{error || "Could not find statistics for this link."}</p>
        <Link href="/dashboard" className="text-indigo-400 font-bold hover:underline">Return to Dashboard</Link>
      </div>
    );
  }

  const maxClicks = Math.max(...stats.daily_timeline.map(d => d.total), 1);

  return (
    <div className="flex-1 p-8 max-w-5xl mx-auto w-full space-y-8">
      <header className="space-y-4">
        <Link href="/dashboard" className="flex items-center gap-2 text-sm text-slate-400 hover:text-white transition-colors">
          <ArrowLeft className="w-4 h-4" />
          Back to Dashboard
        </Link>
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div className="space-y-1">
            <h1 className="text-3xl font-bold text-gradient">Link Analytics</h1>
            <p className="text-slate-400 break-all">{urlInfo.original_url}</p>
          </div>
          <div className="flex items-center gap-2 px-4 py-2 bg-indigo-500/10 border border-indigo-500/20 rounded-xl">
             <span className="text-indigo-400 font-bold">{`${BASE_URL}/${urlInfo.short_code}`}</span>
             <a href={`${BASE_URL}/${urlInfo.short_code}`} target="_blank"><ExternalLink className="w-4 h-4 text-slate-500 hover:text-indigo-400" /></a>
          </div>
        </div>
      </header>

      {/* Main Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div className="glass p-8 rounded-3xl flex flex-col items-center justify-center text-center space-y-2 border-indigo-500/10">
           <div className="p-4 bg-indigo-500/20 rounded-2xl text-indigo-400 mb-2">
             <MousePointer2 className="w-8 h-8" />
           </div>
           <p className="text-3xl md:text-5xl font-black">{stats.total_clicks}</p>
           <p className="text-slate-400 font-semibold uppercase tracking-widest text-xs">Total Clicks</p>
        </div>
        <div className="glass p-8 rounded-3xl flex flex-col items-center justify-center text-center space-y-2 border-cyan-500/10">
           <div className="p-4 bg-cyan-500/20 rounded-2xl text-cyan-400 mb-2">
             <Zap className="w-8 h-8" />
           </div>
           <p className="text-3xl md:text-5xl font-black">{stats.clicks_today}</p>
           <p className="text-slate-400 font-semibold uppercase tracking-widest text-xs">Clicks Today</p>
        </div>
      </div>

      {/* Daily Activity Chart (CSS Bar Chart) */}
      <div className="premium-card space-y-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-slate-800 rounded-lg"><History className="w-5 h-5 text-indigo-400" /></div>
            <h2 className="text-xl font-bold">Daily Activity</h2>
          </div>
          <span className="text-xs text-slate-500 font-mono">Last 7 Entries</span>
        </div>

        {stats.daily_timeline.length === 0 ? (
          <div className="py-20 text-center text-slate-500">
            No activity recorded yet for this timeframe.
          </div>
        ) : (
          <div className="space-y-6">
            <div className="flex items-end justify-between gap-2 h-48 px-2">
               {stats.daily_timeline.slice(-7).map((d, i) => (
                 <div key={i} className="flex-1 flex flex-col items-center gap-3 group">
                    <div className="relative w-full flex flex-col items-center">
                       <div className="absolute -top-8 opacity-0 group-hover:opacity-100 transition-opacity bg-indigo-500 text-white text-[10px] font-bold px-2 py-1 rounded shadow-xl pointer-events-none">
                         {d.total} clicks
                       </div>
                       <div 
                         className="w-full max-w-[40px] bg-indigo-500/20 group-hover:bg-indigo-500/40 rounded-t-lg transition-all duration-500"
                         style={{ height: `${(d.total / maxClicks) * 100}%`, minHeight: '4px' }}
                       ></div>
                    </div>
                    <span className="text-[10px] text-slate-500 font-medium rotate-45 sm:rotate-0">
                      {new Date(d.day).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
                    </span>
                 </div>
               ))}
            </div>
          </div>
        )}
      </div>

      {/* Info Card */}
      <div className="glass p-6 rounded-2xl border-slate-800 flex flex-col sm:flex-row items-center gap-6">
         <div className="flex items-center gap-4 flex-1">
            <div className="p-3 bg-slate-800 rounded-xl"><Calendar className="w-5 h-5 text-slate-400" /></div>
            <div>
              <p className="text-xs text-slate-500 uppercase font-bold tracking-wider">Created On</p>
              <p className="font-semibold">{new Date(urlInfo.created_at).toLocaleString()}</p>
            </div>
         </div>
         <div className="h-10 w-[1px] bg-slate-800 hidden sm:block"></div>
         <div className="flex items-center gap-4 flex-1">
            <div className="p-3 bg-slate-800 rounded-xl"><BarChart2 className="w-5 h-5 text-slate-400" /></div>
            <div>
              <p className="text-xs text-slate-500 uppercase font-bold tracking-wider">Short Code</p>
              <p className="font-mono font-bold text-indigo-400">{urlInfo.short_code}</p>
            </div>
         </div>
      </div>
    </div>
  );
}

function Zap({ className }: { className?: string }) {
  return (
    <svg 
      xmlns="http://www.w3.org/2000/svg" 
      width="24" 
      height="24" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      strokeWidth="2" 
      strokeLinecap="round" 
      strokeLinejoin="round" 
      className={className}
    >
      <path d="M4 14.71 14 3v8h6l-10 11.71V13H4z"/>
    </svg>
  );
}
