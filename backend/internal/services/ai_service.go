package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"url-shortener-backend/internal/errors"
	"url-shortener-backend/internal/models"
)

type AIService interface {
	GenerateTitleAndDescription(ctx context.Context, url string) (*models.AIMetadata, error)
	AnalyzeClickPatterns(ctx context.Context, analytics *models.Analytics) (string, error)
}

type aiService struct {
	client *openai.Client
}

func NewAIService(apiKey string) AIService {
	client := openai.NewClient(apiKey)
	return &aiService{
		client: client,
	}
}

func (s *aiService) GenerateTitleAndDescription(ctx context.Context, targetURL string) (*models.AIMetadata, error) {
	// Validate URL format
	if _, err := url.Parse(targetURL); err != nil {
		return nil, errors.NewAppError(errors.ErrInvalidURL, "Invalid URL format", err.Error())
	}

	// Fetch page content to get context for AI generation
	pageContent, err := s.fetchPageContent(targetURL)
	if err != nil {
		// If we can't fetch content, still try to generate based on URL
		pageContent = fmt.Sprintf("URL: %s", targetURL)
	}

	// Create prompt for OpenAI
	prompt := s.createTitleDescriptionPrompt(targetURL, pageContent)

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that generates concise, descriptive titles and descriptions for web URLs. Always respond in JSON format with 'title' and 'description' fields.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   200,
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, errors.NewAppError(errors.ErrAIServiceError, "Failed to generate title and description", err.Error())
	}

	if len(resp.Choices) == 0 {
		return nil, errors.NewAppError(errors.ErrAIServiceError, "No response from AI service", "")
	}

	// Parse the response
	content := resp.Choices[0].Message.Content
	metadata, err := s.parseAIResponse(content, targetURL)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrAIServiceError, "Failed to parse AI response", err.Error())
	}

	return metadata, nil
}

func (s *aiService) AnalyzeClickPatterns(ctx context.Context, analytics *models.Analytics) (string, error) {
	if analytics == nil {
		return "", errors.NewAppError(errors.ErrInvalidInput, "Analytics data is required", "")
	}

	// Create prompt for analytics insights
	prompt := s.createAnalyticsPrompt(analytics)

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a data analyst that provides concise, actionable insights about URL click patterns. Focus on trends, patterns, and recommendations in 2-3 sentences.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   300,
			Temperature: 0.5,
		},
	)

	if err != nil {
		return "", errors.NewAppError(errors.ErrAIServiceError, "Failed to generate analytics insights", err.Error())
	}

	if len(resp.Choices) == 0 {
		return "", errors.NewAppError(errors.ErrAIServiceError, "No response from AI service", "")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func (s *aiService) fetchPageContent(targetURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "URL-Shortener-Bot/1.0")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Read first 1KB of content for context
	buffer := make([]byte, 1024)
	n, _ := resp.Body.Read(buffer)
	
	return string(buffer[:n]), nil
}

func (s *aiService) createTitleDescriptionPrompt(targetURL, pageContent string) string {
	return fmt.Sprintf(`Generate a title and description for this URL:

URL: %s

Page content preview:
%s

Please provide a JSON response with:
- "title": A concise, descriptive title (max 60 characters)
- "description": A brief description explaining what this link is about (max 160 characters)

Focus on being informative and engaging while keeping it concise.`, targetURL, pageContent)
}

func (s *aiService) createAnalyticsPrompt(analytics *models.Analytics) string {
	var prompt strings.Builder
	
	prompt.WriteString(fmt.Sprintf("Analyze these URL click analytics and provide insights:\n\n"))
	prompt.WriteString(fmt.Sprintf("Total Clicks: %d\n\n", analytics.TotalClicks))
	
	// Daily clicks trend
	if len(analytics.DailyClicks) > 0 {
		prompt.WriteString("Daily Clicks:\n")
		for _, daily := range analytics.DailyClicks {
			prompt.WriteString(fmt.Sprintf("- %s: %d clicks\n", daily.Date, daily.Count))
		}
		prompt.WriteString("\n")
	}
	
	// Top countries
	if len(analytics.CountryBreakdown) > 0 {
		prompt.WriteString("Top Countries:\n")
		for i, country := range analytics.CountryBreakdown {
			if i >= 5 { // Limit to top 5
				break
			}
			prompt.WriteString(fmt.Sprintf("- %s: %d clicks\n", country.Country, country.Count))
		}
		prompt.WriteString("\n")
	}
	
	// Device breakdown
	if len(analytics.DeviceBreakdown) > 0 {
		prompt.WriteString("Device Types:\n")
		for _, device := range analytics.DeviceBreakdown {
			prompt.WriteString(fmt.Sprintf("- %s: %d clicks\n", device.DeviceType, device.Count))
		}
		prompt.WriteString("\n")
	}
	
	// Top referrers
	if len(analytics.ReferrerBreakdown) > 0 {
		prompt.WriteString("Top Referrers:\n")
		for i, referrer := range analytics.ReferrerBreakdown {
			if i >= 3 { // Limit to top 3
				break
			}
			prompt.WriteString(fmt.Sprintf("- %s: %d clicks\n", referrer.Referrer, referrer.Count))
		}
	}
	
	prompt.WriteString("\nProvide 2-3 sentences with key insights, trends, and actionable recommendations based on this data.")
	
	return prompt.String()
}

func (s *aiService) parseAIResponse(content, fallbackURL string) (*models.AIMetadata, error) {
	// Try to extract JSON from the response
	content = strings.TrimSpace(content)
	
	// Simple JSON parsing - look for title and description
	var title, description string
	
	// Extract title
	if titleStart := strings.Index(content, `"title"`); titleStart != -1 {
		titleStart = strings.Index(content[titleStart:], `"`) + titleStart + 1
		titleStart = strings.Index(content[titleStart:], `"`) + titleStart + 1
		titleEnd := strings.Index(content[titleStart:], `"`) + titleStart
		if titleEnd > titleStart {
			title = content[titleStart:titleEnd]
		}
	}
	
	// Extract description
	if descStart := strings.Index(content, `"description"`); descStart != -1 {
		descStart = strings.Index(content[descStart:], `"`) + descStart + 1
		descStart = strings.Index(content[descStart:], `"`) + descStart + 1
		descEnd := strings.Index(content[descStart:], `"`) + descStart
		if descEnd > descStart {
			description = content[descStart:descEnd]
		}
	}
	
	// Fallback if parsing fails
	if title == "" {
		title = s.generateFallbackTitle(fallbackURL)
	}
	if description == "" {
		description = s.generateFallbackDescription(fallbackURL)
	}
	
	// Ensure length limits
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	if len(description) > 160 {
		description = description[:157] + "..."
	}
	
	return &models.AIMetadata{
		Title:       title,
		Description: description,
	}, nil
}

func (s *aiService) generateFallbackTitle(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "Shortened Link"
	}
	
	domain := parsedURL.Hostname()
	if domain == "" {
		return "Shortened Link"
	}
	
	// Remove www. prefix
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	
	return fmt.Sprintf("Link to %s", domain)
}

func (s *aiService) generateFallbackDescription(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "A shortened URL link"
	}
	
	domain := parsedURL.Hostname()
	if domain == "" {
		return "A shortened URL link"
	}
	
	// Remove www. prefix
	if strings.HasPrefix(domain, "www.") {
		domain = domain[4:]
	}
	
	return fmt.Sprintf("Visit %s via this shortened link", domain)
}