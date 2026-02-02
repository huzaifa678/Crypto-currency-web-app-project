# Cryptocurrency Web App

A full-stack cryptocurrency trading platform with real-time market data, order management, and secure user authentication.

## Overview

This is a comprehensive cryptocurrency web application that provides:

- **User Management**: Secure registration, login, and profile management with OAuth2 (Google) integration
- **Market Data**: Real-time cryptocurrency market information with WebSocket streaming (Binance integration)
- **Trading Platform**: Create, manage, and track orders and trades
- **Wallet Management**: Secure wallet functionality for managing cryptocurrency assets
- **Transaction History**: Complete transaction tracking with audit logs
- **Fee Management**: Flexible fee configuration and calculation
- **Email Notifications**: Async email notifications via worker queue
- **Role-Based Access Control**: Admin, user, and guest roles with fine-grained permissions

## Tech Stack

### Backend

- **Language**: Go 1.24.3
- **Framework**: Gin Web Framework
- **APIs**:
  - REST API (HTTP)
  - gRPC with gRPC Gateway
  - WebSocket for real-time data
- **Database**: PostgreSQL
- **Message Queue**: Redis with Asynq (async task processing)
- **Authentication**: JWT (Paseto tokens), OAuth2
- **Real-time Data**: Binance WebSocket Stream

### Frontend

- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **HTTP Client**: Axios
- **Testing**: Playwright (E2E), Vitest (Unit)
- **Authentication**: React OAuth (Google)
- **UI Components**: Lucide React icons, React Hot Toast notifications
- **Charts**: Recharts for data visualization

### Infrastructure

- **Containerization**: Docker & Docker Compose
- **Kubernetes**: EKS (AWS Elastic Kubernetes Service) with ArgoCD and Helm
- **IaC**: Terraform
- **Frontend Deployment**: AWS CloudFront
- **Service Mesh**: Cert-Manager for SSL/TLS
- **Migrations**: golang-migrate

## Project Structure

```
.
├── api/                    # REST API handlers and server setup
├── gapi/                   # gRPC API handlers
├── db/                     # Database layer
│   ├── migrations/         # SQL migrations
│   ├── mock/              # Mock database for testing
│   ├── query/             # SQL query files
│   └── sqlc/              # Generated SQL code
├── frontend/               # React TypeScript frontend
│   ├── src/               # Source code
│   ├── e2e-tests/         # Playwright E2E tests
│   └── playwright-report/  # E2E test results
├── gapi/                   # gRPC service implementations
├── worker/                 # Async job processing (Asynq)
├── oauth2/                 # OAuth2 integration
├── token/                  # Token generation and validation
├── mail/                   # Email service
├── utils/                  # Utility functions
├── val/                    # Validation logic
├── proto-files/            # Protocol Buffer definitions
├── pb/                     # Generated Protocol Buffer code
├── eks/                    # Kubernetes manifests
├── eks-chart/              # Helm chart
├── terraform/              # Terraform infrastructure code
├── docs/                   # Swagger API documentation
├── main.go                 # Application entry point
├── docker-compose.yaml     # Local development setup
├── Dockerfile              # Container image definition
└── Makefile                # Build and development commands
```

## Getting Started

### Prerequisites

- **Go**: 1.24.3 or higher
- **Docker & Docker Compose**
- **PostgreSQL**: 16 (via Docker)
- **Redis**: 7.2 (via Docker)
- **Node.js**: 18+ (for frontend)
- **Protoc**: Protocol Buffer compiler (for gRPC development)

### Installation

#### 1. Clone the repository

```bash
git clone https://github.com/huzaifa678/Crypto-currency-web-app-project.git
cd Crypto-currency-web-app-project
```

#### 2. Set up environment variables

Create an `.env` file in the root directory:

```bash
cp app.env.example app.env  # or configure based on app.env template
```

#### 3. Start Docker services

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, and the API server.

#### 4. Run database migrations

```bash
make migrateup
```

#### 5. Start the backend server

```bash
make server
```

The API will be available at:

- REST API: `http://localhost:8081`
- gRPC: `localhost:9090`
- Swagger Docs: `http://localhost:8081/swagger/index.html`

#### 6. Start the frontend (in separate terminal)

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

## Available Commands

### Backend Commands

```bash
# Database
make createdb                # Create the crypto_db database
make dropdb                  # Drop the database
make migrateup              # Run all migrations
make migratedown            # Rollback migrations

# Development
make server                 # Run the server
make test                   # Run all tests with coverage
make sqlc                   # Generate SQL code from queries
make mock                   # Generate mocks for testing
make proto                  # Generate gRPC code from .proto files

# Infrastructure
make postgres               # Start PostgreSQL container
make redis                  # Start Redis container
make go-backend-compose     # Build and start with Docker Compose

# Tools
make evans                  # Interactive gRPC client
make go-tools               # Install Go development tools
```

### Frontend Commands

```bash
npm run dev                 # Start development server
npm run build               # Build for production
npm run lint                # Run ESLint
npm run test                # Run unit tests with Vitest
npm run test:e2e            # Run E2E tests with Playwright
npm run preview             # Preview production build
```

## Testing

### Backend Unit Tests

Run all tests:

```bash
make test
```

Run specific test file:

```bash
go test -v -cover ./api/... -short
```

Run integration tests:

```bash
go test -v -cover ./...  # Without -short flag
```

### Frontend Unit Tests

```bash
cd frontend
npm run test
```

### E2E Tests

The project includes Playwright end-to-end tests for critical user flows.

