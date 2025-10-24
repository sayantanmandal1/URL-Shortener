# URL Shortener Application

A production-ready end-to-end URL shortener application with analytics and AI-powered insights.

## 🌟 Features

- 🔗 **URL Shortening**: Create short URLs from long URLs with custom or auto-generated slugs
- 📊 **Analytics Dashboard**: Comprehensive analytics with interactive charts and visualizations
- 🤖 **AI Integration**: AI-powered title generation and performance insights using OpenAI
- 📱 **Responsive Design**: Modern web interface built with Next.js and Tailwind CSS
- 🚀 **Production Ready**: Docker deployment support for Render and Vercel
- 🗄️ **Robust Database**: PostgreSQL with optimized queries and connection pooling
- ⚡ **High Performance**: Fast redirects with click tracking and metadata capture
- 🔒 **Secure**: Input validation, rate limiting, and security headers

## 🛠️ Tech Stack

### Backend
- **Go 1.21+** with Fiber v2 framework for high-performance HTTP server
- **PostgreSQL** with pgx driver for reliable data persistence
- **OpenAI API** for intelligent title generation and analytics insights
- **Docker** for containerization and deployment

### Frontend
- **Next.js 14+** with App Router for modern React development
- **TypeScript** for type safety and better developer experience
- **Tailwind CSS** for utility-first styling
- **shadcn/ui** component library for consistent UI components
- **Recharts** for interactive analytics visualizations

### Infrastructure
- **Render** for backend deployment with Docker
- **Vercel** for frontend deployment with edge optimization
- **Docker Compose** for local development environment

## 🚀 Quick Start

### Prerequisites

- **Docker and Docker Compose** (recommended for easiest setup)
- **Node.js 18+** (for frontend development)
- **Go 1.21+** (for backend development)
- **PostgreSQL 13+** (or use Docker)
- **OpenAI API Key** (for AI features)

### 🐳 Local Development with Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd url-shortener
   ```

2. **Set up environment variables**
   ```bash
   # Backend configuration
   cp backend/.env.example backend/.env
   
   # Frontend configuration
   cp frontend/.env.example frontend/.env.local
   ```

3. **Configure your environment files**
   
   Edit `backend/.env`:
   ```env
   DATABASE_URL=postgres://postgres:postgres@db:5432/url_shortener?sslmode=disable
   PORT=8080
   OPENAI_API_KEY=your_openai_api_key_here
   ENVIRONMENT=development
   ```
   
   Edit `frontend/.env.local`:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NODE_ENV=development
   ```

4. **Start the application**
   ```bash
   docker-compose up -d
   ```

5. **Run database migrations**
   ```bash
   # Wait for containers to start, then run migrations
   docker-compose exec backend ./main migrate
   ```

6. **Access the application**
   - 🌐 **Frontend**: http://localhost:3000
   - 🔧 **Backend API**: http://localhost:8080
   - ❤️ **Health Check**: http://localhost:8080/health
   - 📊 **API Documentation**: http://localhost:8080/docs (if Swagger UI is enabled)

### 🔧 Manual Development Setup

#### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Using Docker for PostgreSQL
   docker run --name url-shortener-db \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=url_shortener \
     -p 5432:5432 \
     -d postgres:15
   ```

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env`:
   ```env
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable
   PORT=8080
   OPENAI_API_KEY=your_openai_api_key_here
   ENVIRONMENT=development
   ```

5. **Run database migrations**
   ```bash
   # Run migrations (if migration command is available)
   go run main.go migrate
   
   # Or manually run SQL files
   psql -h localhost -U postgres -d url_shortener -f migrations/001_initial_schema.sql
   ```

6. **Start the backend server**
   ```bash
   go run main.go
   ```

#### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install Node.js dependencies**
   ```bash
   npm install
   # or
   yarn install
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env.local
   ```
   
   Edit `.env.local`:
   ```env
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NODE_ENV=development
   ```

