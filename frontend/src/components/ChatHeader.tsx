import React from 'react';

interface ChatHeaderProps {
  isConnected: boolean;
}

export function ChatHeader({ isConnected }: ChatHeaderProps) {
  return (
    <header className="chat-header">
      <h1>Cosmic Chat {isConnected ? "🟢" : "🔴"}</h1>
    </header>
  );
}
