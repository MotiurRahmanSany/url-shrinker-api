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
  Link2,
  Pencil,
  X
} from "lucide-react";
import Link from "next/link";

type URLItem = {
  id: number;
  original_url: string;
  short_code: string;
  is_active: boolean;
  expires_at?: string;
  max_clicks?: number;
  created_at: string;
  total_clicks: number;
};

type PaginationMeta = {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
};

export default function Dashboard() {
  const { user, loading: authLoading } = useAuth();
  const [urls, setUrls] = useState<URLItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [meta, setMeta] = useState<PaginationMeta>({ page: 1, limit: 10, total: 0, total_pages: 0 });
  const router = useRouter();

  // Edit URL states
  const [editingUrl, setEditingUrl] = useState<URLItem | null>(null);
  const [editOriginalUrl, setEditOriginalUrl] = useState("");
  const [editExpiresAt, setEditExpiresAt] = useState("");
  const [editMaxClicks, setEditMaxClicks] = useState("");
  const [editLoading, setEditLoading] = useState(false);
  const [editError, setEditError] = useState("");

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]);

  const fetchUrls = async (page = 1) => {
    setLoading(true);
    setError("");
    const res = await apiRequest<{ items: URLItem[]; meta: PaginationMeta }>(`/urls?page=${page}&limit=10`);
    if (res.success && res.data) {
      const items = res.data.items || [];
      const initializedItems = items.map(item => ({ ...item, total_clicks: 0 }));
      setUrls(initializedItems);
      
      if (res.data.meta) {
        setMeta(res.data.meta);
      }
      
      // Async fetch clicks stats for each URL code
      items.forEach(async (urlItem) => {
        const statsRes = await apiRequest<{ total_clicks: number }>(`/urls/${urlItem.short_code}/stats`);
        if (statsRes.success && statsRes.data) {
          const clicks = statsRes.data.total_clicks;
          setUrls(prev => prev.map(u => u.short_code === urlItem.short_code ? { ...u, total_clicks: clicks } : u));
        }
      });
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
    if (!confirm("Are you sure you want to deactivate/delete this URL?")) return;
    
    const res = await apiRequest(`/urls/${code}`, {
      method: "DELETE",
    });

    if (res.success) {
      setUrls(prev => prev.map(u => u.short_code === code ? { ...u, is_active: false } : u));
    }
  };

  const handleEditClick = (urlItem: URLItem) => {
    setEditingUrl(urlItem);
    setEditOriginalUrl(urlItem.original_url);
    if (urlItem.expires_at) {
      const d = new Date(urlItem.expires_at);
      const offset = d.getTimezoneOffset();
      const localDate = new Date(d.getTime() - (offset * 60 * 1000));
      setEditExpiresAt(localDate.toISOString().slice(0, 16));
    } else {
      setEditExpiresAt("");
    }
    setEditMaxClicks(urlItem.max_clicks ? urlItem.max_clicks.toString() : "");
    setEditError("");
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingUrl) return;
    setEditLoading(true);
    setEditError("");

    const body: any = {
      original_url: editOriginalUrl,
    };

    body.expires_at = editExpiresAt ? new Date(editExpiresAt).toISOString() : "0001-01-01T00:00:00Z";
    body.max_clicks = editMaxClicks ? parseInt(editMaxClicks, 10) : 0;

    const res = await apiRequest<URLItem>(`/urls/${editingUrl.short_code}`, {
      method: "PATCH",
      body: JSON.stringify(body),
    });

    if (res.success && res.data) {
      const updated = res.data;
      setUrls(urls.map(u => u.short_code === editingUrl.short_code ? { 
        ...u, 
        original_url: updated.original_url, 
        expires_at: updated.expires_at, 
        max_clicks: updated.max_clicks,
        is_active: updated.is_active
      } : u));
      setEditingUrl(null);
    } else {
      setEditError(res.message || "Failed to update URL");
    }
    setEditLoading(false);
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
          value={meta.total.toString()} 
        />
        <StatSummaryCard 
          icon={<MousePointer2 className="w-5 h-5 text-cyan-400" />} 
          label="Total Clicks" 
          value={urls.reduce((acc, curr) => acc + (curr.total_clicks || 0), 0).toString()} 
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
      <div className="premium-card !p-0 overflow-hidden space-y-4">
        {urls.length === 0 ? (
          <div className="p-12 text-center space-y-4">
            <div className="bg-slate-800/50 w-16 h-16 rounded-2xl flex items-center justify-center mx-auto">
              <Link2 className="w-8 h-8 text-slate-500" />
            </div>
            <p className="text-slate-400">You haven't shortened any URLs yet.</p>
          </div>
        ) : (
          <div className="flex flex-col">
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
                        <p className="text-sm font-medium truncate text-slate-300" title={url.original_url}>{url.original_url}</p>
                        <div className="flex flex-wrap gap-x-3 gap-y-1 mt-1 text-[10px] text-slate-500">
                          <span>Created: {new Date(url.created_at).toLocaleDateString()}</span>
                          {url.expires_at && (
                            <span className="text-indigo-400">
                              Expires: {new Date(url.expires_at).toLocaleDateString()}
                            </span>
                          )}
                          {url.max_clicks && (
                            <span className="text-cyan-400">
                              Limit: {url.max_clicks} clicks
                            </span>
                          )}
                        </div>
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
                          <span className="text-sm font-bold">{url.total_clicks || 0}</span>
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
                            onClick={() => handleEditClick(url)}
                            className="p-2 bg-slate-800 hover:bg-indigo-500/20 hover:text-indigo-400 rounded-lg transition-all"
                            title="Edit URL Settings"
                          >
                            <Pencil className="w-4 h-4" />
                          </button>
                          <button 
                            onClick={() => deleteUrl(url.short_code)}
                            className="p-2 bg-slate-800 hover:bg-red-500/20 hover:text-red-400 rounded-lg transition-all"
                            title="Deactivate URL"
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

            {/* Pagination Controls */}
            {meta.total_pages > 1 && (
              <div className="flex items-center justify-between p-6 border-t border-slate-800/50">
                <div className="text-xs text-slate-500">
                  Showing Page <span className="font-bold text-slate-300">{meta.page}</span> of <span className="font-bold text-slate-300">{meta.total_pages}</span> ({meta.total} total links)
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => fetchUrls(meta.page - 1)}
                    disabled={meta.page <= 1}
                    className="bg-slate-800 hover:bg-slate-700 disabled:opacity-40 text-white px-4 py-2 rounded-xl text-xs font-semibold transition-all disabled:pointer-events-none"
                  >
                    Previous
                  </button>
                  <button
                    onClick={() => fetchUrls(meta.page + 1)}
                    disabled={meta.page >= meta.total_pages}
                    className="bg-slate-800 hover:bg-slate-700 disabled:opacity-40 text-white px-4 py-2 rounded-xl text-xs font-semibold transition-all disabled:pointer-events-none"
                  >
                    Next
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Edit URL Modal */}
      {editingUrl && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="glass max-w-lg w-full rounded-3xl p-8 border border-slate-800 shadow-2xl relative space-y-6 animate-in zoom-in duration-300">
            <button 
              onClick={() => setEditingUrl(null)} 
              className="absolute right-6 top-6 p-2 bg-slate-800/50 hover:bg-slate-800 rounded-full text-slate-400 hover:text-white transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
            <div className="space-y-1">
              <h2 className="text-2xl font-bold">Edit URL Settings</h2>
              <p className="text-sm text-slate-400">Short Code: <span className="font-mono text-indigo-400 font-bold">/{editingUrl.short_code}</span></p>
            </div>
            {editError && (
              <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-xl text-sm">
                {editError}
              </div>
            )}
            <form onSubmit={handleEditSubmit} className="space-y-6">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-300 ml-1">Original URL</label>
                <div className="relative">
                  <Link2 className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
                  <input
                    type="url"
                    required
                    value={editOriginalUrl}
                    onChange={(e) => setEditOriginalUrl(e.target.value)}
                    className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 pl-10 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm"
                  />
                </div>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-300 ml-1">Expires At (Optional)</label>
                  <input
                    type="datetime-local"
                    value={editExpiresAt}
                    onChange={(e) => setEditExpiresAt(e.target.value)}
                    className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 px-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm [color-scheme:dark]"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-300 ml-1">Max Clicks (Optional)</label>
                  <input
                    type="number"
                    min="1"
                    value={editMaxClicks}
                    onChange={(e) => setEditMaxClicks(e.target.value)}
                    placeholder="Unlimited"
                    className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 px-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm"
                  />
                </div>
              </div>
              <div className="flex gap-4 pt-4 border-t border-slate-800/50">
                <button
                  type="button"
                  onClick={() => setEditingUrl(null)}
                  className="flex-1 bg-slate-800 hover:bg-slate-700 text-white font-bold py-3 px-6 rounded-xl transition-all"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={editLoading}
                  className="flex-1 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white font-bold py-3 px-6 rounded-xl transition-all flex items-center justify-center gap-2"
                >
                  {editLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : "Save Changes"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
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
