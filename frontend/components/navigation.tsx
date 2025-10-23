"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Link as LinkIcon, BarChart3, Home } from "lucide-react"

const navigation = [
  {
    name: "Home",
    href: "/",
    icon: Home,
  },
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LinkIcon,
  },
  {
    name: "Analytics",
    href: "/analytics",
    icon: BarChart3,
  },
]

export function Navigation() {
  const pathname = usePathname()

  return (
    <nav className="flex items-center space-x-1">
      {navigation.map((item) => {
        const Icon = item.icon
        const isActive = pathname === item.href || (item.href !== "/" && pathname.startsWith(item.href))
        
        return (
          <Button
            key={item.name}
            variant={isActive ? "default" : "ghost"}
            size="sm"
            asChild
          >
            <Link href={item.href} className="flex items-center gap-2">
              <Icon className="h-4 w-4" />
              <span className="hidden sm:inline">{item.name}</span>
            </Link>
          </Button>
        )
      })}
    </nav>
  )
}