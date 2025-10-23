"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer
} from "recharts"
import {
  MousePointer,
  Calendar,
  Globe,
  Smartphone,
  ExternalLink,
  TrendingUp
} from "lucide-react"
import { apiClient } from "@/lib/api-client"
import { Analytics, Link } from "@/lib/types"

interface AnalyticsDashboardProps {
  link: Link
  analytics: Analytics
  isLoading?: boolean
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D']

export function AnalyticsDashboard({ link, analytics, isLoading }: AnalyticsDashboardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
    })
  }

  const getShortUrl = (slug: string) => {
    return apiClient.getShortURL(slug)
  }

  if (isLoading) {
    return (
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
    )
  }

  return (
    <div className="space-y-6">
      {/* Link Info */}
      <div className="space-y-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            {link.title || "Untitled Link"}
          </h1>
          {link.description && (
            <p className="text-gray-600 mt-1">{link.description}</p>
          )}
        </div>
        
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex items-center gap-2">
            <Badge variant="secondary">Short URL</Badge>
            <code className="text-sm bg-gray-100 px-2 py-1 rounded">
              {getShortUrl(link.slug)}
            </code>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="outline">Original</Badge>
            <a
              href={link.original_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm text-blue-600 hover:text-blue-800 flex items-center gap-1 truncate max-w-md"
            >
              {link.original_url}
              <ExternalLink className="h-3 w-3 shrink-0" />
            </a>
          </div>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Total Clicks</p>
                <p className="text-2xl font-bold">{analytics.total_clicks}</p>
              </div>
              <MousePointer className="h-8 w-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Countries</p>
                <p className="text-2xl font-bold">{analytics.country_breakdown?.length || 0}</p>
              </div>
              <Globe className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Devices</p>
                <p className="text-2xl font-bold">{analytics.device_breakdown?.length || 0}</p>
              </div>
              <Smartphone className="h-8 w-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Created</p>
                <p className="text-2xl font-bold">
                  {new Date(link.created_at).toLocaleDateString("en-US", {
                    month: "short",
                    day: "numeric",
                  })}
                </p>
              </div>
              <Calendar className="h-8 w-8 text-orange-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Daily Clicks Chart */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Daily Clicks
            </CardTitle>
            <CardDescription>
              Click activity over time
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={analytics.daily_clicks || []}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis 
                  dataKey="date" 
                  tickFormatter={formatDate}
                />
                <YAxis />
                <Tooltip 
                  labelFormatter={(value) => formatDate(value)}
                  formatter={(value) => [value, "Clicks"]}
                />
                <Area 
                  type="monotone" 
                  dataKey="count" 
                  stroke="#0088FE" 
                  fill="#0088FE" 
                  fillOpacity={0.3}
                />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Country Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Globe className="h-5 w-5" />
              Top Countries
            </CardTitle>
            <CardDescription>
              Clicks by country
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={(analytics.country_breakdown || []).slice(0, 5)}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="country" />
                <YAxis />
                <Tooltip formatter={(value) => [value, "Clicks"]} />
                <Bar dataKey="count" fill="#00C49F" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Device Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Smartphone className="h-5 w-5" />
              Device Types
            </CardTitle>
            <CardDescription>
              Clicks by device type
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={analytics.device_breakdown || []}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ device_type, count }) => {
                    const total = (analytics.device_breakdown || []).reduce((sum, item) => sum + item.count, 0)
                    const percentage = total > 0 ? Math.round((count / total) * 100) : 0
                    return `${device_type} (${percentage}%)`
                  }}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="count"
                >
                  {(analytics.device_breakdown || []).map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value) => [value, "Clicks"]} />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Referrer Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ExternalLink className="h-5 w-5" />
              Top Referrers
            </CardTitle>
            <CardDescription>
              Traffic sources
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {(analytics.referrer_breakdown || []).slice(0, 5).map((referrer, index) => {
                const total = (analytics.referrer_breakdown || []).reduce((sum, item) => sum + item.count, 0)
                const percentage = total > 0 ? Math.round((referrer.count / total) * 100) : 0
                
                return (
                  <div key={referrer.referrer} className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <div 
                        className="w-3 h-3 rounded-full" 
                        style={{ backgroundColor: COLORS[index % COLORS.length] }}
                      />
                      <span className="text-sm font-medium truncate max-w-[200px]">
                        {referrer.referrer || "Direct"}
                      </span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-600">{referrer.count}</span>
                      <Badge variant="secondary" className="text-xs">
                        {percentage}%
                      </Badge>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}