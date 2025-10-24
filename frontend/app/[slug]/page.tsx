import { redirect } from 'next/navigation'
import { notFound } from 'next/navigation'

interface SlugPageProps {
  params: Promise<{
    slug: string
  }>
}

async function getRedirectUrl(slug: string): Promise<string | null> {
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'https://url-shortener-h4gc.onrender.com'
    
    // This is a fallback - normally users should go directly to backend
    // But if they somehow end up here, we'll redirect them properly
    const response = await fetch(`${apiUrl}/api/links`, {
      method: 'GET',
    })

    if (response.ok) {
      const data = await response.json()
      const links = data.data || data
      const link = links.find((l: any) => l.slug === slug)
      
      if (link) {
        // Redirect to backend URL for proper tracking
        return `${apiUrl}/${slug}`
      }
    }

    return null
  } catch (error) {
    console.error('Error fetching redirect URL:', error)
    return null
  }
}

export default async function SlugPage({ params }: SlugPageProps) {
  const { slug } = await params
  
  // Get the redirect URL from the backend (this also tracks the click)
  const redirectUrl = await getRedirectUrl(slug)
  
  if (!redirectUrl) {
    notFound()
  }
  
  // Redirect to the original URL
  redirect(redirectUrl)
}

// Generate metadata for the page
export async function generateMetadata({ params }: SlugPageProps) {
  const { slug } = await params
  
  return {
    title: `Redirecting... | ${slug}`,
    description: 'You are being redirected to the original URL.',
  }
}