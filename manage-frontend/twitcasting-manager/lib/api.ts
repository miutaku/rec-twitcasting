import { Streamer, type StreamerResponse, ApiError } from "@/types/streamer"

const API_BASE = `${process.env.NEXT_PUBLIC_MANAGE_BACKEND_API_FQDN}:${process.env.NEXT_PUBLIC_MANAGE_BACKEND_API_PORT}`

export async function listStreamers(): Promise<StreamerResponse[]> {
  const res = await fetch(`${API_BASE}/list-casting-users`)
  if (!res.ok) throw new Error("ストリーマーリストの取得に失敗しました")
  return res.json()
}

export async function addStreamer(username: string): Promise<StreamerResponse> {
  const res = await fetch(`${API_BASE}/add-casting-user?username=${username}`, {
    method: "POST",
  })
  if (!res.ok) throw new Error("ストリーマーの追加に失敗しました")
  return res.json()
}

export async function deleteStreamer(username: string): Promise<void> {
  const res = await fetch(`${API_BASE}/del-casting-user?username=${username}`, {
    method: "DELETE",
  })
  if (!res.ok) throw new Error("ストリーマーの削除に失敗しました")
}

export async function checkRecordingState(username: string): Promise<StreamerResponse> {
  const res = await fetch(`${API_BASE}/check-recording-state?username=${username}`)
  if (!res.ok) throw new Error("録画状態の確認に失敗しました")
  return res.json()
}

