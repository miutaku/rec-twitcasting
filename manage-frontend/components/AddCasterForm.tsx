"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { addCastingUser } from "../actions"

interface AddCasterFormProps {
  onAdd: (username: string) => void
}

export default function AddCasterForm({ onAdd }: AddCasterFormProps) {
  const [username, setUsername] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (username) {
      const success = await addCastingUser(username)
      if (success) {
        onAdd(username)
        setUsername("")
      }
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <Input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="配信者のユーザー名"
      />
      <Button type="submit">追加</Button>
    </form>
  )
}

