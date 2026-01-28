'use client';

import { Header } from '@/components/layout/Header';
import { Button } from '@/components/ui/Button';
import { endpoints } from '@/lib/api';
import { cn, formatCurrency, formatDate } from '@/lib/utils';
import { Transaction } from '@/types';
import { ArrowRight, Check, AlertCircle, Loader2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function ReviewPage() {
  const router = useRouter();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [classifiedTransactions, setClassifiedTransactions] = useState<Transaction[]>([]);
  const [isClassifying, setIsClassifying] = useState(false);
  const [classificationComplete, setClassificationComplete] = useState(false);

  useEffect(() => {
    // In a real app, fetch from API. Here we load from localStorage for demo.
    const stored = localStorage.getItem('current_upload');
    if (stored) {
      setTransactions(JSON.parse(stored));
    } else {
      router.push('/upload');
    }
  }, [router]);

  const runClassification = async () => {
    setIsClassifying(true);
    try {
      const response = await endpoints.classify(transactions);
      
      // Transform response to match Transaction type
      const classified = response.data.data.transactions.map((tx: any) => ({
        ...tx,
        is_manual: tx.confidence < 0.7 // Flag low confidence for review
      }));
      
      setClassifiedTransactions(classified);
      setClassificationComplete(true);
    } catch (error) {
      console.error('Classification failed:', error);
      alert('Failed to classify transactions.');
    } finally {
      setIsClassifying(false);
    }
  };

  const proceedToTax = () => {
    localStorage.setItem('classified_transactions', JSON.stringify(classifiedTransactions));
    router.push('/report');
  };

  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      
      <main className="flex-1 container mx-auto px-4 py-8 max-w-5xl">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-2xl font-bold">Review Transactions</h1>
            <p className="text-slate-400">
              {transactions.length} transactions uploaded
            </p>
          </div>
          
          <div className="flex gap-3">
            {!classificationComplete ? (
              <Button onClick={runClassification} isLoading={isClassifying}>
                {isClassifying ? 'AI Processing...' : 'Run AI Classification'}
              </Button>
            ) : (
              <Button onClick={proceedToTax} className="gap-2">
                Calculate Tax <ArrowRight size={16} />
              </Button>
            )}
          </div>
        </div>

        {classificationComplete ? (
          <div className="bg-slate-900 border border-slate-800 rounded-xl overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-left">
                <thead className="bg-slate-950 text-slate-400 uppercase text-xs">
                  <tr>
                    <th className="px-6 py-3">Date</th>
                    <th className="px-6 py-3">Description</th>
                    <th className="px-6 py-3 text-right">Amount</th>
                    <th className="px-6 py-3">Category</th>
                    <th className="px-6 py-3 text-center">Ai Conf.</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800">
                  {classifiedTransactions.map((tx, idx) => (
                    <tr key={idx} className="hover:bg-slate-800/50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap text-slate-400">
                        {formatDate(tx.date)}
                      </td>
                      <td className="px-6 py-4 text-white font-medium">
                        {tx.description}
                      </td>
                      <td className={cn(
                        "px-6 py-4 text-right font-mono",
                        tx.type === 'credit' ? "text-emerald-400" : "text-white"
                      )}>
                        {tx.type === 'credit' ? '+' : '-'}{formatCurrency(tx.amount)}
                      </td>
                      <td className="px-6 py-4">
                        <span className={cn(
                          "px-2 py-1 rounded-full text-xs border uppercase",
                          getCategoryStyle(tx.category)
                        )}>
                          {tx.category.replace('_', ' ')}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-center">
                        <div className="flex items-center justify-center gap-1" title={`${(tx.confidence * 100).toFixed(0)}%`}>
                          {tx.confidence > 0.8 ? (
                            <Check size={16} className="text-emerald-500" />
                          ) : (
                            <AlertCircle size={16} className="text-amber-500" />
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        ) : (
          <div className="text-center py-20 bg-slate-900/50 rounded-xl border border-dashed border-slate-800">
            <div className="mx-auto w-16 h-16 bg-slate-800 rounded-full flex items-center justify-center mb-4">
              <Loader2 size={32} className="text-slate-500 animate-spin" />
            </div>
            <h3 className="text-lg font-medium">Ready to Classify</h3>
            <p className="text-slate-400 max-w-md mx-auto mt-2">
              Our AI engine will analyze {transactions.length} transactions to identify taxable income, deductible expenses, and transfers.
            </p>
          </div>
        )}
      </main>
    </div>
  );
}

function getCategoryStyle(category: string) {
  if (category.includes('income')) return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
  if (category.includes('expense')) return 'bg-red-500/10 text-red-400 border-red-500/20';
  if (category.includes('transfer')) return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
  return 'bg-slate-800 text-slate-400 border-slate-700';
}
