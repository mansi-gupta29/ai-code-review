# AI Code Review API

A Go REST API that generates AI-powered code reviews using an LLM (Groq).

## Features
- `POST /review` — submit code, get an AI-generated review
- `GET /reviews` — list past reviews
- `GET /reviews/{id}` — fetch a specific review
- Reviews persisted in PostgreSQL

## Tech Stack
- Go, Chi router
- PostgreSQL (pgx)
- Groq LLM API

## Architecture
Layered design — handlers (HTTP), ai (LLM client), store (database), models (types).

## Running
\`\`\`bash
export GROQ_API_KEY=your-key
export DATABASE_URL=postgres://localhost:5432/codereview
go run cmd/api/main.go
\`\`\`