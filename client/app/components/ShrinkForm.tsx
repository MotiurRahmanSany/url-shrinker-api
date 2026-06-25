"use client";

import { useState } from "react";
import { apiRequest, BASE_URL } from "@/lib/api";
import { Loader2, Link2, Copy, Check, Sparkles, ChevronDown, ChevronUp, Calendar, MousePointer2 } from "lucide-react";
import { useAuth } from "../context/AuthContext";
import Link from "next/link";

type URLCreateResponse = {
  short_code: string;
};

export default function ShrinkForm() {
  const { user } = useAuth();
  const [url, setUrl] = useState("");
  const [customCode, setCustomCode] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [maxClicks, setMaxClicks] = useState("");
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    setLoading(true);
    setError("");
    setResult(null);

    const body: any = {
      original_url: url,
      custom_short_code: customCode || undefined,
    };

    if (expiresAt) {
      body.expires_at = new Date(expiresAt).toISOString();
    }

    if (maxClicks) {
      body.max_clicks = parseInt(maxClicks, 10);
    }

    const res = await apiRequest<URLCreateResponse>("/urls", {
      method: "POST",
      body: JSON.stringify(body),
    });

    if (res.success && res.data) {
      setResult(res.data.short_code);
      setUrl("");
      setCustomCode("");
      setExpiresAt("");
      setMaxClicks("");
      setShowAdvanced(false);
    } else {
      setError(res.message || "Failed to shrink URL");
    }
    setLoading(false);
  };

  const copyToClipboard = () => {
    if (result) {
      navigator.clipboard.writeText(`${BASE_URL}/${result}`);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="w-full max-w-3xl mx-auto space-y-6">
      <form onSubmit={handleSubmit} className="premium-card space-y-6">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Link2 className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" />
            <input
              type="url"
              required
              disabled={!user || loading}
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder={user ? "Paste your long link here..." : "Please sign in to shrink links..."}
              className="w-full bg-slate-950/50 border border-slate-800 rounded-2xl py-4 pl-12 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-lg disabled:opacity-60"
            />
          </div>
          {user ? (
            <button
              type="submit"
              disabled={loading}
              className="bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white font-bold py-4 px-8 rounded-2xl shadow-lg shadow-indigo-500/20 transition-all active:scale-95 flex items-center justify-center gap-2 min-w-[160px]"
            >
              {loading ? <Loader2 className="w-6 h-6 animate-spin" /> : (
                <>
                  <Sparkles className="w-5 h-5" />
                  <span>Shrink It!</span>
                </>
              )}
            </button>
          ) : (
            <Link
              href="/login"
              className="bg-indigo-600 hover:bg-indigo-500 text-white font-bold py-4 px-8 rounded-2xl shadow-lg shadow-indigo-500/20 transition-all active:scale-95 flex items-center justify-center gap-2 min-w-[160px] text-center"
            >
              Get Started
            </Link>
          )}
        </div>

        {user && (
          <>
            <div className="flex flex-col gap-3">
              <button
                type="button"
                onClick={() => setShowAdvanced(!showAdvanced)}
                className="flex items-center gap-2 text-xs font-semibold text-slate-400 hover:text-indigo-400 uppercase tracking-wider ml-1 w-fit transition-colors"
              >
                <span>Advanced Options</span>
                {showAdvanced ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
              </button>
              <div className="border-t border-slate-800/50 w-full"></div>
            </div>

            {showAdvanced && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6 animate-in fade-in slide-in-from-top-2 duration-300">
                <div className="space-y-2">
                  <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider ml-1">Custom Code (Optional)</label>
                  <div className="relative">
                    <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500 text-sm font-mono">/</span>
                    <input
                      type="text"
                      value={customCode}
                      onChange={(e) => setCustomCode(e.target.value)}
                      placeholder="my-cool-link"
                      className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 pl-8 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm font-mono"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider ml-1">Expiration Date (Optional)</label>
                  <div className="relative">
                    <Calendar className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500 pointer-events-none" />
                    <input
                      type="datetime-local"
                      value={expiresAt}
                      onChange={(e) => setExpiresAt(e.target.value)}
                      className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 pl-10 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm [color-scheme:dark]"
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <label className="text-xs font-semibold text-slate-400 uppercase tracking-wider ml-1">Max Clicks (Optional)</label>
                  <div className="relative">
                    <MousePointer2 className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" />
                    <input
                      type="number"
                      min="1"
                      value={maxClicks}
                      onChange={(e) => setMaxClicks(e.target.value)}
                      placeholder="e.g. 100"
                      className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-3 pl-10 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm"
                    />
                  </div>
                </div>
              </div>
            )}
          </>
        )}

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-xl text-sm animate-in fade-in slide-in-from-top-1">
            {error}
          </div>
        )}
      </form>

      {result && (
        <div className="glass rounded-2xl p-6 border border-indigo-500/50 animate-in fade-in zoom-in duration-500">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="space-y-1 text-center sm:text-left">
              <p className="text-xs font-semibold text-slate-400 uppercase tracking-widest">Your short link is ready!</p>
              <p className="text-xl font-bold text-indigo-400 break-all">{`${BASE_URL}/${result}`}</p>
            </div>
            <button
              onClick={copyToClipboard}
              className="flex items-center gap-2 bg-slate-800 hover:bg-slate-700 text-white px-6 py-3 rounded-xl transition-all active:scale-95 group shrink-0"
            >
              {copied ? (
                <>
                  <Check className="w-5 h-5 text-emerald-400" />
                  <span className="text-emerald-400">Copied!</span>
                </>
              ) : (
                <>
                  <Copy className="w-5 h-5 group-hover:text-indigo-400 transition-colors" />
                  <span>Copy Link</span>
                </>
              )}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
