export type TransactionType = 'credit' | 'debit';

export type Category =
    | 'employment_income'
    | 'freelance_income'
    | 'rental_income'
    | 'investment_income'
    | 'crypto_income'
    | 'interest_income'
    | 'other_income'
    | 'expense'
    | 'rent_expense'
    | 'transfer'
    | 'uncategorized';

export interface Transaction {
    id?: string;
    date: string;
    description: string;
    amount: number;
    type: TransactionType;
    category: Category;
    confidence: number;
    is_manual?: boolean;
}

export interface TaxReport {
    id: string;
    tax_year: number;
    total_income: number;
    employment_income: number;
    freelance_income: number;
    rental_income: number;
    crypto_income: number;
    investment_income: number;
    other_income: number;

    // Reliefs
    rent_relief: number;
    pension_deduction: number;
    nhis_deduction: number;
    nhf_deduction: number;
    total_reliefs: number;

    taxable_income: number;
    pit_amount: number;
    total_tax: number;
    breakdown: any;
}

export interface ReliefInput {
    annual_rent: number;
    pension_contribution: number;
    nhis_contribution: number;
    nhf_contribution: number;
}
