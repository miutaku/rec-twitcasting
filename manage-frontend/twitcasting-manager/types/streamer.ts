export interface Streamer {
  target_username: string
  recording_state: boolean
  action_date_time: string
  action: string
}

export interface StreamerResponse {
  recording_state: boolean
  action_date_time: string
  action: string
  target_username: string
}

export interface ApiError {
  error: string
  message: string
}

