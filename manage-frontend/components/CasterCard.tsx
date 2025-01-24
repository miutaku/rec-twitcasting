"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { deleteCastingUser, checkRecordingState, updateRecordingState } from "../actions"

interface CasterCardProps {
  username: string
  onDelete: () => void
}

export default function CasterCard({ username, onDelete }: CasterCardProps) {
  const [isRecording, setIsRecording] = useState(false)

  const handleCheckState = async () => {
    const state = await checkRecordingState(username)
    setIsRecording(state)
  }

  const handleUpdateState = async (newState: boolean) => {
    const success = await updateRecordingState(username, newState)
    if (success) {
      setIsRecording(newState)
    }
  }

  const handleDelete = async () => {
    const success = await deleteCastingUser(username)
    if (success) {
      onDelete()
    }
  }

  return (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>{username}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>配信状態: {isRecording ? "配信中" : "配信していない"}</p>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button onClick={handleCheckState}>状態確認</Button>
        <Button onClick={() => handleUpdateState(!isRecording)}>{isRecording ? "配信停止" : "配信開始"}</Button>
        <Button variant="destructive" onClick={handleDelete}>
          削除
        </Button>
      </CardFooter>
    </Card>
  )
}

