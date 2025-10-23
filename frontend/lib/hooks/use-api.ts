import { useState, useEffect, useCallback } from 'react'
import { apiClient } from '../api-client'
import { APIError } from '../types'

// Generic hook for API calls with loading and error states
export function useAPI<T>(
  apiCall: () => Promise<T>,
  dependencies: any[] = []
) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<APIError | null>(null)

  const fetchData = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await apiCall()
      setData(result)
    } catch (err) {
      setError(err instanceof APIError ? err : new APIError('UNKNOWN_ERROR', 'An unknown error occurred'))
    } finally {
      setLoading(false)
    }
  }, dependencies)

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const refetch = useCallback(() => {
    fetchData()
  }, [fetchData])

  return { data, loading, error, refetch }
}

// Hook for links list
export function useLinks() {
  return useAPI(() => apiClient.getLinks())
}

// Hook for single link
export function useLink(id: string) {
  return useAPI(() => apiClient.getLink(id), [id])
}

// Hook for analytics
export function useAnalytics(linkId: string) {
  return useAPI(() => apiClient.getAnalytics(linkId), [linkId])
}

// Hook for AI insights
export function useInsights(linkId: string) {
  return useAPI(() => apiClient.getInsights(linkId), [linkId])
}

// Hook for creating links with mutation state
export function useCreateLink() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<APIError | null>(null)

  const createLink = useCallback(async (data: { original_url: string; custom_slug?: string }) => {
    try {
      setLoading(true)
      setError(null)
      const result = await apiClient.createLink(data)
      return result
    } catch (err) {
      const apiError = err instanceof APIError ? err : new APIError('UNKNOWN_ERROR', 'An unknown error occurred')
      setError(apiError)
      throw apiError
    } finally {
      setLoading(false)
    }
  }, [])

  return { createLink, loading, error }
}