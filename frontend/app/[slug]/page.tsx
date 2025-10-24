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
    const response = await fetch(`${apiUrl}/${slug}`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
      // Don't follow redirects automatically
      redirect: 'manual'
    })

    if (response.status === 302 || response.status === 301) {
      const location = response.headers.get('location')
      return location
    }

    if (response.status === 404) {
      return null
    }

    // If we get a JSON response, extract the URL
    if (response.headers.get('content-type')?.includes('application/json')) {
      const data = await response.json()
      return data.original_url || data.url || null
    }

    return null
  } catch (error) {
    console.error('Error fetching redirect URL:', error)
    return null
  }
}

export default async function SlugPage({ params }: SlugPageProps) {
  const { slug } = await params
  
  // Get the redirect URL from the backend
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