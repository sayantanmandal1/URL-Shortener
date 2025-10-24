"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { toast } from "sonner"
import { Link as LinkIcon, Loader2, Copy, ExternalLink } from "lucide-react"
import { useCreateLink } from "@/lib/hooks/use-api"
import { Link, APIError } from "@/lib/types"
import { apiClient } from "@/lib/api-client"

interface URLShortenerFormProps {
  onLinkCreated?: (link: Link) => void
}

export function URLShortenerForm({ onLinkCreated }: URLShortenerFormProps) {
  const [originalUrl, setOriginalUrl] = useState("")
  const [customSlug, setCustomSlug] = useState("")
  const [errors, setErrors] = useState<{ originalUrl?: string; customSlug?: string }>({})
  const [createdLink, setCreatedLink] = useState<Link | null>(null)
  const [copied, setCopied] = useState(false)
  
  const { createLink, loading: isLoading, error: apiError } = useCreateLink()

  const validateUrl = (url: string): boolean => {
    try {
      new URL(url)
      return true
    } catch {
      return false
    }
  }

  const validateSlug = (slug: string): boolean => {
    if (!slug) return true // Optional field
    const slugRegex = /^[a-zA-Z0-9_-]+$/
    return slugRegex.test(slug) && slug.length >= 3 && slug.length <= 50
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Reset errors
    setErrors({})
    
    // Validate inputs
    const newErrors: { originalUrl?: string; customSlug?: string } = {}
    
    if (!originalUrl.trim()) {
      newErrors.originalUrl = "URL is required"
    } else if (!validateUrl(originalUrl)) {
      newErrors.originalUrl = "Please enter a valid URL"
    }
    
    if (customSlug && !validateSlug(customSlug)) {
      newErrors.customSlug = "Slug must be 3-50 characters and contain only letters, numbers, hyphens, and underscores"
    }
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors)
      return
    }
    
    try {
      const link = await createLink({
        original_url: originalUrl,
        custom_slug: customSlug || undefined,
      })
      
      // Reset form
      setOriginalUrl("")
      setCustomSlug("")
      
      // Show created link
      setCreatedLink(link)
      
      // Show success message
      toast.success("Short URL created successfully!")
      
      // Notify parent component
      onLinkCreated?.(link)
      
    } catch (error) {
      console.error("Error creating short URL:", error)
      
      // Handle specific API errors
      if (error instanceof APIError) {
        if (error.code === 'SLUG_EXISTS') {
          setErrors({ customSlug: "This custom slug is already taken" })
        } else if (error.code === 'INVALID_URL') {
          setErrors({ originalUrl: error.message })
        } else {
          toast.error(error.message)
        }
      } else {
        toast.error("Failed to create short URL")
      }
    }
  }

  const copyToClipboard = async (shortUrl: string) => {
    try {
      await navigator.clipboard.writeText(shortUrl)
      setCopied(true)
      toast.success("Short URL copied to clipboard!")
      
      // Reset copied state after 2 seconds
      setTimeout(() => setCopied(false), 2000)
    } catch (error) {
      console.error("Failed to copy:", error)
      toast.error("Failed to copy URL")
    }
  }

  const getShortUrl = (slug: string) => {
    return apiClient.getShortURL(slug)
  }

  return (
    <Card className="w-full max-w-2xl mx-auto">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <LinkIcon className="h-5 w-5" />
          Shorten URL
        </CardTitle>
        <CardDescription>
          Enter a long URL to create a short, shareable link
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="originalUrl">URL to shorten *</Label>
            <Input
              id="originalUrl"
              type="url"
              placeholder="https://example.com/very/long/url"
              value={originalUrl}
              onChange={(e) => setOriginalUrl(e.target.value)}
              className={errors.originalUrl ? "border-red-500" : ""}
            />
            {errors.originalUrl && (
              <p className="text-sm text-red-500">{errors.originalUrl}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="customSlug">Custom slug (optional)</Label>
            <Input
              id="customSlug"
              type="text"
              placeholder="my-custom-link"
              value={customSlug}
              onChange={(e) => setCustomSlug(e.target.value)}
              className={errors.customSlug ? "border-red-500" : ""}
            />
            {errors.customSlug && (
              <p className="text-sm text-red-500">{errors.customSlug}</p>
            )}
            <p className="text-sm text-muted-foreground">
              Leave empty for auto-generated slug. Must be 3-50 characters, letters, numbers, hyphens, and underscores only.
            </p>
          </div>
          
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Creating...
              </>
            ) : (
              "Create Short URL"
            )}
          </Button>
        </form>

        {/* Show created link */}
        {createdLink && (
          <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-lg">
            <div className="flex items-center gap-2 mb-3">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-sm font-medium text-green-800">Link created successfully!</span>
            </div>
            
            <div className="space-y-3">
              <div>
                <Label className="text-xs text-green-700">Short URL</Label>
                <div className="flex items-center gap-2 mt-1">
                  <a
                    href={getShortUrl(createdLink.slug)}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex-1 text-sm bg-white border border-green-300 hover:border-green-400 px-3 py-2 rounded text-blue-600 hover:text-blue-800 transition-colors font-mono flex items-center gap-1"
                    title="Click to open in new tab"
                  >
                    {getShortUrl(createdLink.slug)}
                    <ExternalLink className="h-3 w-3 shrink-0 opacity-60" />
                  </a>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => copyToClipboard(getShortUrl(createdLink.slug))}
                    className="shrink-0"
                  >
                    <Copy className="h-4 w-4" />
                    {copied ? "Copied!" : "Copy"}
                  </Button>
                </div>
              </div>
              
              <div>
                <Label className="text-xs text-green-700">Original URL</Label>
                <div className="flex items-center gap-2 mt-1">
                  <a
                    href={createdLink.original_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex-1 text-sm text-gray-600 hover:text-gray-800 truncate flex items-center gap-1"
                  >
                    {createdLink.original_url}
                    <ExternalLink className="h-3 w-3 shrink-0" />
                  </a>
                </div>
              </div>
            </div>
            
            <Button
              size="sm"
              variant="ghost"
              onClick={() => setCreatedLink(null)}
              className="mt-3 text-green-700 hover:text-green-800"
            >
              Create Another Link
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}