import { NextRequest, NextResponse } from 'next/server'

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  
  // Skip middleware for API routes, static files, and app routes
  if (
    pathname.startsWith('/api') ||
    pathname.startsWith('/_next') ||
    pathname.startsWith('/favicon.ico') ||
    pathname.startsWith('/dashboard') ||
    pathname.startsWith('/analytics') ||
    pathname === '/' ||
    pathname.includes('.')
  ) {
    return NextResponse.next()
  }

  // Extract slug from pathname (remove leading slash)
  const slug = pathname.slice(1)
  
  if (!slug) {
    return NextResponse.next()
  }

  try {
    // Make request to backend to get redirect URL
    const backendUrl = process.env.NEXT_PUBLIC_API_URL || 'https://url-shortener-h4gc.onrender.com'
    const response = await fetch(`${backendUrl}/${slug}`, {
      method: 'GET',
      redirect: 'manual'
    })

    if (response.status === 302 || response.status === 301) {
      const location = response.headers.get('Location')
      if (location) {
        return NextResponse.redirect(location)
      }
    }

    // If not found, continue to show 404 page
    return NextResponse.next()
  } catch (error) {
    console.error('Middleware redirect error:', error)
    return NextResponse.next()
  }
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}