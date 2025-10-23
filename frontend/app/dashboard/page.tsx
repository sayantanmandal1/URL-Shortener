"use client"

import { Layout } from "@/components/layout"
import { LinksList } from "@/components/links-list"
import { useLinks } from "@/lib/hooks/use-api"

export default function Dashboard() {
  const { data: links, loading: isLoading } = useLinks()

  const totalClicks = links?.reduce((sum, link) => sum + link.click_count, 0) || 0
  const linkCount = links?.length || 0

  return (
    <Layout>
      <div className="max-w-6xl mx-auto space-y-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600 mt-1">
              Manage and monitor your shortened URLs
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-500">Total Links</p>
            <p className="text-2xl font-bold text-gray-900">{linkCount}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-lg border">
            <h3 className="text-sm font-medium text-gray-500">Total Clicks</h3>
            <p className="text-2xl font-bold text-gray-900 mt-1">{totalClicks}</p>
          </div>
          <div className="bg-white p-6 rounded-lg border">
            <h3 className="text-sm font-medium text-gray-500">Active Links</h3>
            <p className="text-2xl font-bold text-gray-900 mt-1">{linkCount}</p>
          </div>
          <div className="bg-white p-6 rounded-lg border">
            <h3 className="text-sm font-medium text-gray-500">Avg. Clicks/Link</h3>
            <p className="text-2xl font-bold text-gray-900 mt-1">
              {linkCount > 0 ? Math.round(totalClicks / linkCount) : 0}
            </p>
          </div>
        </div>

        <LinksList links={links || []} isLoading={isLoading} />
      </div>
    </Layout>
  )
}