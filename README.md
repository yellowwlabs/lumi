# Lumi

Lumi is a self-hosted infrastructure platform for managing and deploying applications across your own servers.

## Overview

Lumi provides a unified control plane for managing multiple servers, applications, services, deployments, logs, and infrastructure.

Your infrastructure stays under your control. Lumi provides the tools to manage it.

```text
Lumi
├── Control Plane
├── Node Agent
├── Applications
├── Deployments
├── Services
└── Infrastructure
```

## Architecture

```text
Web / Desktop / Mobile / CLI
             |
             v
       Lumi Control Plane
             |
             v
         Node Agent
             |
             v
      Docker / Runtime
             |
       +-----+-----+
       |     |     |
      Apps  DBs  Services
```

## Core Features

- Self-hosted infrastructure management
- Multi-server management
- Application deployments
- Service management
- Real-time logs
- Server and application monitoring
- Secrets and environment variables
- Node registration and management
- Deployment history and rollback
- CLI support

## Tech Stack

- Go
- PostgreSQL
- Valkey
- Docker
- WebSockets / gRPC
- Next.js
- Tauri
- Expo

## Project Structure

```text
lumi/
├── apps/
│   ├── web/
│   ├── desktop/
│   └── mobile/
├── backend/
├── cli/
└── README.md
```

## Status

Lumi is currently under active development.

The current focus is building the backend, authentication, organization management, node registry, and node-agent architecture.

## License

License information will be added later.