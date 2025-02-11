"use client"

import { useState, useEffect } from "react"
import { useToast } from "@/components/ui/use-toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { listStreamers, addStreamer, deleteStreamer } from "@/lib/api"
import type { Streamer } from "@/types/streamer"
import { Loader2, Plus, Trash2 } from "lucide-react"
import { useIsMobile } from "@/hooks/use-mobile"

export default function StreamerList() {
  const [streamers, setStreamers] = useState<Streamer[]>([])
  const [newStreamer, setNewStreamer] = useState("")
  const [loading, setLoading] = useState(false)
  const { toast } = useToast()

  useEffect(() => {
    fetchStreamers()
  }, [])

  async function fetchStreamers() {
    try {
      setLoading(true)
      const data = await listStreamers()
      setStreamers(data)
    } catch (error) {
      toast({
        variant: "destructive",
        title: "エラー",
        description: "ストリーマーリストの取得に失敗しました",
      })
    } finally {
      setLoading(false)
    }
  }

  async function handleAddStreamer(e: React.FormEvent) {
    e.preventDefault()
    if (!newStreamer) return

    try {
      setLoading(true)
      await addStreamer(newStreamer)
      setNewStreamer("")
      await fetchStreamers()
      toast({
        title: "成功",
        description: "ストリーマーを追加しました",
      })
    } catch (error) {
      toast({
        variant: "destructive",
        title: "エラー",
        description: "ストリーマーの追加に失敗しました",
      })
    } finally {
      setLoading(false)
    }
  }

  async function handleDeleteStreamer(username: string) {
    try {
      setLoading(true)
      await deleteStreamer(username)
      await fetchStreamers()
      toast({
        title: "成功",
        description: "ストリーマーを削除しました",
      })
    } catch (error) {
      toast({
        variant: "destructive",
        title: "エラー",
        description: "ストリーマーの削除に失敗しました",
      })
    } finally {
      setLoading(false)
    }
  }

  const isMobile = useIsMobile()

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-2xl font-bold mb-8">TwitCasting ストリーマー管理</h1>

      <form onSubmit={handleAddStreamer} className="flex gap-4 mb-8">
        <Input
          type="text"
          placeholder="ストリーマー名を入力"
          value={newStreamer}
          onChange={(e) => setNewStreamer(e.target.value)}
          className="max-w-sm"
        />
        <Button type="submit" disabled={loading}>
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Plus className="h-4 w-4 mr-2" />}
          ストリーマーを追加
        </Button>
      </form>

      <div className="grid gap-4">
        {streamers.map((streamer) => (
          <Card key={streamer.target_username}>
            <CardContent className="flex items-center justify-between p-6">
              <div className="flex flex-col gap-1">
                <h2 className="text-xl font-semibold">
                  <a 
                  href={`https://twitcasting.tv/${streamer.target_username}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="hover:underline"
                  >
                  {streamer.target_username}
                  </a>
                </h2>

                <p className="text-sm text-muted-foreground">
                  追加日: {new Date(streamer.action_date_time).toLocaleString(window.navigator.language)}
                </p>
                <p className="text-sm text-muted-foreground">状態: {streamer.recording_state ? "配信録画中" : "配信オフライン"}</p>
              </div>
              <div className="flex items-center gap-4">
                <span className={`text-sm ${streamer.recording_state ? "text-red-500" : "text-green-500"}`}>
                  {streamer.recording_state ? "配信録画中" : "配信オフライン"}
                </span>
                <Button
                  variant="destructive"
                  size="icon"
                  onClick={() => handleDeleteStreamer(streamer.target_username)}
                  disabled={loading}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
