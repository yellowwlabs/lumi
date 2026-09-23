"use client";

import { useCallback, useEffect, useState } from "react";
import { Dashboard } from "@/components/dashboard";
import { NodeSetup } from "@/components/node-setup";
import { SessionGate } from "@/components/session-gate";
import { ApiError, listNodes, type Node, type Session } from "@/lib/api";
import {
  clearSession,
  loadActiveNodeId,
  loadSession,
  saveActiveNodeId,
  saveSession,
} from "@/lib/session";

export default function Home() {
  const [session, setSession] = useState<Session | null>(null);
  const [nodes, setNodes] = useState<Node[] | null>(null);
  const [activeNodeId, setActiveNodeId] = useState<string>("");
  const [showSetup, setShowSetup] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setSession(loadSession());
  }, []);

  const refreshNodes = useCallback(async (s: Session) => {
    try {
      const list = await listNodes(s);
      setNodes(list);
      setError(null);
      if (list.length > 0) {
        const stored = loadActiveNodeId();
        const stillExists = list.some((n) => n.ID === stored);
        setActiveNodeId(stillExists ? (stored as string) : list[0].ID);
        setShowSetup(false);
      }
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        clearSession();
        setSession(null);
        return;
      }
      setError(err instanceof Error ? err.message : "failed to load nodes");
    }
  }, []);

  useEffect(() => {
    if (!session) return;
    refreshNodes(session);
    const interval = setInterval(() => refreshNodes(session), 5000);
    return () => clearInterval(interval);
  }, [session, refreshNodes]);

  if (!session) {
    return (
      <SessionGate
        onConnect={(s) => {
          saveSession(s);
          setSession(s);
        }}
      />
    );
  }

  if (nodes === null) {
    return (
      <div className="flex flex-1 items-center justify-center bg-zinc-50 dark:bg-black">
        <p className="text-sm text-zinc-500">Connecting…</p>
        {error && <p className="ml-2 text-sm text-red-500">{error}</p>}
      </div>
    );
  }

  if (nodes.length === 0 || showSetup) {
    return (
      <div className="flex flex-1 items-center justify-center bg-zinc-50 p-8 dark:bg-black">
        <NodeSetup
          session={session}
          onPaired={() => refreshNodes(session)}
          onCancel={() => setShowSetup(false)}
        />
      </div>
    );
  }

  return (
    <Dashboard
      nodes={nodes}
      activeNodeId={activeNodeId}
      onSelectNode={(id) => {
        setActiveNodeId(id);
        saveActiveNodeId(id);
      }}
      onAddNode={() => setShowSetup(true)}
    />
  );
}
