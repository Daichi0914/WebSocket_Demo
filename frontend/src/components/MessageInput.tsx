import React, { useState } from 'react';

interface MessageInputProps {
  disabled: boolean;
  onSend: (content: string) => void;
}

export function MessageInput({ disabled, onSend }: MessageInputProps) {
  const [inputValue, setInputValue] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || disabled) return;
    
    onSend(inputValue);
    setInputValue("");
  };

  return (
    <form className="chat-input-area" onSubmit={handleSubmit}>
      <input
        className="chat-input"
        type="text"
        placeholder="Type a message..."
        value={inputValue}
        onChange={e => setInputValue(e.target.value)}
        disabled={disabled}
      />
      <button 
        className="chat-send-btn" 
        type="submit"
        disabled={disabled || !inputValue.trim()}
      >
        Send
      </button>
    </form>
  );
}
