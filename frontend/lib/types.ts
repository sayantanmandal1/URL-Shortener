// TypeScript interfaces matching backend models

export interface Link {
  id: string
  original_url: string
  slug: string
  title: string
  description: string
  click_count: number
  created_at: string
  updated_at: string
}

export interface CreateLinkRequest {
  original_url: string
  custom_slug?: string
}

export interface Click {
  id: string
  link_id: string
  clicked_at: string
  ip_address: string
  user_agent: string
  referrer: string
  country: string
  device_type: string
}

export interface DailyClick {
  date: string
  count: number
}

export interface CountryStats {
  country: string
  count: number
}

export interface DeviceStats {
  device_type: string
  count: number
}

export interface ReferrerStats {
  referrer: string
  count: number
}

export interface Analytics {
  total_clicks: number
  daily_clicks: DailyClick[]
  country_breakdown: CountryStats[]
  device_breakdown: DeviceStats[]
  referrer_breakdown: ReferrerStats[]
}

// API Response types
export interface APIResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
    details?: string
  }
}

// Error types
export class APIError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: string
  ) {
    super(message)
    this.name = 'APIError'
  }
}