4. **Start the development server**
   ```bash
   npm run dev
   # or
   yarn dev
   ```

5. **Access the frontend**
   - Open http://localhost:3000 in your browser

## ⚙️ Environment Variables

### Backend Configuration (`backend/.env`)

```env
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable

# Server Configuration
PORT=8080
HOST=0.0.0.0

# OpenAI API Configuration (Required for AI features)
OPENAI_API_KEY=your_openai_api_key_here

# Environment Settings
ENVIRONMENT=development

# Optional: Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60

# Optional: CORS Settings
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
```

### Frontend Configuration (`frontend/.env.local`)

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080

# Environment
NODE_ENV=development

# Optional: Analytics
NEXT_PUBLIC_VERCEL_ANALYTICS_ID=your_analytics_id

# Optional: Error Tracking
NEXT_PUBLIC_SENTRY_DSN=your_sentry_dsn
```

### Production Environment Variables

#### Backend (Render)
```env
DATABASE_URL=postgres://user:password@host:port/database?sslmode=require
PORT=8080
OPENAI_API_KEY=your_openai_api_key_here
ENVIRONMENT=production
CORS_ALLOWED_ORIGINS=https://yourdomain.vercel.app
```

#### Frontend (Vercel)
```env
NEXT_PUBLIC_API_URL=https://your-backend-url.onrender.com
NODE_ENV=production
```

## 📡 API Endpoints

### 🔗 Link Management
| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `POST` | `/api/links` | Create a new short link | `{"original_url": "https://example.com", "custom_slug": "optional"}` |
| `GET` | `/api/links` | Get all links | - |
| `GET` | `/api/links/:id` | Get a specific link by ID | - |
| `GET` | `/:slug` | Redirect to original URL (with click tracking) | - |

### 📊 Analytics & Insights
| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/api/analytics/:id` | Get detailed analytics for a link | Click counts, geographic data, device breakdown |
| `GET` | `/api/insights/:id` | Get AI-generated insights | AI analysis of click patterns and trends |

### ❤️ Health & Monitoring
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check endpoint |

### 📖 API Documentation

Complete API documentation is available in OpenAPI format:
- **OpenAPI Spec**: `backend/api/openapi.yaml`
- **Interactive Docs**: Available when Swagger UI is configured

#### Example API Usage

**Create a short link:**
```bash
curl -X POST http://localhost:8080/api/links \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://www.example.com/very/long/url"}'
```

**Get analytics:**
```bash
curl http://localhost:8080/api/analytics/123e4567-e89b-12d3-a456-426614174000
```

**Get AI insights:**
```bash
curl http://localhost:8080/api/insights/123e4567-e89b-12d3-a456-426614174000
```

## 🗄️ Database Schema

The application uses PostgreSQL with optimized schema design:

### Main Tables

#### `links` table
- Stores URL mappings, metadata, and basic statistics
- Indexed on `slug` for fast redirects
- Includes AI-generated titles and descriptions

#### `clicks` table  
- Stores detailed click analytics and metadata
- Captures IP, user agent, referrer, country, device type
- Indexed on `link_id` and `clicked_at` for efficient queries

### Schema Details

```sql
-- Links table
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_url TEXT NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(200),
    description TEXT,
    click_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Clicks table for detailed analytics
CREATE TABLE clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link_id UUID REFERENCES links(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(2),
    device_type VARCHAR(20)
);

-- Performance indexes
CREATE INDEX idx_links_slug ON links(slug);
CREATE INDEX idx_clicks_link_id ON clicks(link_id);
CREATE INDEX idx_clicks_clicked_at ON clicks(clicked_at);
```

### Migrations

Database migrations are located in `backend/migrations/`:
- `001_initial_schema.sql` - Initial database schema
- Run migrations with: `go run main.go migrate` (if implemented)

## 🚀 Deployment

### 🔧 Backend Deployment (Render)

