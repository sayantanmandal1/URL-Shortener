"use client"

import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { toast } from "sonner"
import { 
  Brain, 
  Loader2, 
  RefreshCw, 
  Lightbulb,
  TrendingUp,
  AlertCircle
} from "lucide-react"
import { useInsights } from "@/lib/hooks/use-api"
import { APIError } from "@/lib/types"

interface AIInsightsProps {
  linkId: string
  className?: string
}

export function AIInsights({ linkId, className }: AIInsightsProps) {
  const { data: insights, loading: isLoading, error, refetch } = useInsights(linkId)
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null)

  const handleRefresh = async () => {
    try {
      await refetch()
      setLastUpdated(new Date())
      toast.success("Insights refreshed successfully")
    } catch (error) {
      if (error instanceof APIError) {
        toast.error(error.message)
      } else {
        toast.error("Failed to refresh insights")
      }
    }
  }

  const formatInsights = (text: string) => {
    // Split by common patterns and format as bullet points
    const sentences = text.split(/[.!?]+/).filter(s => s.trim().length > 0)
    
    if (sentences.length <= 2) {
      return <p className="text-gray-700 leading-relaxed">{text}</p>
    }
    
    return (
      <ul className="space-y-2">
        {sentences.map((sentence, index) => (
          <li key={index} className="flex items-start gap-2">
            <Lightbulb className="h-4 w-4 text-yellow-500 mt-0.5 shrink-0" />
            <span className="text-gray-700">{sentence.trim()}.</span>
          </li>
        ))}
      </ul>
    )
  }

  const formatLastUpdated = (date: Date) => {
    const now = new Date()
    const diffInMinutes = Math.floor((now.getTime() - date.getTime()) / (1000 * 60))
    
    if (diffInMinutes < 1) {
      return "Just now"
    } else if (diffInMinutes < 60) {
      return `${diffInMinutes} minute${diffInMinutes > 1 ? 's' : ''} ago`
    } else if (diffInMinutes < 1440) {
      const hours = Math.floor(diffInMinutes / 60)
      return `${hours} hour${hours > 1 ? 's' : ''} ago`
    } else {
      return date.toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      })
    }
  }

  return (
    <Card className={className}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Brain className="h-5 w-5 text-purple-500" />
            <CardTitle>AI Insights</CardTitle>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isLoading}
          >
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <RefreshCw className="h-4 w-4" />
            )}
          </Button>
        </div>
        <CardDescription>
          AI-powered analysis of your link performance and patterns
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading && !insights ? (
          <div className="flex items-center justify-center py-8">
            <div className="text-center">
              <Loader2 className="h-8 w-8 animate-spin mx-auto text-purple-500" />
              <p className="text-sm text-gray-500 mt-2">
                Analyzing your link data...
              </p>
            </div>
          </div>
        ) : error ? (
          <div className="flex items-center gap-3 p-4 bg-red-50 border border-red-200 rounded-lg">
            <AlertCircle className="h-5 w-5 text-red-500 shrink-0" />
            <div>
              <p className="text-sm font-medium text-red-800">
                Unable to generate insights
              </p>
              <p className="text-sm text-red-600 mt-1">{error.message}</p>
            </div>
          </div>
        ) : insights ? (
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                <TrendingUp className="h-3 w-3 mr-1" />
                Analysis
              </Badge>
              {lastUpdated && (
                <span className="text-xs text-gray-500">
                  Updated {formatLastUpdated(lastUpdated)}
                </span>
              )}
            </div>
            
            <div className="prose prose-sm max-w-none">
              {formatInsights(insights)}
            </div>
            
            <div className="pt-3 border-t">
              <p className="text-xs text-gray-500">
                💡 Insights are generated based on your link's click patterns, 
                geographic distribution, and device usage data.
              </p>
            </div>
          </div>
        ) : (
          <div className="text-center py-8">
            <Brain className="h-12 w-12 text-gray-300 mx-auto" />
            <p className="text-sm text-gray-500 mt-2">
              No insights available yet
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}