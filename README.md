# TaxSmart Nigeria

AI-powered tax calculator for Nigerian bank statements and crypto transactions.

## Project Structure

```
taxsmart-dev/
├── taxsmart-api/     # Go backend (parsing, AI classification, tax calculations)
└── taxsmart-web/     # Next.js frontend
```

## Tech Stack

- **Backend**: Go + Chi router
- **Database**: Supabase (PostgreSQL)
- **Auth**: Supabase Auth
- **Storage**: Supabase Storage
- **Frontend**: Next.js 14 + TypeScript + Tailwind CSS
- **AI**: Gemini/OpenAI/Claude (configurable)

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Supabase account

### Backend Setup

```bash
cd taxsmart-api
cp .env.example .env
# Edit .env with your Supabase credentials
go run cmd/server/main.go
```

### Frontend Setup

```bash
cd taxsmart-web
npm install
cp .env.example .env.local
# Edit .env.local with your Supabase credentials
npm run dev
```

## Features

- 📄 Upload bank statements (CSV, PDF)
- 🤖 AI-powered transaction classification
- 🧮 Nigeria 2026 tax calculations (PIT, CGT, Crypto)
- 📊 Detailed tax reports with export