The backend is configured for deployment on Render using Docker:

1. **Create a new Web Service on Render**
   - Connect your GitHub repository
   - Select "Docker" as the environment
   - Use the root directory

2. **Configure Environment Variables**
   ```env
   DATABASE_URL=postgres://user:password@host:port/database?sslmode=require
   PORT=8080
   OPENAI_API_KEY=your_openai_api_key_here
   ENVIRONMENT=production
   CORS_ALLOWED_ORIGINS=https://yourdomain.vercel.app
   ```

3. **Deploy**
   - Render will automatically build and deploy using `backend/Dockerfile`
   - The service will be available at `https://your-service-name.onrender.com`

4. **Database Setup**
   - Create a PostgreSQL database on Render
   - Run migrations after first deployment
   - Update `DATABASE_URL` with the connection string

### 🌐 Frontend Deployment (Vercel)

The frontend is configured for deployment on Vercel:

1. **Connect Repository to Vercel**
   - Import your GitHub repository
   - Vercel will auto-detect Next.js configuration

2. **Configure Build Settings**
   - Root Directory: `frontend`
   - Build Command: `npm run build`
   - Output Directory: `.next`

3. **Set Environment Variables**
   ```env
   NEXT_PUBLIC_API_URL=https://your-backend-url.onrender.com
   NODE_ENV=production
   ```

4. **Deploy**
   - Vercel will automatically deploy on push to main branch
   - Available at `https://your-project.vercel.app`

### 🔄 CI/CD Pipeline

**Automatic Deployments:**
- Frontend: Deploys automatically on push to `main` branch
- Backend: Deploys automatically on push to `main` branch
- Both services support preview deployments for pull requests

**Manual Deployment:**
```bash
# Deploy frontend to Vercel
cd frontend && vercel --prod

# Deploy backend (push to main branch or manual deploy on Render)
git push origin main
```

## 🏗️ Development

### Project Structure

```
url-shortener/
├── 📁 backend/                    # Go backend application
│   ├── 📁 internal/
│   │   ├── 📁 handlers/          # HTTP request handlers
│   │   ├── 📁 services/          # Business logic layer
│   │   ├── 📁 repositories/      # Data access layer
│   │   ├── 📁 models/            # Data models and structs
│   │   ├── 📁 middleware/        # HTTP middleware
│   │   ├── 📁 config/            # Configuration management
│   │   ├── 📁 errors/            # Error handling
│   │   └── 📁 utils/             # Utility functions
│   ├── 📁 migrations/            # Database migration files
│   ├── 📁 api/                   # API documentation
│   ├── 📄 Dockerfile            # Production Docker config
│   ├── 📄 .dockerignore         # Docker ignore file
│   ├── 📄 go.mod                # Go module definition
│   └── 📄 main.go               # Application entry point
├── 📁 frontend/                   # Next.js frontend application
│   ├── 📁 app/                   # Next.js app router pages
│   ├── 📁 components/            # Reusable React components
│   ├── 📁 lib/                   # Utility functions and API client
│   ├── 📁 public/                # Static assets
│   ├── 📄 vercel.json           # Vercel deployment config
│   ├── 📄 Dockerfile.dev        # Development Docker config
│   ├── 📄 package.json          # Node.js dependencies
│   └── 📄 tailwind.config.ts    # Tailwind CSS configuration
├── 📄 docker-compose.yml         # Local development setup
├── 📄 .gitignore                # Git ignore file
└── 📄 README.md                 # This documentation
```

### 🔧 Development Workflow

1. **Setup Development Environment**
   ```bash
   # Clone and setup
   git clone <repository-url>
   cd url-shortener
   
   # Start with Docker (recommended)
   docker-compose up -d
   
   # Or setup manually (see Manual Development Setup above)
   ```

2. **Making Changes**
   ```bash
   # Create feature branch
   git checkout -b feature/your-feature-name
   
   # Make your changes
   # Backend: Edit files in backend/
   # Frontend: Edit files in frontend/
   
   # Test your changes
   # Backend: go test ./...
   # Frontend: npm test
   ```

