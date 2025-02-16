"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/router"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export default function PollCreation() {
  const router = useRouter()
  const [question, setQuestion] = useState("")
  const [options, setOptions] = useState(["", ""])

  const handleOptionChange = (index: number, value: string) => {
    const newOptions = [...options]
    newOptions[index] = value
    setOptions(newOptions)
  }

  const addOption = () => {
    setOptions([...options, ""])
  }

  const createPoll = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const response = await fetch("/api/polls", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ question, options: options.filter((opt) => opt.trim() !== "") }),
      })
      if (response.ok) {
        const data = await response.json()
        router.push(`/poll/${data.id}`)
      } else {
        console.error("Failed to create poll")
      }
    } catch (error) {
      console.error("Error creating poll:", error)
    }
  }

  return (
    <form onSubmit={createPoll} className="space-y-6">
      <div>
        <Label htmlFor="question">Question</Label>
        <Input id="question" value={question} onChange={(e) => setQuestion(e.target.value)} required />
      </div>
      {options.map((option, index) => (
        <div key={index}>
          <Label htmlFor={`option-${index}`}>Option {index + 1}</Label>
          <Input
            id={`option-${index}`}
            value={option}
            onChange={(e) => handleOptionChange(index, e.target.value)}
            required
          />
        </div>
      ))}
      <Button type="button" onClick={addOption} variant="outline">
        Add Option
      </Button>
      <Button type="submit">Create Poll</Button>
    </form>
  )
}

