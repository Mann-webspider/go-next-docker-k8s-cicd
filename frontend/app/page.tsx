"use client";

import { useState, useEffect, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";

// Interface for floating +1 badges
interface FloatingText {
  id: number;
}

// Fallback to local Go backend port if environment variable is not defined
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function Home() {
  const [clickCount, setClickCount] = useState<number>(0);
  const [floatingTexts, setFloatingTexts] = useState<FloatingText[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  // 1. Fetch initial click count from Backend (GET /clicks)
  const fetchCount = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE_URL}/clicks`);
      if (res.ok) {
        const data = await res.json();
        setClickCount(data.count ?? 0);
      }
    } catch (err) {
      console.error("Failed to fetch click count:", err);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCount();
  }, [fetchCount]);

  // 2. Trigger click action (Calls POST /clicked and triggers animation)
  const handleClick = useCallback(async () => {
    // Spawn floating "+1" animation
    const newId = Date.now() + Math.random();
    setFloatingTexts((prev) => [...prev, { id: newId }]);

    try {
      const res = await fetch(`${API_BASE_URL}/clicked`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      if (res.ok) {
        const data = await res.json();
        // Update count from backend response
        setClickCount(data.count);
      }
    } catch (err) {
      console.error("Failed to register click:", err);
    }
  }, []);

  // 3. Handle Keyboard "Spacebar" presses
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Prevent page from scrolling down when pressing spacebar
      if (event.code === "Space") {
        event.preventDefault();
        handleClick();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleClick]);

  const removeFloatingText = (id: number) => {
    setFloatingTexts((prev) => prev.filter((item) => item.id !== id));
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-zinc-50 font-sans dark:bg-zinc-950 select-none">
      <div className="flex flex-col items-center justify-center space-y-8">
        
        {/* Count Display */}
        <div className="flex flex-col items-center">
          <span className="text-xs uppercase tracking-widest text-zinc-400 dark:text-zinc-500 font-semibold mb-1">
            Total Telemetry Clicks
          </span>
          <motion.h1
            key={clickCount}
            initial={{ scale: 0.9, opacity: 0.8 }}
            animate={{ scale: 1, opacity: 1 }}
            className="text-6xl font-extrabold text-zinc-900 dark:text-zinc-100 tracking-tight"
          >
            {isLoading ? "..." : clickCount.toLocaleString()}
          </motion.h1>
        </div>

        {/* Playful SaaS Button Wrapper */}
        <div className="relative flex items-center justify-center">
          
          {/* Floating "+1" Animations rising from center */}
          <AnimatePresence>
            {floatingTexts.map((item) => (
              <motion.span
                key={item.id}
                initial={{ opacity: 1, y: 0, scale: 0.8 }}
                animate={{ opacity: 0, y: -90, scale: 1.4 }}
                exit={{ opacity: 0 }}
                transition={{ duration: 0.7, ease: "easeOut" }}
                onAnimationComplete={() => removeFloatingText(item.id)}
                className="absolute font-black text-2xl text-indigo-500 dark:text-indigo-400 pointer-events-none z-20"
              >
                +1
              </motion.span>
            ))}
          </AnimatePresence>

          {/* Interactive Button */}
          <motion.button
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.92 }}
            onClick={handleClick}
            className="relative inline-flex items-center justify-center px-8 py-4 text-lg font-bold text-white transition-all duration-150 bg-indigo-600 rounded-2xl shadow-lg shadow-indigo-500/30 hover:bg-indigo-500 active:shadow-sm focus:outline-none focus:ring-4 focus:ring-indigo-300 dark:focus:ring-indigo-800"
          >
            <span className="flex items-center space-x-2">
              <span>Push Event</span>
              <span className="text-xs px-2 py-0.5 bg-indigo-700/60 rounded-md font-mono text-indigo-200">
                [Space]
              </span>
            </span>
          </motion.button>
        </div>

        {/* Hint text */}
        <p className="text-xs text-zinc-400 dark:text-zinc-600">
          Click the button or tap <kbd className="font-mono bg-zinc-200 dark:bg-zinc-800 px-1 rounded">Spacebar</kbd> repeatedly
        </p>
      </div>
    </div>
  );
}