import { 
  Link, 
  CreateLinkRequest, 
  Analytics, 
  Click, 
  APIResponse, 
  APIError 
} from './types'

class URLShortenerAPIClient {
  private baseURL: string

  constructor(baseURL?: string) {
    // Use environment variable or default to localhost for development
    this.baseURL = baseURL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
  }

  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    }

    try {
      const response = await fetch(url, config)
      
      // Handle non-JSON responses (like redirects)
      const contentType = response.headers.get('content-type')
      if (!contentType?.includes('application/json')) {
        if (!response.ok) {
          throw new APIError(
            'HTTP_ERROR',
            `HTTP ${response.status}: ${response.statusText}`
          )
        }
        return response as unknown as T
      }

      const data: APIResponse<T> = await response.json()

      if (!response.ok) {
        throw new APIError(
          data.error?.code || 'HTTP_ERROR',
          data.error?.message || `HTTP ${response.status}: ${response.statusText}`,
          data.error?.details
        )
      }

      if (!data.success && data.error) {
        throw new APIError(
          data.error.code,
          data.error.message,
          data.error.details
        )
      }

      return data.data as T
    } catch (error) {
      if (error instanceof APIError) {
        throw error
      }
      
      // Handle network errors
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new APIError(
          'NETWORK_ERROR',
          'Unable to connect to the server. Please check your internet connection.'
        )
      }
      
      throw new APIError(
        'UNKNOWN_ERROR',
        error instanceof Error ? error.message : 'An unknown error occurred'
      )
    }
  }

  // Link management methods
  async createLink(data: CreateLinkRequest): Promise<Link> {
    return this.request<Link>('/api/links', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getLinks(): Promise<Link[]> {
    return this.request<Link[]>('/api/links')
  }

  async getLink(id: string): Promise<Link> {
    return this.request<Link>(`/api/links/${id}`)
  }

  // Analytics methods
  async getAnalytics(linkId: string): Promise<Analytics> {
    return this.request<Analytics>(`/api/analytics/${linkId}`)
  }

  async getClicks(linkId: string): Promise<Click[]> {
    return this.request<Click[]>(`/api/analytics/${linkId}/clicks`)
  }

  // AI insights methods
  async getInsights(linkId: string): Promise<string> {
    return this.request<string>(`/api/insights/${linkId}`)
  }

  // Health check methods
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    return this.request<{ status: string; timestamp: string }>('/health')
  }

  // Utility method to get the short URL
  getShortURL(slug: string): string {
    // For client-side, use window.location.origin if available
    if (typeof window !== 'undefined') {
      return `${window.location.origin}/${slug}`
    }
    
    // For server-side rendering, use the base URL
    const baseURL = process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000'
    return `${baseURL}/${slug}`
  }
}

// Create and export a singleton instance
export const apiClient = new URLShortenerAPIClient()

// Export the class for custom instances if needed
export { URLShortenerAPIClient }