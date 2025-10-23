"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { toast } from "sonner"
import { Link as LinkIcon, Loader2 } from "lucide-react"
import { useCreateLink } from "@/lib/hooks/use-api"
import { Link, APIError } from "@/lib/types"

interface URLShortenerFormProps {
  onLinkCreated?: (link: Link) => void
}

export function URLShortenerForm({ onLinkCreated }: URLShortenerFormProps) {
  const [originalUrl, setOriginalUrl] = useState("")
  const [customSlug, setCustomSlug] = useState("")
  const [errors, setErrors] = useState<{ originalUrl?: string; customSlug?: string }>({})
  
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
      </CardContent>
    </Card>
  )
}