import { useEffect, useState, useRef } from 'react'

export interface AggregateUpdate {
  time: string
  cpuTemp: number
  coolantTemp: number
  vibration: number
  leakDetected: boolean
  smokeDetected: boolean
  flameDetected: boolean
  motorOn: boolean
  coolantLevel: number
  inFlowrate: number
  outFlowrate: number
}

// Internal command message shape
interface CommandMessage {
  type: string
  payload: any
}

export function useSensorStream() {
  const [latest, setLatest] = useState<AggregateUpdate | null>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const base = process.env.REACT_APP_API_BASE_URL || window.location.origin
    const protocol = base.startsWith('https') ? 'wss' : 'ws'
    const host = base.replace(/^https?:\/\//, '')
    const ws = new WebSocket(`${protocol}://${host}/stream`)
    wsRef.current = ws

    ws.onopen = () => console.log('WebSocket connected to', ws.url)
    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        // Distinguish between aggregate updates and command responses
        if (msg.type === 'motorStatus' && latest) {
          // update only motorOn field
          setLatest({ ...latest, motorOn: msg.payload.on })
        } else if (!msg.type) {
          // raw aggregate update
          setLatest(msg as AggregateUpdate)
        }
      } catch (e) {
        console.error('Failed to parse WS message', e)
      }
    }
    ws.onerror = (err) => console.error('WebSocket error', err)
    ws.onclose = () => console.log('WebSocket closed')

    return () => ws.close()
  }, [])

  const toggleMotor = (on: boolean) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const cmd = { type: 'setMotor', payload: { on } }
      wsRef.current.send(JSON.stringify(cmd))
    }
  }

  const autoMotor = () => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const cmd: CommandMessage = { type: 'autoMotor', payload: null }
      wsRef.current.send(JSON.stringify(cmd))
    }
  }

  return { latest, toggleMotor, autoMotor }
}
