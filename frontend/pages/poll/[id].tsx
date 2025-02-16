"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/router"
import Layout from "@/components/Layout"
import PollVoting from "@/components/PollVoting"
import PollResults from "@/components/PollResults"

interface Poll {
  id: string
  question: string
  options: string[]
}

export default function PollPage() {
  const router = useRouter()
  const { id } = router.query
  const [poll, setPoll] = useState<Poll | null>(null)
  const [hasVoted, setHasVoted] = useState(false)

  useEffect(() => {
    if (id) {
      fetchPoll()
    }
  }, [id])

  const fetchPoll = async () => {
    try {
      const response = await fetch(`/api/polls/${id}`)
      const data = await response.json()
      setPoll(data)
    } catch (error) {
      console.error("Error fetching poll:", error)
    }
  }

  const handleVote = () => {
    setHasVoted(true)
  }

  if (!poll) {
    return <Layout>Loading...</Layout>
  }

  return (
    <Layout>
      <h1 className="text-3xl font-bold mb-6">{poll.question}</h1>
      {!hasVoted ? (
        <PollVoting pollId={poll.id} question={poll.question} options={poll.options} onVote={handleVote} />
      ) : (
        <PollResults pollId={poll.id} />
      )}
    </Layout>
  )
}

