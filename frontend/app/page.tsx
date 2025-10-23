"use client"

import { Layout } from "@/components/layout"
import { URLShortenerForm } from "@/components/url-shortener-form"
import { LinksList } from "@/components/links-list"
import { useLinks } from "@/lib/hooks/use-api"
import { Link } from "@/lib/types"

export default function Home() {
  const { data: links, loading: isLoading, refetch } = useLinks()

  const handleLinkCreated = (newLink: Link) => {
    // Refetch the links to get the updated list
    refetch()
  }

  return (
    <Layout>
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-4xl font-bold text-gray-900">
            URL Shortener
          </h1>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Create short, shareable links with detailed analytics and AI-powered insights. 
            Track clicks, analyze traffic patterns, and optimize your link performance.
          </p>
        </div>

        <URLShortenerForm onLinkCreated={handleLinkCreated} />
        
        <LinksList links={links || []} isLoading={isLoading} />
      </div>
    </Layout>
  )
}
