# URL Shortener

A production-ready URL shortener with analytics and AI insights.

## Features

- 🔗 **URL Shortening** - Create short, shareable links
- 📊 **Analytics** - Track clicks, countries, devices, and referrers
- 🤖 **AI Insights** - Get intelligent analysis of your link performance
- 🎯 **Custom Slugs** - Use your own custom short URLs
- � ***Responsive UI** - Works on desktop and mobile

## Quick Start

### Backend (Go)
```bash
cd backend
cp .env.example .env
# Edit .env with your database and OpenAI API key
go mod tidy
go run main.go
```

### Frontend (Next.js)
```bash
cd frontend
npm install
cp .env.example .env.local
# Edit .env.local with your backend URL
npm run dev
```

### Docker
```bash
docker-compose up
```

## Environment Variables

### Backend (.env)
```env
DATABASE_URL=postgres://user:pass@localhost:5432/dbname
PORT=8000
OPENAI_API_KEY=your_openai_key
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8000
```

## API Usage

**Create a short link:**
```bash
curl -X POST http://localhost:8000/api/links \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com"}'
```

**Get analytics:**
```bash
curl http://localhost:8000/api/analytics/{link_id}
```

## Tech Stack

- **Backend**: Go, Fiber, PostgreSQL
- **Frontend**: Next.js, TypeScript, Tailwind CSS
- **AI**: OpenAI GPT-5
- **Analytics**: Real-time click tracking

## License

MIT License - see LICENSE file for details.