'use client';

import { Header } from '@/components/layout/Header';
import { FileUploader } from '@/components/upload/FileUploader';
import { endpoints } from '@/lib/api';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function UploadPage() {
    const router = useRouter();
    const [isUploading, setIsUploading] = useState(false);

    const handleUpload = async (file: File) => {
        setIsUploading(true);
        try {
            const response = await endpoints.parse(file);
            // In a real app, we'd store the transaction IDs or upload ID to fetch in the review page
            // For now, let's assume we pass data via state or query params (simplified)

            // Store transactions in localStorage for demo purposes
            // In production, this would be fetched from the backend via UploadID
            localStorage.setItem('current_upload', JSON.stringify(response.data.data));

            router.push('/review');
        } catch (error) {
            console.error('Upload failed:', error);
            alert('Failed to process file. Please try again.');
        } finally {
            setIsUploading(false);
        }
    };

    return (
        <div className="flex flex-col min-h-screen">
            <Header />

            <main className="flex-1 container mx-auto px-4 py-12 max-w-4xl">
                <div className="text-center mb-12">
                    <h1 className="text-3xl font-bold mb-4">Upload Your Statement</h1>
                    <p className="text-slate-400 max-w-2xl mx-auto">
                        We support PDF and CSV statements from major Nigerian banks (GTBank, Access, Zenith, UBA) and crypto exchanges.
                    </p>
                </div>

                <FileUploader onUpload={handleUpload} isUploading={isUploading} />

                <div className="mt-16 border-t border-white/5 pt-8">
                    <h3 className="text-sm font-medium text-slate-500 mb-4 uppercase tracking-wider">Supported Formats</h3>
                    <div className="flex flex-wrap gap-4 text-sm text-slate-400">
                        <span className="px-3 py-1 bg-slate-900 rounded-full border border-slate-800">GTBank CSV/PDF</span>
                        <span className="px-3 py-1 bg-slate-900 rounded-full border border-slate-800">Access Bank CSV</span>
                        <span className="px-3 py-1 bg-slate-900 rounded-full border border-slate-800">Binance Export</span>
                        <span className="px-3 py-1 bg-slate-900 rounded-full border border-slate-800">Cowrywise CSV</span>
                    </div>
                </div>
            </main>
        </div>
    );
}
