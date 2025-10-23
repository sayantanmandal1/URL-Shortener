import Link from "next/link"
import { Navigation } from "./navigation"
import { Link as LinkIcon } from "lucide-react"

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container flex h-14 items-center justify-between">
        <Link href="/" className="flex items-center space-x-2">
          <LinkIcon className="h-6 w-6" />
          <span className="font-bold">URL Shortener</span>
        </Link>
        <Navigation />
      </div>
    </header>
  )
}