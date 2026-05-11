import React from 'react';

interface NamePromptProps {
  senderName: string;
  setSenderName: (name: string) => void;
  onJoin: () => void;
}

export function NamePrompt({ senderName, setSenderName, onJoin }: NamePromptProps) {
  const handleNameSubmit = (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (senderName.trim()) {
      onJoin();
    }
  };

  return (
    <div className="name-input-overlay">
      <form className="name-input-box" onSubmit={handleNameSubmit}>
        <h2>Enter your name to join</h2>
        <div className="chat-input-area" style={{ border: 'none', padding: 0, background: 'transparent' }}>
          <input
            className="chat-input"
            type="text"
            placeholder="Your Name..."
            value={senderName}
            onChange={e => setSenderName(e.target.value)}
            autoFocus
          />
          <button className="chat-send-btn" type="submit">Join</button>
        </div>
      </form>
    </div>
  );
}
