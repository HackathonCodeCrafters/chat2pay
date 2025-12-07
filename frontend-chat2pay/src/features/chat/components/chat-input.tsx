"use client";

import * as React from "react";
import { Send, Sparkles } from "lucide-react";
import { cn } from "@/shared/lib/utils";

interface ChatInputProps {
  onSend: (message: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

export function ChatInput({
  onSend,
  disabled = false,
  placeholder = "Ketik pesan Anda...",
}: ChatInputProps) {
  const [value, setValue] = React.useState("");
  const [isFocused, setIsFocused] = React.useState(false);
  const inputRef = React.useRef<HTMLInputElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (value.trim() && !disabled) {
      onSend(value.trim());
      setValue("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  React.useEffect(() => {
    inputRef.current?.focus();
  }, []);

  return (
    <form onSubmit={handleSubmit} className="relative">
      <div 
        className={cn(
          "relative flex items-center gap-3 glass-card rounded-2xl px-4 py-2 transition-all duration-300",
          isFocused && "border-violet-500/50 shadow-lg shadow-violet-500/10"
        )}
      >
        <Sparkles className={cn(
          "h-5 w-5 shrink-0 transition-colors duration-300",
          isFocused ? "text-violet-400" : "text-muted-foreground"
        )} />
        
        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          placeholder={placeholder}
          disabled={disabled}
          className="flex-1 bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground disabled:opacity-50 py-2"
        />
        
        <button
          type="submit"
          disabled={disabled || !value.trim()}
          className={cn(
            "flex items-center justify-center h-10 w-10 rounded-xl shrink-0 transition-all duration-300",
            value.trim() && !disabled
              ? "bg-gradient-to-r from-violet-500 to-cyan-500 text-white shadow-lg shadow-violet-500/30 hover:shadow-violet-500/50"
              : "bg-white/5 text-muted-foreground"
          )}
        >
          <Send className="h-4 w-4" />
        </button>
      </div>
    </form>
  );
}
