# Contributing to URL Shortener

Thank you for your interest in contributing to the URL Shortener project! This document provides guidelines and information for contributors.

## 🤝 How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Use the issue template** if available
3. **Provide detailed information** including:
   - Steps to reproduce the problem
   - Expected vs actual behavior
   - Environment details (OS, browser, versions)
   - Screenshots or error messages if applicable

### Suggesting Features

We welcome feature suggestions! Please:

1. **Check existing feature requests** first
2. **Describe the use case** and why it would be valuable
3. **Provide implementation ideas** if you have them
4. **Consider the scope** - smaller, focused features are easier to implement

### Code Contributions

#### Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
   ```

2. **Set up development environment**
   ```bash
   # Using Docker (recommended)
   docker-compose up -d
   
   # Or manual setup (see README.md)
   ```

3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

#### Development Workflow

1. **Make your changes**
   - Follow the coding standards (see below)
   - Write tests for new functionality
   - Update documentation if needed

2. **Test your changes**
   ```bash
   # Backend tests
   cd backend && go test ./...
   
   # Frontend tests
   cd frontend && npm test
   
   # Integration tests
   docker-compose up -d
   # Test the full application
   ```

3. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add amazing new feature"
   ```

4. **Push and create pull request**
   ```bash
   git push origin feature/your-feature-name
   ```

## 📋 Coding Standards

### Backend (Go)

**Code Style:**
- Follow standard Go formatting (`gofmt`)
- Use `golint` and `go vet` for code quality
- Follow Go naming conventions
- Use meaningful variable and function names

**Structure:**
```go
// Good: Clear, descriptive function name
func CreateShortLink(originalURL string, customSlug string) (*Link, error) {
    // Implementation
}

// Bad: Unclear abbreviations
func CrtLnk(url string, slug string) (*Link, error) {
    // Implementation
}
```

**Error Handling:**
```go
// Good: Proper error handling
result, err := someOperation()
if err != nil {
    return nil, fmt.Errorf("failed to perform operation: %w", err)
}

// Bad: Ignoring errors
result, _ := someOperation()
```

**Comments:**
```go
// CreateLink creates a new shortened link with optional custom slug.
// It validates the URL, generates metadata using AI, and stores the link.
func CreateLink(ctx context.Context, req CreateLinkRequest) (*Link, error) {
    // Implementation
}
```

### Frontend (TypeScript/React)

**Code Style:**
- Use TypeScript for all new code
- Follow React best practices and hooks patterns
- Use meaningful component and variable names
- Prefer functional components over class components

**Component Structure:**
```tsx
// Good: Well-structured component
interface LinkCardProps {
  link: Link
  onCopy: (url: string) => void
}

export function LinkCard({ link, onCopy }: LinkCardProps) {
  const handleCopyClick = () => {
    onCopy(link.short_url)
  }

  return (
    <div className="border rounded-lg p-4">
      {/* Component content */}
    </div>
  )
}

// Bad: Unclear props and structure
export function Card({ data, fn }: any) {
  return <div>{/* content */}</div>
}
```

**Hooks Usage:**
```tsx
// Good: Custom hooks for reusable logic
function useApi() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  
  const createLink = async (data: CreateLinkRequest) => {
    setLoading(true)
    try {
      const result = await apiClient.createLink(data)
      return result
    } catch (err) {
      setError(err.message)
      throw err
    } finally {
      setLoading(false)
    }
  }
  
  return { createLink, loading, error }
}
```

**Styling:**
- Use Tailwind CSS classes consistently
- Follow the existing design system
- Use shadcn/ui components when possible
- Ensure responsive design

### Database

**Migrations:**
- Always create reversible migrations
- Use descriptive migration names
- Include proper indexes for performance
- Test migrations on sample data

```sql
-- Good: Clear, descriptive migration
-- 002_add_analytics_indexes.sql
CREATE INDEX CONCURRENTLY idx_clicks_link_id_created_at 
ON clicks(link_id, created_at DESC);

-- Bad: Unclear migration
-- 002_indexes.sql
CREATE INDEX idx1 ON clicks(link_id);
```

## 🧪 Testing Guidelines

### Backend Testing

**Unit Tests:**
```go
func TestLinkService_CreateLink(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateLinkRequest
        want    *Link
        wantErr bool
    }{
        {
            name: "valid URL creates link",
            input: CreateLinkRequest{
                OriginalURL: "https://example.com",
            },
            wantErr: false,
        },
        {
            name: "invalid URL returns error",
            input: CreateLinkRequest{
                OriginalURL: "not-a-url",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Integration Tests:**
- Test API endpoints with real database
- Use test database for isolation
- Clean up test data after each test

### Frontend Testing

**Component Tests:**
```tsx
import { render, screen, fireEvent } from '@testing-library/react'
import { LinkCard } from './link-card'

