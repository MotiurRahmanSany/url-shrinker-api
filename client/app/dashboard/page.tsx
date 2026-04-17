"use client";

import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useRouter } from "next/navigation";
import { apiRequest, BASE_URL } from "@/lib/api";
import { 
  Loader2, 
  Plus, 
  ExternalLink, 
  BarChart2, 
  Trash2, 
  Calendar, 
  MousePointer2,
  AlertCircle,
  Link2
} from "lucide-react";
import Link from "next/link";

type URLItem = {
  id: string;
  original_url: string;
  short_code: string;
  short_url: string;
  total_clicks: number;
  created_at: string;
  is_active: boolean;
};

export default function Dashboard() {
  const { user, loading: authLoading } = useAuth();
  const [urls, setUrls] = useState<URLItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const router = useRouter();

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]);

  const fetchUrls = async () => {
    setLoading(true);
    const res = await apiRequest<{ items: URLItem[] }>("/urls");
    if (res.success && res.data && Array.isArray(res.data.items)) {
      setUrls(res.data.items);
    } else if (res.success && !res.data) {
      setUrls([]);
    } else {
      setError(res.message || "Failed to load URLs");
    }
    setLoading(false);
  };

  useEffect(() => {
    if (user) {
      fetchUrls();
    }
  }, [user]);

  const deleteUrl = async (code: string) => {
    if (!confirm("Are you sure you want to delete this URL?")) return;
    
    const res = await apiRequest(`/urls/${code}`, {
      method: "DELETE",
    });

    if (res.success) {
      setUrls(urls.filter(u => u.short_code !== code));
    }
  };

  if (authLoading || loading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <Loader2 className="w-10 h-10 animate-spin text-indigo-500" />
      </div>
    );
  }

  return (
    <div className="flex-1 p-8 max-w-7xl mx-auto w-full space-y-8">
      <header className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold">My Dashboard</h1>
          <p className="text-slate-400">Manage your shortened links and track performance</p>
        </div>
        <Link 
          href="/" 
          className="bg-indigo-600 hover:bg-indigo-500 text-white font-bold py-3 px-6 rounded-xl flex items-center gap-2 transition-all shadow-lg shadow-indigo-500/20 active:scale-95 w-fit"
        >
          <Plus className="w-5 h-5" />
          <span>New Link</span>
        </Link>
      </header>

      {/* Stats Summary Panel */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <StatSummaryCard 
          icon={<Link2 className="w-5 h-5 text-indigo-400" />} 
          label="Total Links" 
          value={urls.length.toString()} 
        />
        <StatSummaryCard 
          icon={<MousePointer2 className="w-5 h-5 text-cyan-400" />} 
          label="Total Clicks" 
          value={urls.reduce((acc, curr) => acc + curr.total_clicks, 0).toString()} 
        />
        <StatSummaryCard 
          icon={<Calendar className="w-5 h-5 text-emerald-400" />} 
          label="Active Links" 
          value={urls.filter(u => u.is_active).length.toString()} 
        />
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/20 text-red-500 p-4 rounded-xl flex items-center gap-3">
          <AlertCircle className="w-5 h-5" />
          <p className="text-sm">{error}</p>
        </div>
      )}

      {/* URL Table Area */}
      <div className="premium-card !p-0 overflow-hidden">
        {urls.length === 0 ? (
          <div className="p-12 text-center space-y-4">
            <div className="bg-slate-800/50 w-16 h-16 rounded-2xl flex items-center justify-center mx-auto">
              <Link2 className="w-8 h-8 text-slate-500" />
            </div>
            <p className="text-slate-400">You haven't shortened any URLs yet.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="bg-slate-800/30 border-b border-slate-800 transition-colors">
                  <th className="px-6 py-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">Original URL</th>
                  <th className="px-6 py-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">Short Link</th>
                  <th className="px-6 py-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">Clicks</th>
                  <th className="px-6 py-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-4 text-xs font-semibold text-slate-400 uppercase tracking-wider text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/50">
                {urls.map((url) => (
                  <tr key={url.id} className="hover:bg-slate-800/20 transition-colors group">
                    <td className="px-6 py-4 max-w-xs md:max-w-sm">
                      <p className="text-sm font-medium truncate text-slate-300">{url.original_url}</p>
                      <p className="text-[10px] text-slate-500 mt-1">{new Date(url.created_at).toLocaleDateString()}</p>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2 group/link">
                        <span className="text-sm font-bold text-indigo-400">{`${BASE_URL}/${url.short_code}`}</span>
                        <a 
                          href={`${BASE_URL}/${url.short_code}`} 
                          target="_blank" 
                          className="opacity-0 group-hover/link:opacity-100 transition-opacity p-1 text-slate-400 hover:text-white"
                        >
                          <ExternalLink className="w-3.5 h-3.5" />
                        </a>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-2">
                        <BarChart2 className="w-4 h-4 text-cyan-500" />
                        <span className="text-sm font-bold">{url.total_clicks}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`inline-flex px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider ${
                        url.is_active ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20" : "bg-red-500/10 text-red-400 border border-red-500/20"
                      }`}>
                        {url.is_active ? "Active" : "Inactive"}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Link 
                          href={`/stats/${url.short_code}`}
                          className="p-2 bg-slate-800 hover:bg-indigo-500/20 hover:text-indigo-400 rounded-lg transition-all"
                          title="View Statistics"
                        >
                          <BarChart2 className="w-4 h-4" />
                        </Link>
                        <button 
                          onClick={() => deleteUrl(url.short_code)}
                          className="p-2 bg-slate-800 hover:bg-red-500/20 hover:text-red-400 rounded-lg transition-all"
                          title="Delete URL"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

function StatSummaryCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="glass p-6 rounded-2xl flex items-center gap-4">
      <div className="p-3 bg-slate-800/50 rounded-xl">{icon}</div>
      <div>
        <p className="text-xs font-semibold text-slate-400 uppercase tracking-widest">{label}</p>
        <p className="text-2xl font-bold">{value}</p>
      </div>
    </div>
  );
}
