"use client";

import { useState } from "react";
import type { Session } from "@/lib/api";

const DEFAULT_API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export function SessionGate({
  onConnect,
}: {
  onConnect: (session: Session) => void;
}) {
  const [apiBaseUrl, setApiBaseUrl] = useState(DEFAULT_API_BASE_URL);
  const [accessToken, setAccessToken] = useState("");
  const [organizationId, setOrganizationId] = useState("");

  return (
    <div className="flex flex-1 items-center justify-center bg-zinc-50 dark:bg-black">
      <form
        className="flex w-full max-w-sm flex-col gap-4 rounded-lg border border-black/10 bg-white p-8 dark:border-white/10 dark:bg-zinc-950"
        onSubmit={(e) => {
          e.preventDefault();
          onConnect({ apiBaseUrl, accessToken, organizationId });
        }}
      >
        <div>
          <h1 className="text-lg font-semibold text-black dark:text-zinc-50">
            Connect to your Lumi server
          </h1>
          <p className="mt-1 text-sm text-zinc-500">
            Self-hosted backend. Point this dashboard at it.
          </p>
        </div>

        <label className="flex flex-col gap-1 text-sm text-zinc-700 dark:text-zinc-300">
          Server URL
          <input
            required
            className="rounded border border-black/10 bg-transparent px-3 py-2 text-sm dark:border-white/10"
            value={apiBaseUrl}
            onChange={(e) => setApiBaseUrl(e.target.value)}
            placeholder="http://localhost:8080"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm text-zinc-700 dark:text-zinc-300">
          Access token
          <input
            required
            className="rounded border border-black/10 bg-transparent px-3 py-2 text-sm dark:border-white/10"
            value={accessToken}
            onChange={(e) => setAccessToken(e.target.value)}
            placeholder="from /login"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm text-zinc-700 dark:text-zinc-300">
          Organization ID
          <input
            required
            className="rounded border border-black/10 bg-transparent px-3 py-2 text-sm dark:border-white/10"
            value={organizationId}
            onChange={(e) => setOrganizationId(e.target.value)}
            placeholder="org uuid"
          />
        </label>

        <button
          type="submit"
          className="mt-2 rounded-full bg-foreground px-5 py-2 text-sm font-medium text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc]"
        >
          Connect
        </button>
      </form>
    </div>
  );
}
