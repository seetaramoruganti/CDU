import React from 'react'
import { useSensorStream } from '../hooks/useSensorStream'
import { StatusCard } from './statusCard'

export default function Dashboard() {
  const { latest, toggleMotor, autoMotor } = useSensorStream()

  if (!latest) {
    return <div className="p-4">Connecting to sensor stream...</div>
  }

  // Boolean status cards
  const statusCards = [
    { label: 'Leak Detected', active: latest.leakDetected },
    { label: 'Smoke Detected', active: latest.smokeDetected },
    { label: 'Flame Detected', active: latest.flameDetected },
    { label: 'Motor On', active: latest.motorOn },
  ]

  // Numeric metric cards
  const metricCards = [
    { label: 'CPU Temp (°C)', value: latest.cpuTemp.toFixed(1) },
    { label: 'Coolant Temp (°C)', value: latest.coolantTemp.toFixed(1) },
    { label: 'Vibration', value: latest.vibration.toString() },
    { label: 'Coolant Level', value: latest.coolantLevel.toString() },
    { label: 'In Flow (L/min)', value: latest.inFlowrate.toFixed(2) },
    { label: 'Out Flow (L/min)', value: latest.outFlowrate.toFixed(2) },
  ]

  return (
    <div className="p-4 space-y-8 bg-gray-100 min-h-screen">
      {/* System Status Section */}
      <section>
        <h2 className="text-2xl font-semibold mb-4">System Status</h2>
        <div className="grid grid-cols-4 gap-4 mb-6">
          {statusCards.map((card, idx) => (
            <StatusCard key={idx} label={card.label} active={card.active} />
          ))}
        </div>
        <button
          onClick={() => toggleMotor(!latest.motorOn)}
          className={`px-4 py-2 rounded ${latest.motorOn ? 'bg-red-500' : 'bg-green-500'} text-white mr-4`}
        >
          {latest.motorOn ? 'Stop Motor' : 'Start Motor'}
        </button>
        <button
          onClick={() => autoMotor()}
          className="px-4 py-2 rounded bg-blue-500 text-white mb-8"
        >
          Auto Motor
        </button>
      </section>

      {/* Metrics Section */}
      <section>
        <h2 className="text-2xl font-semibold mb-4">Live Metrics</h2>
        <div className="grid grid-cols-3 gap-4">
          {metricCards.map((card, idx) => (
            <div key={idx} className="p-4 rounded-xl shadow-md bg-white">
              <div className="text-gray-600 mb-2">{card.label}</div>
              <div className="text-3xl font-bold">{card.value}</div>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
