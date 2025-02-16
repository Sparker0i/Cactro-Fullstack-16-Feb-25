"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { Button } from "@/components/ui/button"

interface Poll {
  id: string
  question: string
}

export default function PollList() {
  const [polls, setPolls] = useState<Poll[]>([])

  useEffect(() => {
    fetchPolls()
  }, [])

  const fetchPolls = async () => {
    try {
      const response = await fetch("/api/polls")
      const data = await response.json()
      setPolls(data)
    } catch (error) {
      console.error("Error fetching polls:", error)
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-semibold">Available Polls</h2>
        <Link href="/create">
          <Button>Create New Poll</Button>
        </Link>
      </div>
      <ul className="space-y-4">
        {polls.map((poll) => (
          <li key={poll.id} className="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
            <Link href={`/poll/${poll.id}`}>
              <span className="text-lg hover:underline">{poll.question}</span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}

