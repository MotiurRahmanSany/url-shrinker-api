"use client";

import { useState } from "react";
import { apiRequest, BASE_URL } from "@/lib/api";
import { Loader2, Link2, Copy, Check, Sparkles } from "lucide-react";

export default function ShrinkForm() {
  const [url, setUrl] = useState("");
  const [customCode, setCustomCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setResult(null);

    const res = await apiRequest<{ short_url: string }>("/urls", {
      method: "POST",
      body: JSON.stringify({
        original_url: url,
        custom_short_code: customCode || undefined,
      }),
    });

    if (res.success && res.data) {
      setResult(res.data.short_url);
      setUrl("");
      setCustomCode("");
    } else {
      setError(res.message || "Failed to shrink URL");
    }
    setLoading(false);
  };

  const copyToClipboard = () => {
    if (result) {
      navigator.clipboard.writeText(result);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="w-full max-w-3xl mx-auto space-y-6">
      <form onSubmit={handleSubmit} className="premium-card space-y-4">
        <div className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Link2 className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-500" />
            <input
              type="url"
              required
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Paste your long link here..."
              className="w-full bg-slate-950/50 border border-slate-800 rounded-2xl py-4 pl-12 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-lg"
            />
          </div>
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
        </div>

        <div className="flex items-center gap-2">
          <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider ml-1">Custom Code (Optional)</span>
          <div className="flex-1 border-t border-slate-800/50"></div>
        </div>

        <div className="max-w-xs relative">
          <span className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500 text-sm font-mono">/</span>
          <input
            type="text"
            value={customCode}
            onChange={(e) => setCustomCode(e.target.value)}
            placeholder="my-cool-link"
            className="w-full bg-slate-950/50 border border-slate-800 rounded-xl py-2 pl-8 pr-4 focus:ring-2 focus:ring-indigo-500/50 focus:border-indigo-500 outline-none transition-all text-sm font-mono"
          />
        </div>

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 text-red-400 p-3 rounded-xl text-sm animate-in fade-in slide-in-from-top-1">
            {error}
          </div>
        )}
      </form>

      {result && (
        <div className="glass rounded-2xl p-6 border-indigo-500/50 animate-in fade-in zoom-in duration-500">
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
