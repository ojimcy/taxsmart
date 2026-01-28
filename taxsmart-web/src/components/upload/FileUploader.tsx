'use client';

import { useCallback, useState } from 'react';
import { useDropzone, FileRejection } from 'react-dropzone';
import { UploadCloud, FileText, X, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/Button';

interface FileUploaderProps {
    onUpload: (file: File) => void;
    isUploading?: boolean;
}

export function FileUploader({ onUpload, isUploading }: FileUploaderProps) {
    const [file, setFile] = useState<File | null>(null);
    const [error, setError] = useState<string | null>(null);

    const onDrop = useCallback((acceptedFiles: File[], fileRejections: FileRejection[]) => {
        setError(null);

        if (fileRejections.length > 0) {
            setError('Please upload a valid CSV or PDF file (max 10MB)');
            return;
        }

        if (acceptedFiles.length > 0) {
            setFile(acceptedFiles[0]);
        }
    }, []);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: {
            'text/csv': ['.csv'],
            'application/pdf': ['.pdf'],
        },
        maxSize: 10 * 1024 * 1024, // 10MB
        multiple: false,
        disabled: isUploading,
    });

    const handleUpload = () => {
        if (file) {
            onUpload(file);
        }
    };

    const removeFile = () => {
        setFile(null);
        setError(null);
    };

    return (
        <div className="w-full max-w-xl mx-auto">
            {!file ? (
                <div
                    {...getRootProps()}
                    className={cn(
                        "border-2 border-dashed rounded-xl p-10 text-center cursor-pointer transition-all duration-200 group",
                        isDragActive
                            ? "border-emerald-500 bg-emerald-500/10"
                            : "border-slate-700 hover:border-emerald-500/50 hover:bg-slate-800/50"
                    )}
                >
                    <input {...getInputProps()} />
                    <div className="flex flex-col items-center gap-4">
                        <div className={cn(
                            "p-4 rounded-full bg-slate-800 transition-colors",
                            isDragActive ? "bg-emerald-500/20 text-emerald-400" : "text-slate-400 group-hover:text-emerald-400"
                        )}>
                            <UploadCloud size={32} />
                        </div>
                        <div>
                            <p className="text-lg font-medium text-slate-200">
                                {isDragActive ? "Drop your statement here" : "Click to upload or drag and drop"}
                            </p>
                            <p className="text-sm text-slate-400 mt-1">
                                Supports CSV and PDF bank statements (max 10MB)
                            </p>
                        </div>
                    </div>
                </div>
            ) : (
                <div className="bg-slate-900 border border-slate-700 rounded-xl p-6">
                    <div className="flex items-center justify-between mb-6">
                        <div className="flex items-center gap-4">
                            <div className="p-3 bg-emerald-500/10 text-emerald-400 rounded-lg">
                                <FileText size={24} />
                            </div>
                            <div>
                                <p className="font-medium text-slate-200">{file.name}</p>
                                <p className="text-sm text-slate-400">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                            </div>
                        </div>
                        <button
                            onClick={removeFile}
                            disabled={isUploading}
                            className="text-slate-400 hover:text-white transition-colors"
                        >
                            <X size={20} />
                        </button>
                    </div>

                    <Button
                        className="w-full"
                        onClick={handleUpload}
                        isLoading={isUploading}
                    >
                        {isUploading ? 'Processing Statement...' : 'Analyze Statement'}
                    </Button>
                </div>
            )}

            {error && (
                <div className="mt-4 p-4 bg-red-500/10 border border-red-500/20 rounded-lg flex items-center gap-3 text-red-400">
                    <AlertCircle size={20} />
                    <p className="text-sm">{error}</p>
                </div>
            )}
        </div>
    );
}
