"use client"

import { useState } from "react"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { toast } from "sonner"
import { 
  Copy, 
  ExternalLink, 
  BarChart3, 
  Calendar,
  MousePointer,
  Link as LinkIcon
} from "lucide-react"
import { apiClient } from "@/lib/api-client"
import { Link as LinkType } from "@/lib/types"

interface LinksListProps {
  links: LinkType[]
  isLoading?: boolean
}

export function LinksList({ links, isLoading }: LinksListProps) {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const copyToClipboard = async (shortUrl: string, linkId: string) => {
    try {
      await navigator.clipboard.writeText(shortUrl)
      setCopiedId(linkId)
      toast.success("Short URL copied to clipboard!")
      
      // Reset copied state after 2 seconds
      setTimeout(() => setCopiedId(null), 2000)
    } catch (error) {
      console.error("Failed to copy:", error)
      toast.error("Failed to copy URL")
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    })
  }

  const getShortUrl = (slug: string) => {
    return apiClient.getShortURL(slug)
  }

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Your Links</CardTitle>
          <CardDescription>Loading your shortened URLs...</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  if (links.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Your Links</CardTitle>
          <CardDescription>No shortened URLs yet</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8">
            <LinkIcon className="mx-auto h-12 w-12 text-gray-400" />
            <h3 className="mt-2 text-sm font-medium text-gray-900">No links</h3>
            <p className="mt-1 text-sm text-gray-500">
              Get started by creating your first short URL above.
            </p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Your Links</CardTitle>
        <CardDescription>
          Manage and track your shortened URLs
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {links.map((link) => {
            const shortUrl = getShortUrl(link.slug)
            const isCopied = copiedId === link.id
            
            return (
              <div
                key={link.id}
                className="border rounded-lg p-4 space-y-3 hover:bg-gray-50 transition-colors"
              >
                {/* Title and Description */}
                <div>
                  <h3 className="font-medium text-gray-900 truncate">
                    {link.title || "Untitled"}
                  </h3>
                  {link.description && (
                    <p className="text-sm text-gray-600 mt-1 line-clamp-2">
                      {link.description}
                    </p>
                  )}
                </div>

                {/* URLs */}
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <Badge variant="secondary" className="text-xs">
                      Short URL
                    </Badge>
                    <code className="text-sm bg-gray-100 px-2 py-1 rounded flex-1 truncate">
                      {shortUrl}
                    </code>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => copyToClipboard(shortUrl, link.id)}
                      className="shrink-0"
                    >
                      <Copy className="h-4 w-4" />
                      {isCopied ? "Copied!" : "Copy"}
                    </Button>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <Badge variant="outline" className="text-xs">
                      Original
                    </Badge>
                    <a
                      href={link.original_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm text-blue-600 hover:text-blue-800 flex-1 truncate flex items-center gap-1"
                    >
                      {link.original_url}
                      <ExternalLink className="h-3 w-3 shrink-0" />
                    </a>
                  </div>
                </div>

                {/* Stats and Actions */}
                <div className="flex items-center justify-between pt-2 border-t">
                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <div className="flex items-center gap-1">
                      <MousePointer className="h-4 w-4" />
                      <span>{link.click_count} clicks</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      <span>{formatDate(link.created_at)}</span>
                    </div>
                  </div>
                  
                  <Button size="sm" variant="outline" asChild>
                    <Link href={`/analytics/${link.id}`}>
                      <BarChart3 className="h-4 w-4 mr-1" />
                      Analytics
                    </Link>
                  </Button>
                </div>
              </div>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}