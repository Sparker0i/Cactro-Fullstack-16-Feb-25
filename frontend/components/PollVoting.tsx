import { useState } from "react"
import { Button } from "@/components/ui/button"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Label } from "@/components/ui/label"

interface PollVotingProps {
  pollId: string
  question: string
  options: string[]
  onVote: () => void
}

export default function PollVoting({ pollId, question, options, onVote }: PollVotingProps) {
  const [selectedOption, setSelectedOption] = useState<string | null>(null)

  const handleVote = async () => {
    if (!selectedOption) return

    try {
      const response = await fetch(`/api/polls/${pollId}/vote`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ option: selectedOption }),
      })
      if (response.ok) {
        onVote()
      } else {
        console.error("Failed to submit vote")
      }
    } catch (error) {
      console.error("Error submitting vote:", error)
    }
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-semibold">{question}</h2>
      <RadioGroup value={selectedOption || ""} onValueChange={setSelectedOption}>
        {options.map((option, index) => (
          <div key={index} className="flex items-center space-x-2">
            <RadioGroupItem value={option} id={`option-${index}`} />
            <Label htmlFor={`option-${index}`}>{option}</Label>
          </div>
        ))}
      </RadioGroup>
      <Button onClick={handleVote} disabled={!selectedOption}>
        Submit Vote
      </Button>
    </div>
  )
}

