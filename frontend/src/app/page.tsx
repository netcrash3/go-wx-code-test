"use client";

import dynamic from "next/dynamic";

const Map = dynamic(() => import("@/components/Map"), { ssr: false });

export default function Home() {
  return (
    <main className="flex h-screen flex-col p-6">
      <h1 className="mb-4 text-3xl font-bold">Go-Weather</h1>
      <div className="flex-1 overflow-hidden rounded-lg border border-gray-200 shadow-sm">
        <Map />
      </div>
    </main>
  );
}