3. **Database Changes**
   ```bash
   # Create new migration
   # Add SQL file to backend/migrations/
   
   # Apply migrations
   docker-compose exec backend ./main migrate
   ```

### 🧪 Testing

**Backend Testing:**
```bash
cd backend
go test ./...                    # Run all tests
go test -v ./internal/services/  # Run specific package tests
go test -cover ./...             # Run with coverage
```

**Frontend Testing:**
```bash
cd frontend
npm test                         # Run Jest tests
npm run test:watch              # Run tests in watch mode
npm run test:coverage           # Run with coverage
```

### 🐛 Troubleshooting

#### Common Issues

**1. Database Connection Issues**
```bash
# Check if PostgreSQL is running
docker-compose ps

# Check database logs
docker-compose logs db

# Reset database
docker-compose down -v
docker-compose up -d
```

**2. Backend Build Issues**
```bash
# Clean Go module cache
go clean -modcache
go mod tidy

# Rebuild Docker image
docker-compose build backend
```

**3. Frontend Build Issues**
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Rebuild Docker image
docker-compose build frontend
```

**4. OpenAI API Issues**
- Verify your API key is correct in `.env`
- Check your OpenAI account has sufficient credits
- Ensure API key has proper permissions

**5. CORS Issues**
- Check `CORS_ALLOWED_ORIGINS` in backend `.env`
- Verify frontend URL matches allowed origins
- Check browser developer tools for CORS errors

#### Debug Mode

**Backend Debug:**
```bash
# Run with debug logging
ENVIRONMENT=development go run main.go

# Or with Docker
docker-compose up backend
docker-compose logs -f backend
```

**Frontend Debug:**
```bash
# Run with debug mode
npm run dev

# Check browser developer tools for errors
# Network tab for API call issues
# Console tab for JavaScript errors
```

### 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
   - Follow existing code style
   - Add tests for new functionality
   - Update documentation if needed
4. **Test your changes**
   ```bash
   # Backend
   cd backend && go test ./...
   
   # Frontend  
   cd frontend && npm test
   ```
5. **Commit your changes**
   ```bash
   git commit -m "Add amazing feature"
   ```
6. **Push to your branch**
   ```bash
   git push origin feature/amazing-feature
   ```
7. **Submit a pull request**

### 📋 Code Style Guidelines

**Backend (Go):**
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors appropriately

**Frontend (TypeScript/React):**
- Use TypeScript for type safety
- Follow React best practices
- Use meaningful component and variable names
- Add JSDoc comments for complex functions

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Express-inspired web framework for Go
- [Next.js](https://nextjs.org/) - React framework for production
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [shadcn/ui](https://ui.shadcn.com/) - Beautiful UI components
- [OpenAI](https://openai.com/) - AI-powered features
## 
🔗 How Click Tracking Works

### Shortened URL Flow
1. **Clickable Links**: All shortened URLs in the interface are clickable and open in new tabs
2. **Click Recording**: When clicked, the request goes to the backend at `/:slug`
3. **Analytics Capture**: Backend records click metadata (IP, user agent, referrer, device type, country)
4. **Redirect**: After recording analytics, backend redirects to the original URL
5. **Real-time Updates**: Click counts and analytics update immediately

### Click Metadata Captured
- **IP Address**: For geographic analysis
- **User Agent**: For device type detection (Mobile, Desktop, Tablet, Bot)
- **Referrer**: To track traffic sources
- **Timestamp**: For daily/hourly analytics
- **Country**: Extracted from IP (simplified implementation)

### Analytics Features
- Total click counts
- Daily click trends
- Geographic distribution
- Device type breakdown
- Referrer analysis
- AI-powered insights

All shortened URLs are fully functional and will properly track clicks when accessed from any browser or application.