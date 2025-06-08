# Codex Sample Project

This repository contains a simple Go backend and a React frontend. To run the project locally:

1. Start the Go backend on port `8080`:

```bash
go run ./backend
```

2. In another terminal, start the React frontend with Vite:

```bash
cd frontend
npm install # first time only
npm run dev
```

The Vite dev server proxies API requests under `/api` to the Go backend, so posting comments on the board page should work when both servers are running.
