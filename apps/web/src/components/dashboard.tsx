"use client";

import type { Node } from "@/lib/api";

const STATUS_COLOR: Record<Node["Status"], string> = {
  online: "bg-green-500",
  pending: "bg-yellow-500",
  offline: "bg-zinc-400",
};

export function Dashboard({
  nodes,
  activeNodeId,
  onSelectNode,
  onAddNode,
}: {
  nodes: Node[];
  activeNodeId: string;
  onSelectNode: (nodeId: string) => void;
  onAddNode: () => void;
}) {
  const activeNode = nodes.find((n) => n.ID === activeNodeId) ?? nodes[0];

  return (
    <div className="flex flex-1 flex-col bg-zinc-50 dark:bg-black">
      <header className="flex items-center justify-between border-b border-black/10 px-8 py-4 dark:border-white/10">
        <h1 className="text-lg font-semibold text-black dark:text-zinc-50">
          Lumi
        </h1>

        <div className="flex items-center gap-3">
          <select
            className="rounded border border-black/10 bg-transparent px-3 py-1.5 text-sm dark:border-white/10"
            value={activeNode?.ID}
            onChange={(e) => onSelectNode(e.target.value)}
          >
            {nodes.map((n) => (
              <option key={n.ID} value={n.ID}>
                {n.Name} · {n.Status}
              </option>
            ))}
          </select>
          <button
            type="button"
            onClick={onAddNode}
            className="rounded-full border border-black/10 px-4 py-1.5 text-sm dark:border-white/10"
          >
            + Add node
          </button>
        </div>
      </header>

      <main className="flex flex-1 flex-col gap-4 p-8">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {nodes.map((n) => (
            <button
              type="button"
              key={n.ID}
              onClick={() => onSelectNode(n.ID)}
              className={`flex flex-col gap-2 rounded-lg border p-4 text-left transition-colors ${
                n.ID === activeNode?.ID
                  ? "border-black/30 dark:border-white/30"
                  : "border-black/10 dark:border-white/10"
              }`}
            >
              <div className="flex items-center justify-between">
                <span className="font-medium text-black dark:text-zinc-50">
                  {n.Name}
                </span>
                <span
                  className={`h-2 w-2 rounded-full ${STATUS_COLOR[n.Status]}`}
                />
              </div>
              <span className="text-sm text-zinc-500">
                {n.Hostname || "not registered yet"}
              </span>
              <span className="text-xs text-zinc-400">
                {n.OperatingSystem || "—"}
              </span>
            </button>
          ))}
        </div>

        {activeNode && (
          <div className="mt-4 rounded-lg border border-black/10 p-6 dark:border-white/10">
            <p className="text-sm text-zinc-500">
              Managing{" "}
              <span className="font-medium text-black dark:text-zinc-50">
                {activeNode.Name}
              </span>{" "}
              — status: {activeNode.Status}.
            </p>
          </div>
        )}
      </main>
    </div>
  );
}
