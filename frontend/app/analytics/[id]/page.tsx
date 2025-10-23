"use client"

import { useParams } from "next/navigation"
import { Layout } from "@/components/layout"
import { AnalyticsDashboard } from "@/components/analytics-dashboard"
import { AIInsights } from "@/components/ai-insights"
import { Button } from "@/components/ui/button"
import { ArrowLeft } from "lucide-react"
import Link from "next/link"
import { useLink, useAnalytics } from "@/lib/hooks/use-api"

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
            <AnalyticsDashboard 
              link={link} 
              analytics={analytics} 
              isLoading={isLoading} 
            />
            
            <AIInsights linkId={linkId} />
          </>
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