```bash
cd frontend
npm run test:e2e
```

View test results:

```bash
cd frontend
npm run test:e2e -- --ui  # Open UI mode
```

Test reports are generated in `frontend/playwright-report/`

## API Documentation

### REST API

Access Swagger UI documentation at `http://localhost:8081/swagger/index.html`

Main endpoints include:

- **Users**: `/users` - Create, read, update, delete users
- **Authentication**: `/token/renew_token` - Token refresh
- **Markets**: `/markets` - Get market data
- **Orders**: `/orders` - Manage orders
- **Trades**: `/trades` - View trades
- **Wallets**: `/wallets` - Wallet management
- **Transactions**: `/transactions` - Transaction history
- **Fees**: `/fees` - Fee configuration
- **Audit Logs**: `/audit_logs` - Activity tracking
- **WebSocket**: `/ws` - Real-time market data

### gRPC API

Connect via:

<<<<<<< HEAD
```bash
make evans
```

=======
>>>>>>> 2def249 (added gateway nginx in terraform)
## Database Schema

The project uses PostgreSQL with the following key tables:

- **users**: User accounts and profiles
- **sessions**: Active user sessions
- **verify_emails**: Email verification records
- **google_auth**: OAuth2 Google authentication
- **wallets**: User cryptocurrency wallets
- **markets**: Market data and symbols
- **orders**: Buy/sell orders
- **trades**: Completed trades
- **transactions**: Transaction records
- **fees**: Fee configurations
- **audit_logs**: System activity logs

Run migrations to set up the schema:

```bash
make migrateup
```

## Docker Deployment

### Build and Run Locally

```bash
docker compose up --build
```

### Build Docker Image

```bash
docker build -t crypto-app:latest .
```

### Run Container

```bash
docker run -p 8080:8081 -p 9090:9090 \
  -e DB_SOURCE=postgresql://root:secret@postgres:5432/crypto_db \
  -e REDIS_ADDR=redis:6379 \
  crypto-app:latest
```

## Kubernetes Deployment (EKS)

### Prerequisites

- AWS EKS cluster
- kubectl configured
- Helm 3+

### Deploy with Helm

```bash
cd eks-chart
helm upgrade --install crypto-app . -n crypto-app --create-namespace
```

### Manual Deployment

```bash
cd eks
./install.sh
```

Apply individual manifests:

```bash
kubectl apply -f eks/deployment.yaml
kubectl apply -f eks/service.yaml
kubectl apply -f eks/ingress.yaml
```

## Development Workflow

### Generate Protocol Buffers

After modifying `.proto` files:

```bash
make proto
```

### Generate SQL Code

After modifying `.sql` query files:

```bash
make sqlc
```

### Generate Mocks

For testing:

```bash
make mock
```

## Configuration

Configuration is loaded from environment variables and the `app.env` file. Key variables:

```
DB_SOURCE=postgresql://...        # Database connection
REDIS_ADDR=localhost:6379         # Redis address
SERVER_ADDRESS=0.0.0.0:8081       # REST API server
GRPC_SERVER_ADDRESS=0.0.0.0:9090  # gRPC server
PASETO_SYMMETRIC_KEY=...          # Token encryption key
SMTP_HOST=...                     # Email service host
SMTP_PASSWORD=...                 # Email password
GOOGLE_CLIENT_ID=...              # OAuth2 Google client
GOOGLE_CLIENT_SECRET=...           # OAuth2 Google secret
```

## Architecture

### API Layers

1. **REST API** (Gin): Traditional HTTP endpoints for web clients
2. **gRPC Gateway:** Converting REST API request to gRPC request
3. **gRPC API**: High-performance RPC for backend services
4. **WebSocket**: Real-time market data streaming
5. **Middleware**: Authentication, CORS, logging, error handling

### Data Flow

```
Frontend (React) → REST (gRPC Gateway)
                ↓
         Token Validation
                ↓
       Database (PostgreSQL)
   
Backend Workers → Redis Queue → Async Tasks (Email, Notifications)
```

## Security

- **JWT/Paseto Tokens**: Secure token-based authentication
- **OAuth2**: Google login integration
- **Password Hashing**: Bcrypt encryption
- **CORS**: Cross-origin request configuration
- **HTTPS/TLS**: SSL/TLS certificates via cert-manager
- **Role-Based Access Control**: Admin, User, Guest roles
- **Audit Logging**: Complete activity tracking

## Performance

- **Connection Pooling**: pgxpool for efficient database access
- **Caching**: Redis for session and temporary data
- **Async Processing**: Asynq for non-blocking tasks
- **WebSocket**: Real-time data without polling
- **Pagination**: Efficient data retrieval with pagination

## Monitoring & Logging

The application uses structured logging with rs/zerolog:

```go
log.Info().Msg("Server started")
log.Error().Err(err).Send()
```

Logs are output to stdout and can be aggregated in production environments.

## Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes
3. Run tests: `make test` (backend), `npm run test` (frontend)
4. Run E2E tests: `npm run test:e2e` (frontend)
5. Commit and push
6. Create a pull request

## License

This project is licensed under the LICENSE file included in the repository.

## Support

For issues and questions, please open an issue on GitHub.

## CI/CD Status

- ✅ **Unit Tests**: Backend and frontend tests passing
- ✅ **Integration Tests**: Database and API integration tests passing
- ✅ **E2E Tests**: User flow tests passing via Playwright
- ✅ **Lint**: Code quality checks passing
- ✅ **Build**: Docker image building successfully
