"use client"

import { useParams } from "next/navigation"
import { Layout } from "@/components/layout"
import { AnalyticsDashboard } from "@/components/analytics-dashboard"
import { AIInsights } from "@/components/ai-insights"
import { Button } from "@/components/ui/button"
import { ArrowLeft } from "lucide-react"
import Link from "next/link"
import { useLink, useAnalytics } from "@/lib/hooks/use-api"
import ErrorBoundary from "@/components/error-boundary"

export default function AnalyticsPage() {
  const params = useParams()
  const linkId = params.id as string
  
  const { data: link, loading: linkLoading, error: linkError } = useLink(linkId)
  const { data: analytics, loading: analyticsLoading, error: analyticsError } = useAnalytics(linkId)
  
  const isLoading = linkLoading || analyticsLoading
  const error = linkError || analyticsError

  if (error) {
    return (
      <Layout>
        <div className="max-w-4xl mx-auto">
          <div className="text-center py-12">
            <h1 className="text-2xl font-bold text-gray-900 mb-4">
              Error Loading Analytics
            </h1>
            <p className="text-gray-600 mb-6">{error.message}</p>
            <Button asChild>
              <Link href="/dashboard">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Dashboard
              </Link>
            </Button>
          </div>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="outline" size="sm" asChild>
            <Link href="/dashboard">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back
            </Link>
          </Button>
        </div>

        {link && analytics ? (
          <>
            <ErrorBoundary>
              <AnalyticsDashboard 
                link={link} 
                analytics={analytics} 
                isLoading={isLoading} 
              />
            </ErrorBoundary>
            
            <ErrorBoundary>
              <AIInsights linkId={linkId} />
            </ErrorBoundary>
          </>
        ) : link && !analytics ? (
          <div className="space-y-6">
            <div className="text-center py-12">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">
                No Analytics Data Yet
              </h2>
              <p className="text-gray-600 mb-6">
                This link hasn't received any clicks yet. Share your shortened URL to start collecting analytics data.
              </p>
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 max-w-md mx-auto mb-6">
                <p className="text-sm text-blue-800">
                  <strong>Your shortened URL:</strong>
                </p>
                <a
                  href={`${window.location.origin}/${link.slug}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 hover:text-blue-800 font-mono text-sm break-all"
                >
                  {`${window.location.origin}/${link.slug}`}
                </a>
              </div>
              
              {process.env.NODE_ENV === 'development' && (
                <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 max-w-md mx-auto">
                  <p className="text-sm text-yellow-800 mb-3">
                    <strong>Development Mode:</strong> Generate sample data for testing
                  </p>
                  <Button
                    onClick={async () => {
                      try {
                        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/test/generate-clicks/${linkId}`, {
                          method: 'POST'
                        })
                        if (response.ok) {
                          window.location.reload()
                        }
                      } catch (error) {
                        console.error('Failed to generate test data:', error)
                      }
                    }}
                    size="sm"
                    variant="outline"
                  >
                    Generate Test Data
                  </Button>
                </div>
              )}
            </div>
          </div>
        ) : (
          <div className="space-y-6">
            <div className="animate-pulse">
              <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="h-24 bg-gray-200 rounded"></div>
                ))}
              </div>
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {[...Array(4)].map((_, i) => (
                  <div key={i} className="h-80 bg-gray-200 rounded"></div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  )
}