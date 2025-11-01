"use client";

import { useEffect, useState } from "react";
import Image from "next/image";

export default function Home() {
  const [backendStatus, setBackendStatus] = useState<string>("Checking...");

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    
    fetch(apiUrl)
      .then((res) => res.json())
      .then((data) => setBackendStatus(data.message))
      .catch(() => setBackendStatus("Backend unreachable"));
  }, []);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <Image
          className="dark:invert"
          src="/next.svg"
          alt="Next.js logo"
          width={100}
          height={20}
          priority
        />
        
        {/* Backend Status Check */}
      <div className="flex items-center gap-4 p-4 bg-blue-100 dark:bg-blue-900 rounded">
        <div>
          <p className="font-bold">Backend Status:</p>
          <p>{backendStatus}</p>
        </div>
        <button
          onClick={() => window.location.href = "/order/1"}
          className="ml-4 rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 dark:bg-blue-800 dark:hover:bg-blue-700"
        >
          Track Order
        </button>
      </div>

        <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
          <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50">
            Tracking Status Frontend
          </h1>
        </div>
      </main>
    </div>
  );
}