describe('LinkCard', () => {
  const mockLink = {
    id: '1',
    original_url: 'https://example.com',
    slug: 'test',
    click_count: 5,
    // ... other properties
  }

  it('displays link information correctly', () => {
    render(<LinkCard link={mockLink} onCopy={jest.fn()} />)
    
    expect(screen.getByText('https://example.com')).toBeInTheDocument()
    expect(screen.getByText('5 clicks')).toBeInTheDocument()
  })

  it('calls onCopy when copy button is clicked', () => {
    const mockOnCopy = jest.fn()
    render(<LinkCard link={mockLink} onCopy={mockOnCopy} />)
    
    fireEvent.click(screen.getByText('Copy'))
    expect(mockOnCopy).toHaveBeenCalledWith(expect.stringContaining('test'))
  })
})
```

## 📝 Documentation

### Code Documentation

**Go Documentation:**
```go
// Package handlers provides HTTP request handlers for the URL shortener API.
package handlers

// LinkHandler handles HTTP requests related to link management.
type LinkHandler struct {
    linkService services.LinkService
}

// CreateLink handles POST /api/links requests to create new shortened links.
// It validates the request, creates the link, and returns the result.
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
    // Implementation
}
```

**TypeScript Documentation:**
```tsx
/**
 * API client for the URL shortener backend service.
 * Provides methods for link management and analytics.
 */
export class ApiClient {
  /**
   * Creates a new shortened link.
   * @param data - The link creation request data
   * @returns Promise resolving to the created link
   * @throws ApiError when the request fails
   */
  async createLink(data: CreateLinkRequest): Promise<Link> {
    // Implementation
  }
}
```

### README Updates

When adding new features:
- Update the features list
- Add new API endpoints to the documentation
- Update environment variables if needed
- Add troubleshooting information for common issues

## 🔄 Pull Request Process

### Before Submitting

1. **Ensure all tests pass**
   ```bash
   # Backend
   cd backend && go test ./...
   
   # Frontend
   cd frontend && npm test
   ```

2. **Run linting and formatting**
   ```bash
   # Backend
   gofmt -w .
   golint ./...
   go vet ./...
   
   # Frontend
   npm run lint
   npm run format
   ```

3. **Update documentation** if needed

4. **Test manually** in development environment

### Pull Request Template

```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows the style guidelines
- [ ] Self-review of code completed
- [ ] Code is commented, particularly in hard-to-understand areas
- [ ] Documentation has been updated
- [ ] No new warnings or errors introduced
```

### Review Process

1. **Automated checks** must pass (CI/CD pipeline)
2. **Code review** by maintainers
3. **Testing** in staging environment if applicable
4. **Approval** from at least one maintainer
5. **Merge** to main branch

## 🐛 Debugging Guidelines

### Backend Debugging

**Logging:**
```go
// Use structured logging
log.Info("Creating new link",
    "original_url", req.OriginalURL,
    "custom_slug", req.CustomSlug,
    "user_id", userID,
)

// Log errors with context
log.Error("Failed to create link",
    "error", err,
    "original_url", req.OriginalURL,
)
```

**Error Handling:**
```go
// Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to validate URL %s: %w", url, err)
}
```

### Frontend Debugging

**Console Logging:**
```tsx
// Use appropriate log levels
console.debug('API request started', { url, method })
console.info('Link created successfully', { linkId })
console.warn('Slow API response', { duration })
console.error('API request failed', { error, url })
```

**Error Boundaries:**
```tsx
class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error) {
    return { hasError: true }
  }

  componentDidCatch(error, errorInfo) {
    console.error('Component error:', error, errorInfo)
    // Report to error tracking service
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback />
    }

    return this.props.children
  }
}
```

## 🏷️ Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

**Examples:**
```
feat(api): add link analytics endpoint

fix(frontend): resolve mobile responsive issues

docs: update deployment guide

refactor(backend): simplify error handling logic

test(api): add integration tests for link creation
```

## 📞 Getting Help

If you need help or have questions:

1. **Check existing documentation** (README, this guide, code comments)
2. **Search existing issues** for similar problems
3. **Create a new issue** with detailed information
4. **Join discussions** in existing issues or pull requests

## 🎉 Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Special mentions for outstanding contributions

Thank you for contributing to the URL Shortener project! 🚀