# 🛡️ Kylon

**"Building Unbreakable Foundations for K8s Data."**

---

## 🚀 Project Overview

**Kylon** is an open-source platform focused on delivering intelligent, reliable, and secure Kubernetes backup and recovery. It leverages a modern, full-stack **TypeScript** architecture using **Next.js**, **tRPC**, **Supabase**, and **PostgreSQL**—designed for multi-cluster and multi-cloud environments.

Kylon’s mission is to enable **proactive data resilience** through intuitive interfaces, fine-grained control, and future-ready AI-powered insights for stateful and stateless Kubernetes workloads.

---

## 🧰 Tech Stack

| Layer              | Technology                                 |
|--------------------|---------------------------------------------|
| Monorepo Tooling   | [NX](https://nx.dev)                        |
| Frontend/Backend   | [Next.js](https://nextjs.org)               |
| API Layer          | [tRPC](https://trpc.io)                     |
| Authentication     | [Supabase Auth](https://supabase.com/auth) |
| Database           | PostgreSQL (via Supabase)                   |
| State Management   | [Zustand](https://zustand-demo.pmnd.rs/)    |
| Remote Data Sync   | [React Query](https://tanstack.com/query)  |
| Styling/UX         | SASS, Framer Motion                         |
| Hosting/Infra      | Supabase, Kubernetes, Object Storage (S3/GCS/etc.) |

---

## ✨ Key Features

### ✅ Multi-Cluster Support
- Manage `kubeconfig`s securely from one dashboard
- Run scoped backups/restores across multiple Kubernetes clusters

### ✅ Cloud-Agnostic Object Storage
- Support for AWS S3, Azure Blob, GCP Storage, and other S3-compatible endpoints

### ✅ Stateless & Stateful K8s Backups
- CRDs, Secrets, ConfigMaps, Deployments, Services, etc.
- Persistent Volume snapshots using CSI
- App-aware DB backups (via `kubectl exec`)

### ✅ Unified Dashboard
- Real-time dashboards for operations, logs, and retention management
- UI for creating backup policies and triggers

### 🔐 Supabase Auth Integration
- User sessions, role-based access control (planned)
- Simple onboarding and management via JWT-based sessions

### 🤖 AI-Driven Features (Planned)
- **Anomaly Detection**: Predict and highlight backup failures
- **Resource Selection**: Recommend backup priorities
- **NLP Assistant**: Ask for backup actions in natural language
- **Restore Path Optimization**: Smart restore strategies based on metrics

---

## 🏗️ Architecture

Kylon follows a unified full-stack TypeScript architecture:

```mermaid
graph TD
    A["🧭 Next.js Frontend (React + SASS + Zustand)"] -->|tRPC| B["⚙️ Next.js Backend (API Routes + tRPC)"]
    B -->|Securely Stores| C1["📂 Supabase DB (PostgreSQL)"]
    B -->|Connects via Kubeconfig| C2["☸️ K8s Clusters"]
    B -->|Uploads/Fetches Snapshots| D["🪣 Object Storage (AWS/GCP/Azure)"]
    B -->|Inference Calls (Planned)| E["🧠 AI/ML Service (External/In-Pod)"]
    E -->|Training Logs| F["📈 Backup Metrics / Failure History"]
