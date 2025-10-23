"use client"

import { useState, useEffect } from "react"
import { Layout } from "@/components/layout"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { 
  BarChart3, 
  MousePointer, 
  Link as LinkIcon, 
  Calendar,
  TrendingUp,
  Globe
} from "lucide-react"

interface Link {
  id: string
  original_url: string
  slug: string
  title: string
  description: string
  click_count: number
  created_at: string
  updated_at: string
}

export default function AnalyticsOverview() {
  const [links, setLinks] = useState<Link[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchLinks = async () => {
      try {
        const response = await fetch("/api/links")
        if (response.ok) {
          const data = await response.json()
          setLinks(data)
        }
      } catch (error) {
        console.error("Error fetching links:", error)
      } finally {
        setIsLoading(false)
      }
    }

    fetchLinks()
  }, [])

  const totalClicks = links.reduce((sum, link) => sum + link.click_count, 0)
  const topLinks = [...links].sort((a, b) => b.click_count - a.click_count).slice(0, 5)
  const recentLinks = [...links].sort((a, b) => 
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  ).slice(0, 5)

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    })
  }

  if (isLoading) {
    return (
      <Layout>
        <div className="max-w-6xl mx-auto space-y-6">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="h-24 bg-gray-200 rounded"></div>
              ))}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {[...Array(2)].map((_, i) => (
                <div key={i} className="h-80 bg-gray-200 rounded"></div>
              ))}
            </div>
          </div>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className="max-w-6xl mx-auto space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Analytics Overview</h1>
          <p className="text-gray-600 mt-1">
            Monitor performance across all your shortened URLs
          </p>
        </div>

        {/* Key Metrics */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total Links</p>
                  <p className="text-2xl font-bold">{links.length}</p>
                </div>
                <LinkIcon className="h-8 w-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total Clicks</p>
                  <p className="text-2xl font-bold">{totalClicks}</p>
                </div>
                <MousePointer className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Avg. Clicks</p>
                  <p className="text-2xl font-bold">
                    {links.length > 0 ? Math.round(totalClicks / links.length) : 0}
                  </p>
                </div>
                <TrendingUp className="h-8 w-8 text-purple-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Active Links</p>
                  <p className="text-2xl font-bold">
                    {links.filter(link => link.click_count > 0).length}
                  </p>
                </div>
                <Globe className="h-8 w-8 text-orange-500" />
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top Performing Links */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <TrendingUp className="h-5 w-5" />
                Top Performing Links
              </CardTitle>
              <CardDescription>
                Links with the most clicks
              </CardDescription>
            </CardHeader>
            <CardContent>
              {topLinks.length > 0 ? (
                <div className="space-y-4">
                  {topLinks.map((link, index) => (
                    <div key={link.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <Badge variant="secondary" className="text-xs">
                            #{index + 1}
                          </Badge>
                          <h4 className="font-medium truncate">
                            {link.title || "Untitled"}
                          </h4>
                        </div>
                        <p className="text-sm text-gray-500 truncate mt-1">
                          /{link.slug}
                        </p>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="text-right">
                          <p className="font-medium">{link.click_count}</p>
                          <p className="text-xs text-gray-500">clicks</p>
                        </div>
                        <Button size="sm" variant="outline" asChild>
                          <Link href={`/analytics/${link.id}`}>
                            <BarChart3 className="h-4 w-4" />
                          </Link>
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <TrendingUp className="h-12 w-12 text-gray-300 mx-auto" />
                  <p className="text-sm text-gray-500 mt-2">
                    No click data available yet
                  </p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Recent Links */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calendar className="h-5 w-5" />
                Recent Links
              </CardTitle>
              <CardDescription>
                Your most recently created links
              </CardDescription>
            </CardHeader>
            <CardContent>
              {recentLinks.length > 0 ? (
                <div className="space-y-4">
                  {recentLinks.map((link) => (
                    <div key={link.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex-1 min-w-0">
                        <h4 className="font-medium truncate">
                          {link.title || "Untitled"}
                        </h4>
                        <p className="text-sm text-gray-500 truncate mt-1">
                          /{link.slug} • {formatDate(link.created_at)}
                        </p>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="text-right">
                          <p className="font-medium">{link.click_count}</p>
                          <p className="text-xs text-gray-500">clicks</p>
                        </div>
                        <Button size="sm" variant="outline" asChild>
                          <Link href={`/analytics/${link.id}`}>
                            <BarChart3 className="h-4 w-4" />
                          </Link>
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <Calendar className="h-12 w-12 text-gray-300 mx-auto" />
                  <p className="text-sm text-gray-500 mt-2">
                    No links created yet
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </Layout>
  )
}