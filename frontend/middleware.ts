import { NextRequest, NextResponse } from 'next/server'

export async function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  
  // Skip middleware for API routes, static files, and known app routes
  if (
    pathname.startsWith('/api') ||
    pathname.startsWith('/_next') ||
    pathname.startsWith('/favicon.ico') ||
    pathname.startsWith('/dashboard') ||
    pathname.startsWith('/analytics') ||
    pathname === '/' ||
    pathname.includes('.') ||
    pathname.startsWith('/public')
  ) {
    return NextResponse.next()
  }

  // For all other paths, let the [slug] route handle it
  // This allows the client-side redirect logic to work properly
  return NextResponse.next()
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