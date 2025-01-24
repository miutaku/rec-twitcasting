"use server"

const API_BASE_URL = "http://manage-backend-rec-twitcasting:8080"

export async function addCastingUser(username: string) {
  const response = await fetch(`${API_BASE_URL}/add-casting-user?username=${username}`)
  return response.ok
}

export async function deleteCastingUser(username: string) {
  const response = await fetch(`${API_BASE_URL}/del-casting-user?username=${username}`)
  return response.ok
}

export async function checkRecordingState(username: string) {
  const response = await fetch(`${API_BASE_URL}/check-recording-state?username=${username}`)
  const data = await response.json()
  return data.isRecording
}

export async function updateRecordingState(username: string, state: boolean) {
  const response = await fetch(`${API_BASE_URL}/update-recording-state?username=${username}&state=${state}`)
  return response.ok
}

