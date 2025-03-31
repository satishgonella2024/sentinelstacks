# SentinelStacks Setup Guide

This guide will help you set up and run the SentinelStacks application, which consists of a Go backend API and a React frontend.

## Prerequisites

- Go 1.23+ installed
- Node.js and npm installed
- Git

## Quick Start

### 1. Setup Scripts

First, make the helper scripts executable:

```bash
# Make all helper scripts executable
bash make-run-with-backend-executable.sh
bash make-toggle-mode-executable.sh
```

### 2. Running the Application

You can run the application in two modes:

#### Mock Mode (No Backend Required)

In mock mode, the frontend uses mock data without requiring the backend server:

```bash
# Ensure the application is in mock mode
./toggle-mode.sh

# Run the frontend only
cd web-ui
npm run dev
```

#### Real Backend Mode

To use the real backend API:

```bash
# Switch to real backend mode if needed
./toggle-mode.sh

# Run the application with the backend
./run-with-backend.sh
```

### 3. Toggling Between Modes

You can easily switch between mock data and real backend:

```bash
./toggle-mode.sh
```

This script will toggle the mode and provide instructions for running the application in the selected mode.

## Troubleshooting

### Backend Build Issues

If you encounter issues building the backend:

1. Check Go dependencies:
   ```bash
   go mod tidy
   ```

2. Fix missing dependencies:
   ```bash
   go get github.com/mattn/go-isatty@v0.0.20
   ```

### Frontend Issues

If you encounter frontend issues:

1. Reinstall dependencies:
   ```bash
   cd web-ui
   npm install
   ```

2. Clear browser cache or try in incognito mode

3. Check browser console for errors

## Project Structure

- `internal/` - Backend Go packages
- `web-ui/` - Frontend React application
- `scripts/` - Helper scripts
- `cmd/` - Go command applications

## API Documentation

The backend API docs are available at:
- Swagger UI: `http://localhost:8081/swagger/`
- Health check: `http://localhost:8080/v1/health`

## Contributing

Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on contributing to the project.
