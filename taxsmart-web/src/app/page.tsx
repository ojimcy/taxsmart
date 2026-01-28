'use client';

import { Header } from '@/components/layout/Header';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { ArrowRight, CheckCircle2, ShieldCheck, Zap } from 'lucide-react';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen">
      <Header />

      <main className="flex-1">
        {/* Hero Section */}
        <section className="relative py-20 md:py-32 px-4 overflow-hidden">
          <div className="absolute inset-0 bg-gradient-to-b from-emerald-500/10 to-transparent pointer-events-none" />
          <div className="container mx-auto max-w-6xl relative z-10 text-center">
            <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6 bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">
              Smart Tax Calculation <br /> for Nigerians
            </h1>
            <p className="text-xl text-slate-400 mb-10 max-w-2xl mx-auto leading-relaxed">
              Upload your bank statements or crypto transaction history.
              Our AI automatically classifies income, calculates tax liability,
              and applies all legal reliefs for the 2026 tax year.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/upload">
                <Button size="lg" className="w-full sm:w-auto gap-2">
                  Calculate My Tax <ArrowRight size={18} />
                </Button>
              </Link>
              <Button variant="outline" size="lg" className="w-full sm:w-auto">
                View Demo Report
              </Button>
            </div>
          </div>
        </section>

        {/* Features Grid */}
        <section className="py-20 bg-slate-900/50 border-y border-white/5">
          <div className="container mx-auto px-4 max-w-6xl">
            <div className="grid md:grid-cols-3 gap-8">
              <Card>
                <div className="w-12 h-12 bg-emerald-500/10 rounded-lg flex items-center justify-center text-emerald-400 mb-4">
                  <Zap size={24} />
                </div>
                <h3 className="text-xl font-semibold mb-3">AI Classification</h3>
                <p className="text-slate-400">
                  Automatically detects salary, freelance income, crypto gains, and business expenses from generic bank narrations.
                </p>
              </Card>
              <Card>
                <div className="w-12 h-12 bg-blue-500/10 rounded-lg flex items-center justify-center text-blue-400 mb-4">
                  <ShieldCheck size={24} />
                </div>
                <h3 className="text-xl font-semibold mb-3">2026 Tax Compliant</h3>
                <p className="text-slate-400">
                  Updated with the latest Nigeria Tax Act 2025 rules, including new progressive brackets and rent relief.
                </p>
              </Card>
              <Card>
                <div className="w-12 h-12 bg-purple-500/10 rounded-lg flex items-center justify-center text-purple-400 mb-4">
                  <CheckCircle2 size={24} />
                </div>
                <h3 className="text-xl font-semibold mb-3">Detailed Reports</h3>
                <p className="text-slate-400">
                  Get a comprehensive breakdown of your tax liability, reliefs applied, and category-wise income distribution.
                </p>
              </Card>
            </div>
          </div>
        </section>
      </main>

      <footer className="py-8 border-t border-white/5 text-center text-slate-500 text-sm">
        <p>&copy; {new Date().getFullYear()} TaxSmart Nigeria. All rights reserved.</p>
      </footer>
    </div>
  );
}
