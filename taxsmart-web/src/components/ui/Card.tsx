import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface CardProps {
    children: ReactNode;
    className?: string;
}

export function Card({ children, className }: CardProps) {
    return (
        <div className={cn("bg-slate-900/50 border border-white/5 rounded-xl p-6 backdrop-blur-sm", className)}>
            {children}
        </div>
    );
}
