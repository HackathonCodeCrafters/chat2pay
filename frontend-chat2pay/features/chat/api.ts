import { apiClient, endpoints, BackendResponse } from "@/shared/api";
import { LLMResponse, AskProductRequest, ChatMessage } from "./types";

export async function askProduct(data: AskProductRequest): Promise<LLMResponse> {
  const response = await apiClient.post<BackendResponse<LLMResponse>>(
    endpoints.products.ask,
    { json: data }
  );

  if (!response.data.status) {
    throw new Error(response.data.error || "Failed to get response");
  }

  return response.data.data;
}

interface ChatHistoryResponse {
  messages: {
    id: string;
    role: "user" | "assistant";
    content: string;
    products?: unknown;
    created_at: string;
  }[];
}

export async function getChatHistory(token: string): Promise<ChatMessage[]> {
  const response = await apiClient.get<BackendResponse<ChatHistoryResponse>>(
    endpoints.chat.history,
    { headers: { Authorization: `Bearer ${token}` } }
  );

  if (!response.data.status || !response.data.data) {
    return [];
  }

  return response.data.data.messages.map((msg) => ({
    id: msg.id,
    role: msg.role,
    content: msg.content,
    products: msg.products as ChatMessage["products"],
    timestamp: new Date(msg.created_at),
  }));
}

export async function saveChatMessage(
  token: string,
  message: { role: string; content: string; products?: unknown }
): Promise<void> {
  await apiClient.post(endpoints.chat.saveMessage, {
    json: message,
    headers: { Authorization: `Bearer ${token}` },
  });
}

export async function clearChatHistory(token: string): Promise<void> {
  await apiClient.delete(endpoints.chat.clearHistory, {
    headers: { Authorization: `Bearer ${token}` },
  });
}
