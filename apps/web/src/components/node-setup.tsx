"use client";

import { useState } from "react";
import { createNode, type Session } from "@/lib/api";

export function NodeSetup({
  session,
  onPaired,
  onCancel,
}: {
  session: Session;
  onPaired: () => void;
  onCancel: () => void;
}) {
  const [name, setName] = useState("");
  const [pairing, setPairing] = useState<{
    nodeId: string;
    token: string;
  } | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const { node, pairing_token } = await createNode(session, name);
      setPairing({ nodeId: node.ID, token: pairing_token });
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to create node");
    } finally {
      setSubmitting(false);
    }
  }

  if (pairing) {
    const registerCmd = `curl -s -X POST ${session.apiBaseUrl}/agent/register \\
  -H 'Content-Type: application/json' \\
  -d '{"pairing_token":"${pairing.token}","hostname":"'"$(hostname)"'","operating_system":"'"$(uname -s)"'"}'`;

    const heartbeatCmd = `AGENT_TOKEN=<agent_token from register above>
while true; do
  curl -s -X POST ${session.apiBaseUrl}/agent/heartbeat -H "X-Agent-Token: $AGENT_TOKEN"
  sleep 30
done`;

    return (
      <div className="flex w-full max-w-xl flex-col gap-4 rounded-lg border border-black/10 bg-white p-8 dark:border-white/10 dark:bg-zinc-950">
        <div>
          <h2 className="text-lg font-semibold text-black dark:text-zinc-50">
            Run this on the server you want to connect
          </h2>
          <p className="mt-1 text-sm text-zinc-500">
            The pairing token is one-time and expires in 15 minutes. This
            dashboard will unlock automatically once the node comes online.
          </p>
        </div>

        <div className="flex flex-col gap-1">
          <span className="text-xs font-medium text-zinc-500">
            1. Register the agent
          </span>
          <pre className="overflow-x-auto rounded bg-black/[.05] p-3 text-xs dark:bg-white/[.06]">
            {registerCmd}
          </pre>
        </div>

        <div className="flex flex-col gap-1">
          <span className="text-xs font-medium text-zinc-500">
            2. Keep it alive (heartbeat)
          </span>
          <pre className="overflow-x-auto rounded bg-black/[.05] p-3 text-xs dark:bg-white/[.06]">
            {heartbeatCmd}
          </pre>
        </div>

        <div className="flex gap-3">
          <button
            type="button"
            onClick={onPaired}
            className="rounded-full bg-foreground px-5 py-2 text-sm font-medium text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc]"
          >
            I've run it, check status
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="rounded-full border border-black/10 px-5 py-2 text-sm dark:border-white/10"
          >
            Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <form
      onSubmit={handleCreate}
      className="flex w-full max-w-sm flex-col gap-4 rounded-lg border border-black/10 bg-white p-8 dark:border-white/10 dark:bg-zinc-950"
    >
      <div>
        <h2 className="text-lg font-semibold text-black dark:text-zinc-50">
          Connect a node
        </h2>
        <p className="mt-1 text-sm text-zinc-500">
          A node is a server that runs your self-hosted Lumi agent. You need at
          least one connected before the dashboard unlocks.
        </p>
      </div>

      <label className="flex flex-col gap-1 text-sm text-zinc-700 dark:text-zinc-300">
        Node name
        <input
          required
          className="rounded border border-black/10 bg-transparent px-3 py-2 text-sm dark:border-white/10"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="prod-server-1"
        />
      </label>

      {error && <p className="text-sm text-red-500">{error}</p>}

      <div className="flex gap-3">
        <button
          type="submit"
          disabled={submitting}
          className="rounded-full bg-foreground px-5 py-2 text-sm font-medium text-background transition-colors hover:bg-[#383838] disabled:opacity-50 dark:hover:bg-[#ccc]"
        >
          {submitting ? "Creating…" : "Generate pairing token"}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="rounded-full border border-black/10 px-5 py-2 text-sm dark:border-white/10"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
