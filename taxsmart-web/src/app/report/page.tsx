'use client';

import { Header } from '@/components/layout/Header';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { endpoints } from '@/lib/api';
import { formatCurrency } from '@/lib/utils';
import { TaxReport } from '@/types';
import { ArrowLeft, Download, FileText, Share2 } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';

export default function ReportPage() {
    const router = useRouter();
    const [report, setReport] = useState<TaxReport | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const storedTx = localStorage.getItem('classified_transactions');
        if (!storedTx) {
            router.push('/upload');
            return;
        }

        const calculate = async () => {
            try {
                const transactions = JSON.parse(storedTx);
                const response = await endpoints.calculateTax({
                    tax_year: 2026,
                    user_id: "00000000-0000-0000-0000-000000000000", // Demo ID
                    transactions,
                    reliefs: {
                        annual_rent: 2000000, // Demo value
                        pension_contribution: 0,
                        nhis_contribution: 0,
                        nhf_contribution: 0
                    }
                });
                setReport(response.data.data);
            } catch (error) {
                console.error('Calculation failed:', error);
            } finally {
                setIsLoading(false);
            }
        };

        calculate();
    }, [router]);

    if (isLoading) {
        return (
            <div className="flex flex-col min-h-screen">
                <Header />
                <div className="flex-1 flex items-center justify-center">
                    <div className="text-center">
                        <div className="animate-spin w-12 h-12 border-4 border-emerald-500 border-t-transparent rounded-full mx-auto mb-4" />
                        <p className="text-slate-400">Calculatin Tax Liability...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (!report) return null;

    const chartData = [
        { name: 'Income Tax (PIT)', value: report.pit_amount, color: '#10b981' }, // Emerald-500
        { name: 'Net Income', value: report.total_income - report.total_tax, color: '#3b82f6' }, // Blue-500
    ];

    return (
        <div className="flex flex-col min-h-screen bg-slate-950">
            <Header />

            <main className="flex-1 container mx-auto px-4 py-8 max-w-5xl">
                <div className="flex items-center justify-between mb-8">
                    <Button variant="ghost" className="gap-2" onClick={() => router.back()}>
                        <ArrowLeft size={16} /> Back
                    </Button>
                    <div className="flex gap-3">
                        <Button variant="outline" className="gap-2">
                            <Share2 size={16} /> Share
                        </Button>
                        <Button className="gap-2">
                            <Download size={16} /> Download PDF
                        </Button>
                    </div>
                </div>

                <div className="grid md:grid-cols-3 gap-6 mb-8">
                    <Card className="bg-gradient-to-br from-emerald-500/20 to-emerald-900/10 border-emerald-500/20">
                        <p className="text-emerald-400 font-medium mb-1">Total Tax Liability</p>
                        <h2 className="text-3xl font-bold text-white">{formatCurrency(report.total_tax)}</h2>
                        <p className="text-xs text-emerald-400/70 mt-2">Personal Income Tax (2026)</p>
                    </Card>

                    <Card>
                        <p className="text-slate-400 font-medium mb-1">Total Taxable Income</p>
                        <h2 className="text-3xl font-bold text-white">{formatCurrency(report.taxable_income)}</h2>
                        <p className="text-xs text-slate-500 mt-2">After {formatCurrency(report.total_reliefs)} in reliefs</p>
                    </Card>

                    <Card>
                        <p className="text-slate-400 font-medium mb-1">Effective Tax Rate</p>
                        <h2 className="text-3xl font-bold text-white">
                            {((report.total_tax / report.total_income) * 100).toFixed(1)}%
                        </h2>
                        <p className="text-xs text-slate-500 mt-2">Of gross income</p>
                    </Card>
                </div>

                <div className="grid md:grid-cols-3 gap-8">
                    <div className="md:col-span-2 space-y-6">
                        <Card>
                            <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
                                <FileText size={20} className="text-emerald-500" />
                                Tax Breakdown
                            </h3>

                            <div className="space-y-4">
                                <div className="flex justify-between py-3 border-b border-white/5">
                                    <span className="text-slate-400">Total Gross Income</span>
                                    <span className="font-medium">{formatCurrency(report.total_income)}</span>
                                </div>

                                <div className="pl-4 space-y-2 text-sm">
                                    {report.employment_income > 0 && (
                                        <div className="flex justify-between text-slate-500">
                                            <span>• Employment</span>
                                            <span>{formatCurrency(report.employment_income)}</span>
                                        </div>
                                    )}
                                    {report.freelance_income > 0 && (
                                        <div className="flex justify-between text-slate-500">
                                            <span>• Freelance/Biz</span>
                                            <span>{formatCurrency(report.freelance_income)}</span>
                                        </div>
                                    )}
                                    {report.crypto_income > 0 && (
                                        <div className="flex justify-between text-slate-500">
                                            <span>• Crypto Gains</span>
                                            <span>{formatCurrency(report.crypto_income)}</span>
                                        </div>
                                    )}
                                </div>

                                <div className="flex justify-between py-3 border-b border-white/5">
                                    <span className="text-slate-400">Consolidated Reliefs</span>
                                    <span className="font-medium text-emerald-400">-{formatCurrency(report.total_reliefs)}</span>
                                </div>

                                <div className="flex justify-between py-3 pt-4">
                                    <span className="text-white font-medium">Net Taxable Income</span>
                                    <span className="font-bold text-white">{formatCurrency(report.taxable_income)}</span>
                                </div>
                            </div>
                        </Card>

                        <div className="bg-slate-900 rounded-xl p-6 border border-slate-800">
                            <h4 className="font-medium mb-4 text-slate-300">Tax Calculation logic (PIT 2026)</h4>
                            <div className="space-y-2 text-sm text-slate-500">
                                {report.breakdown.pit_breakdown.map((bracket: any, idx: number) => (
                                    <div key={idx} className="flex justify-between items-center bg-slate-950/50 p-3 rounded">
                                        <div>
                                            <span className="block text-slate-400">
                                                {bracket.rate === 0 ? 'Exempt (0%)' : `Tier ${(bracket.rate * 100).toFixed(0)}%`}
                                            </span>
                                            <span className="text-xs">
                                                On first {formatCurrency(bracket.bracket_max - bracket.bracket_min)}
                                            </span>
                                        </div>
                                        <span className="font-mono text-slate-300">
                                            {formatCurrency(bracket.tax_amount)}
                                        </span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>

                    <div className="space-y-6">
                        <Card>
                            <h3 className="text-lg font-semibold mb-4">Distribution</h3>
                            <div className="h-64 w-full">
                                <ResponsiveContainer width="100%" height="100%">
                                    <PieChart>
                                        <Pie
                                            data={chartData}
                                            cx="50%"
                                            cy="50%"
                                            innerRadius={60}
                                            outerRadius={80}
                                            paddingAngle={5}
                                            dataKey="value"
                                        >
                                            {chartData.map((entry, index) => (
                                                <Cell key={`cell-${index}`} fill={entry.color} />
                                            ))}
                                        </Pie>
                                        <Tooltip
                                            contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b' }}
                                            itemStyle={{ color: '#f8fafc' }}
                                            formatter={(value: number) => formatCurrency(value)}
                                        />
                                        <Legend />
                                    </PieChart>
                                </ResponsiveContainer>
                            </div>
                        </Card>

                        <Card className="bg-blue-500/5 border-blue-500/10">
                            <h4 className="font-medium text-blue-400 mb-2">Tax Saving Tip</h4>
                            <p className="text-sm text-slate-400">
                                Maximizing your pension contributions can reduce your taxable income significantly.
                                You utilized {((report.pension_deduction / report.total_income) * 100).toFixed(1)}% of your allowed relief.
                            </p>
                        </Card>
                    </div>
                </div>
            </main>
        </div>
    );
}
