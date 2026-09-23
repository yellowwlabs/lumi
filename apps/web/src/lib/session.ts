"use client";

import type { Session } from "./api";

const KEY = "lumi.session";
const ACTIVE_NODE_KEY = "lumi.activeNodeId";

export function loadSession(): Session | null {
  if (typeof window === "undefined") return null;
  const raw = window.localStorage.getItem(KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Session;
  } catch {
    return null;
  }
}

export function saveSession(session: Session) {
  window.localStorage.setItem(KEY, JSON.stringify(session));
}

export function clearSession() {
  window.localStorage.removeItem(KEY);
  window.localStorage.removeItem(ACTIVE_NODE_KEY);
}

export function loadActiveNodeId(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(ACTIVE_NODE_KEY);
}

export function saveActiveNodeId(nodeId: string) {
  window.localStorage.setItem(ACTIVE_NODE_KEY, nodeId);
}
