import React, { useEffect, useRef } from 'react';
import { Message } from '../types/chat';

interface MessageListProps {
  messages: Message[];
  currentUser: string;
}

export function MessageList({ messages, currentUser }: MessageListProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="chat-messages">
      {messages.map((msg, idx) => {
        const isMine = msg.sender === currentUser;
        return (
          <div key={idx} className={`message-wrapper ${isMine ? 'mine' : 'other'}`}>
            <span className="message-sender">{msg.sender}</span>
            <div className="message-bubble">{msg.content}</div>
          </div>
        );
      })}
      <div ref={messagesEndRef} />
    </div>
  );
}
