import React from 'react';

interface StatusCardProps {
  label: string;
  active: boolean;
}

export const StatusCard: React.FC<StatusCardProps> = ({ label, active }) => (
  <div
    className={`p-4 rounded-xl shadow-md flex flex-col items-center justify-center space-y-2
      ${active ? 'bg-green-100' : 'bg-red-100'}`}
  >
    <div className="text-lg font-semibold">{label}</div>
    <div
      className={`w-8 h-8 rounded-full flex items-center justify-center
        ${active ? 'bg-green-500 text-white' : 'bg-red-500 text-white'}`}
    >
      {active ? '✔️' : '❌'}
    </div>
  </div>
);
