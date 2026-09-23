export type NodeStatus = "pending" | "online" | "offline";

export type Node = {
  ID: string;
  OrganizationID: string;
  Name: string;
  Hostname: string;
  Status: NodeStatus;
  OperatingSystem: string;
  AgentVersion: string;
  LastActive: string;
  CreatedAt: string;
  UpdatedAt: string;
};

export type CreateNodeResponse = {
  node: Node;
  pairing_token: string;
};

export type Session = {
  apiBaseUrl: string;
  accessToken: string;
  organizationId: string;
};

class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function request<T>(
  session: Session,
  path: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetch(`${session.apiBaseUrl}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${session.accessToken}`,
      ...init?.headers,
    },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new ApiError(res.status, body.error ?? res.statusText);
  }

  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export function listNodes(session: Session): Promise<Node[]> {
  return request<Node[]>(
    session,
    `/organizations/${session.organizationId}/nodes`,
  );
}

export function createNode(
  session: Session,
  name: string,
): Promise<CreateNodeResponse> {
  return request<CreateNodeResponse>(
    session,
    `/organizations/${session.organizationId}/nodes`,
    { method: "POST", body: JSON.stringify({ name }) },
  );
}

export { ApiError };
