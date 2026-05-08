"use client";

import { useEffect, useState, useRef } from "react";
import "./globals.css";

interface Message {
  id?: number;
  sender: string;
  content: string;
  created_at?: string;
}

export default function Home() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState("");
  const [senderName, setSenderName] = useState("");
  const [isNameSet, setIsNameSet] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  
  const ws = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";

  useEffect(() => {
    if (!isNameSet) return;

    // Fetch history
    fetch(`${API_URL}/api/messages`)
      .then(res => res.json())
      .then(data => {
        if(Array.isArray(data)) {
          setMessages(data);
        }
      })
      .catch(console.error);

    // Setup WebSocket
    const socket = new WebSocket(WS_URL);
    
    socket.onopen = () => {
      setIsConnected(true);
      console.log("Connected to WS");
    };

    socket.onmessage = (event) => {
      const msg: Message = JSON.parse(event.data);
      setMessages(prev => [...prev, msg]);
    };

    socket.onclose = () => {
      setIsConnected(false);
      console.log("Disconnected from WS");
    };

    ws.current = socket;

    return () => {
      socket.close();
    };
  }, [isNameSet, API_URL, WS_URL]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const sendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputValue.trim() || !ws.current) return;

    const msg = {
      sender: senderName,
      content: inputValue
    };

    ws.current.send(JSON.stringify(msg));
    setInputValue("");
  };

  const handleNameSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (senderName.trim()) {
      setIsNameSet(true);
    }
  };

  return (
    <main>
      {!isNameSet && (
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
      )}

      <div className="chat-container">
        <header className="chat-header">
          <h1>Cosmic Chat {isConnected ? "🟢" : "🔴"}</h1>
        </header>

        <div className="chat-messages">
          {messages.map((msg, idx) => {
            const isMine = msg.sender === senderName;
            return (
              <div key={idx} className={`message-wrapper ${isMine ? 'mine' : 'other'}`}>
                <span className="message-sender">{msg.sender}</span>
                <div className="message-bubble">{msg.content}</div>
              </div>
            );
          })}
          <div ref={messagesEndRef} />
        </div>

        <form className="chat-input-area" onSubmit={sendMessage}>
          <input
            className="chat-input"
            type="text"
            placeholder="Type a message..."
            value={inputValue}
            onChange={e => setInputValue(e.target.value)}
            disabled={!isNameSet || !isConnected}
          />
          <button 
            className="chat-send-btn" 
            type="submit"
            disabled={!isNameSet || !isConnected || !inputValue.trim()}
          >
            Send
          </button>
        </form>
      </div>
    </main>
  );
}
