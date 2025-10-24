package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repositories"
	"url-shortener-backend/internal/services"
	"url-shortener-backend/internal/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) (*fiber.App, *testutils.TestDB, *LinkHandler) {
	testDB := testutils.SetupTestDB(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	aiService := testutils.NewMockAIService()

	linkService := services.NewLinkService(linkRepo, aiService)
	analyticsService := services.NewAnalyticsService(analyticsRepo, linkRepo, aiService)

	handler := NewLinkHandler(linkService, analyticsService)

	app := fiber.New()

	// Setup routes
	api := app.Group("/api")
	links := api.Group("/links")
	links.Post("/", handler.CreateLink)
	links.Get("/", handler.GetLinks)
	links.Get("/:id", handler.GetLink)

	app.Get("/:slug", handler.RedirectLink)

	return app, testDB, handler
}

func TestLinkHandler_CreateLink(t *testing.T) {
	app, testDB, _ := setupTestApp(t)
	defer testDB.Cleanup(t)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid request without custom slug",
			requestBody: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "valid request with custom slug",
			requestBody: models.CreateLinkRequest{
				OriginalURL: "https://example.com/path",
				CustomSlug:  "my-custom-slug",
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "invalid URL",
			requestBody: models.CreateLinkRequest{
				OriginalURL: "not-a-valid-url",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing original URL",
			requestBody: models.CreateLinkRequest{
				CustomSlug: "test-slug",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "custom slug too short",
			requestBody: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
				CustomSlug:  "ab",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB.ClearTables(t)

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/api/links", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.expectError {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)

				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])

				data := response["data"].(map[string]interface{})
				assert.NotEmpty(t, data["id"])
				assert.Equal(t, tt.requestBody.(models.CreateLinkRequest).OriginalURL, data["original_url"])
				assert.NotEmpty(t, data["slug"])

				if tt.requestBody.(models.CreateLinkRequest).CustomSlug != "" {
					assert.Equal(t, tt.requestBody.(models.CreateLinkRequest).CustomSlug, data["slug"])
				}
			}
		})
	}
}

func TestLinkHandler_GetLinks(t *testing.T) {
	app, testDB, _ := setupTestApp(t)
	defer testDB.Cleanup(t)

	// Test empty database
	req, err := http.NewRequest("GET", "/api/links", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Empty(t, response["data"])

	// Create test links
	linkID1 := uuid.New().String()
	linkID2 := uuid.New().String()
	testDB.CreateTestLink(t, linkID1, "https://example1.com", "slug1")
	testDB.CreateTestLink(t, linkID2, "https://example2.com", "slug2")

	// Test with data
	req, err = http.NewRequest("GET", "/api/links", nil)
	require.NoError(t, err)

	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestLinkHandler_GetLink(t *testing.T) {
	app, testDB, _ := setupTestApp(t)
	defer testDB.Cleanup(t)

	// Test getting non-existent link
	req, err := http.NewRequest("GET", "/api/links/non-existent-id", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Create test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test getting existing link
	req, err = http.NewRequest("GET", "/api/links/"+linkID, nil)
	require.NoError(t, err)

	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, linkID, data["id"])
	assert.Equal(t, "https://example.com", data["original_url"])
	assert.Equal(t, "test-slug", data["slug"])
}

func TestLinkHandler_RedirectLink(t *testing.T) {
	app, testDB, _ := setupTestApp(t)
	defer testDB.Cleanup(t)

	// Test redirecting non-existent slug
	req, err := http.NewRequest("GET", "/non-existent-slug", nil)
	require.NoError(t, err)

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Create test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test redirecting existing slug
	req, err = http.NewRequest("GET", "/test-slug", nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "Test Browser")
	req.Header.Set("Referer", "https://google.com")

	resp, err = app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "https://example.com", resp.Header.Get("Location"))

	// Note: Click counting and analytics recording happen asynchronously,
	// so we can't easily test them in this unit test without adding delays
	// or making the operations synchronous for testing
}
