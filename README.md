# reproGolangSDKChange

Tests SDK behavior with Google Cloud Storage destination creation and retrieval.

## Quick Start

```bash
# Install dependencies
go mod tidy

# Start Cribl server
make start

# Run example
make run
```

## What It Does

The example demonstrates:
1. Creating a Google Cloud Storage destination
2. Reading it back
3. Comparing request vs response
4. Cleaning up

## Environment Variables

- `CRIBL_SERVER_URL` - Server URL (default: `http://localhost:9000`)
- `CRIBL_USERNAME` or `CRIBL_USER` - Username for authentication (default: `admin`)
- `CRIBL_PASSWORD` or `CRIBL_PASS` - Password for authentication (default: `admin`)

The code will automatically authenticate using the `/api/v1/auth/login` endpoint to obtain a bearer token.

## Docker Commands

```bash
make start    # Start server
make stop     # Stop server
make restart  # Restart server
make clean    # Remove containers and volumes
```
