import type { ReactNode } from "react"
import { useTheme } from "next-themes"
import { Button } from "@/components/ui/button"
import { MoonIcon, SunIcon } from "@radix-ui/react-icons"

interface LayoutProps {
  children: ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const { theme, setTheme } = useTheme()

  return (
    <div className="min-h-screen bg-white dark:bg-gray-900 text-black dark:text-white">
      <nav className="p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">Quick Polling</h1>
        <Button variant="outline" size="icon" onClick={() => setTheme(theme === "dark" ? "light" : "dark")}>
          {theme === "dark" ? <SunIcon /> : <MoonIcon />}
        </Button>
      </nav>
      <main className="container mx-auto px-4 py-8">{children}</main>
    </div>
  )
}

