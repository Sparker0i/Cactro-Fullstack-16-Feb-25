"use client"

import { useState, useEffect } from "react"

interface PollResult {
  option: string
  votes: number
}

interface PollResultsProps {
  pollId: string
}

export default function PollResults({ pollId }: PollResultsProps) {
  const [results, setResults] = useState<PollResult[]>([])

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const response = await fetch(`/api/polls/${pollId}/results`)
        const data = await response.json()
        setResults(data)
      } catch (error) {
        console.error("Error fetching poll results:", error)
      }
    }

    fetchResults()
    const interval = setInterval(fetchResults, 5000) // Refresh every 5 seconds

    return () => clearInterval(interval)
  }, [pollId])

  const totalVotes = results.reduce((sum, result) => sum + result.votes, 0)

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-semibold">Poll Results</h2>
      {results.map((result, index) => (
        <div key={index} className="space-y-2">
          <div className="flex justify-between">
            <span>{result.option}</span>
            <span>{((result.votes / totalVotes) * 100).toFixed(1)}%</span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
            <div
              className="bg-blue-600 h-2.5 rounded-full"
              style={{ width: `${(result.votes / totalVotes) * 100}%` }}
            ></div>
          </div>
        </div>
      ))}
      <p className="text-sm text-gray-500 dark:text-gray-400">Total votes: {totalVotes}</p>
    </div>
  )
}

