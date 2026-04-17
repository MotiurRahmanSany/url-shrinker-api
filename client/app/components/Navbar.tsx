"use client";

import Link from "next/link";
import { useAuth } from "../context/AuthContext";
import { Link2 } from "lucide-react";

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <nav className="glass sticky top-0 z-50 py-4 px-8 flex items-center justify-between border-b border-slate-800/50">
      <Link href="/" className="flex items-center gap-2">
        <div className="bg-indigo-500 p-2 rounded-lg">
          <Link2 className="w-5 h-5 text-white" />
        </div>
        <span className="text-xl font-bold tracking-tight text-gradient">URL Shrinker</span>
      </Link>

      <div className="flex items-center gap-6">
        {user ? (
          <>
            <Link href="/dashboard" className="text-sm font-medium hover:text-indigo-400 transition-colors">
              Dashboard
            </Link>
            <div className="flex items-center gap-4 border-l border-slate-800 pl-6 ml-2">
              <span className="text-sm text-slate-400">{user.email}</span>
              <button
                onClick={logout}
                className="text-sm font-medium bg-slate-800 hover:bg-slate-700 px-4 py-2 rounded-full transition-all"
              >
                Logout
              </button>
            </div>
          </>
        ) : (
          <>
            <Link href="/login" className="text-sm font-medium hover:text-indigo-400 transition-colors">
              Login
            </Link>
            <Link
              href="/register"
              className="bg-indigo-600 hover:bg-indigo-500 text-white px-5 py-2.5 rounded-full text-sm font-semibold transition-all shadow-lg shadow-indigo-500/20 active:scale-95"
            >
              Sign Up
            </Link>
          </>
        )}
      </div>
    </nav>
  );
}
