import Link from 'next/link';
import { Button } from '@/components/ui/Button';

export function Header() {
    return (
        <header className="border-b border-white/10 bg-slate-950/50 backdrop-blur-md sticky top-0 z-50">
            <div className="container mx-auto px-4 h-16 flex items-center justify-between">
                <Link href="/" className="flex items-center gap-2">
                    <div className="text-2xl font-bold bg-gradient-to-r from-emerald-400 to-teal-400 bg-clip-text text-transparent">
                        TaxSmart
                    </div>
                    <span className="text-xs text-slate-400 bg-slate-800 px-2 py-0.5 rounded-full border border-slate-700">
                        NG 2026
                    </span>
                </Link>
                <div className="flex items-center gap-4">
                    <nav className="hidden md:flex items-center gap-6 text-sm font-medium text-slate-300">
                        <Link href="/" className="hover:text-white transition-colors">How it Works</Link>
                        <Link href="/" className="hover:text-white transition-colors">Security</Link>
                    </nav>
                    <div className="flex items-center gap-3">
                        <Button variant="ghost" size="sm">Sign In</Button>
                        <Button size="sm">Get Started</Button>
                    </div>
                </div>
            </div>
        </header>
    );
}
