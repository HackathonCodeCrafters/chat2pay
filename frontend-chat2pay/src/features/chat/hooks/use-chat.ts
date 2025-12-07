"use client";

import * as React from "react";
import { ChatMessage, Product } from "../types";
import { getChatHistory, saveChatMessage, clearChatHistory } from "../api";
import { chatWebSocket, getWebSocketBaseUrl } from "../websocket";
import { useAuth } from "@/features/auth";

function generateId() {
  return Math.random().toString(36).substring(2, 9);
}

export function useChat() {
  const { token, isAuthenticated, user } = useAuth();
  const [messages, setMessages] = React.useState<ChatMessage[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [isConnected, setIsConnected] = React.useState(false);
  const [isHistoryLoaded, setIsHistoryLoaded] = React.useState(false);

  // Connect to WebSocket
  React.useEffect(() => {
    if (!user?.id) return;

    const wsBaseUrl = getWebSocketBaseUrl();

    chatWebSocket.connect({
      url: wsBaseUrl,
      userId: user.id,
      onMessage: (data) => {
        setIsLoading(false);
        
        let content = data;
        let products: Product[] | undefined;

        // Try to parse as JSON (products array)
        try {
          const parsed = JSON.parse(data);
          if (Array.isArray(parsed)) {
            products = parsed;
            content = "Berikut produk yang saya temukan untuk Anda:";
          }
        } catch {
          // Not JSON, use as plain text
        }

        const assistantMessage: ChatMessage = {
          id: generateId(),
          role: "assistant",
          content,
          products,
          timestamp: new Date(),
        };

        setMessages((prev) => [...prev, assistantMessage]);

        // Save assistant message to database
        if (token) {
          saveChatMessage(token, {
            role: "assistant",
            content,
            products,
          }).catch(console.error);
        }
      },
      onConnect: () => {
        setIsConnected(true);
        setError(null);
      },
      onDisconnect: () => {
        setIsConnected(false);
      },
      onError: () => {
        setError("Koneksi terputus. Mencoba menghubungkan ulang...");
      },
    });

    return () => {
      chatWebSocket.disconnect();
    };
  }, [user?.id, token]);

  // Load messages from database on mount
  React.useEffect(() => {
    async function loadHistory() {
      if (!token || !isAuthenticated || isHistoryLoaded) return;
      
      try {
        const history = await getChatHistory(token);
        if (history.length > 0) {
          setMessages(history);
        }
      } catch (err) {
        console.error("Failed to load chat history:", err);
      } finally {
        setIsHistoryLoaded(true);
      }
    }
    loadHistory();
  }, [token, isAuthenticated, isHistoryLoaded]);

  const sendMessage = React.useCallback(async (content: string) => {
    if (!content.trim() || isLoading) return;

    setError(null);

    // Add user message
    const userMessage: ChatMessage = {
      id: generateId(),
      role: "user",
      content: content.trim(),
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setIsLoading(true);

    // Save user message to database
    if (token) {
      saveChatMessage(token, { role: "user", content: content.trim() }).catch(console.error);
    }

    // Send via WebSocket
    if (chatWebSocket.isConnected()) {
      chatWebSocket.send(content.trim());
    } else {
      setError("Tidak terhubung ke server. Silakan refresh halaman.");
      setIsLoading(false);
      
      const errorMessage: ChatMessage = {
        id: generateId(),
        role: "assistant",
        content: "Maaf, koneksi terputus. Silakan refresh halaman dan coba lagi.",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    }
  }, [isLoading, token]);

  const clearMessages = React.useCallback(async () => {
    setMessages([]);
    setError(null);
    
    // Clear from database
    if (token) {
      try {
        await clearChatHistory(token);
      } catch (err) {
        console.error("Failed to clear chat history:", err);
      }
    }
  }, [token]);

  return {
    messages,
    isLoading,
    error,
    isConnected,
    sendMessage,
    clearMessages,
  };
}
