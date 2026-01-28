import axios from 'axios';
import { auth } from './firebase';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add auth token to requests
// Add auth token to requests
api.interceptors.request.use(async (config) => {
    const user = auth.currentUser;
    if (user) {
        const token = await user.getIdToken();
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const endpoints = {
    parse: (file: File) => {
        const formData = new FormData();
        formData.append('file', file);
        return api.post('/parse', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
    },

    classify: (transactions: any[]) => api.post('/classify', transactions),

    calculateTax: (data: any) => api.post('/tax/calculate', data),

    quickPit: (annualIncome: number) => api.post('/tax/quick-pit', { annual_income: annualIncome }),
};
