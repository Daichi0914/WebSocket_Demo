"use client";

import { useState } from "react";
import { useChat } from "../hooks/useChat";
import { NamePrompt } from "../components/NamePrompt";
import { ChatHeader } from "../components/ChatHeader";
import { MessageList } from "../components/MessageList";
import { MessageInput } from "../components/MessageInput";
import "./globals.css";

export default function Home() {
  const [senderName, setSenderName] = useState("");
  const [isNameSet, setIsNameSet] = useState(false);
  
  const { messages, isConnected, sendMessage } = useChat(isNameSet, senderName);

  return (
    <main>
      {!isNameSet && (
        <NamePrompt 
          senderName={senderName} 
          setSenderName={setSenderName} 
          onJoin={() => setIsNameSet(true)} 
        />
      )}

      <div className="chat-container">
        <ChatHeader isConnected={isConnected} />
        
        <MessageList 
          messages={messages} 
          currentUser={senderName} 
        />

        <MessageInput 
          disabled={!isNameSet || !isConnected} 
          onSend={sendMessage} 
        />
      </div>
    </main>
  );
}
