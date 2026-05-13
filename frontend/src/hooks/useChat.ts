import { useEffect, useState, useRef, useCallback } from "react";
import { Message } from "../types/chat";

export function useChat(isNameSet: boolean, senderName: string) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const ws = useRef<WebSocket | null>(null);

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";
  const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws";

  useEffect(() => {
    if (!isNameSet) return;

    // Fetch history
    fetch(`${API_URL}/messages`)
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
      try {
        const msg: Message = JSON.parse(event.data);
        setMessages(prev => [...prev, msg]);
      } catch (err) {
        console.error("Failed to parse message:", err);
      }
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

  const sendMessage = useCallback((content: string) => {
    if (!content.trim() || !ws.current) return;

    const msg = {
      sender: senderName,
      content: content
    };

    ws.current.send(JSON.stringify(msg));
  }, [senderName]);

  return {
    messages,
    isConnected,
    sendMessage
  };
}
