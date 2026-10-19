# Cloud-Native Telemetry & Event Pipeline 

![CI/CD Workflow](https://img.shields.io/badge/CI%2FCD-GitHub_Actions-blue?logo=githubactions)
![Docker](https://img.shields.io/badge/Containerization-Docker_Buildx-2496ED?logo=docker)
![Kubernetes](https://img.shields.io/badge/Orchestration-Kubernetes-326CE5?logo=kubernetes)

![Go](https://img.shields.io/badge/Backend-Golang-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Frontend-Nextjs_14-000000?logo=nextdotjs)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-4169E1?logo=postgresql)

A high-concurrency event telemetry system designed to showcase production-grade DevOps practices—from multi-stage container optimization and automated CI/CD pipelines to Kubernetes ingress routing and zero-downtime deployment strategies.

---

##  Architecture & Component Flow

```text
               +-------------------------------------------------+
               |                Nginx Ingress                    |
               +-----------------------+-------------------------+
                                       |
                   +-------------------+-------------------+
                   | /                                     | /api
                   v                                       v
      +-------------------------+             +-------------------------+
      |    Next.js Frontend     |             |      Golang API         |
      |   (Standalone Cluster)  |             |  (Concurrency Scaled)   |
      +-------------------------+             +------------+------------+
                                                           |
                                                           v
                                              +-------------------------+
                                              |    PostgreSQL DB        |
                                              | (StatefulSet + Volume)  |
                                              +-------------------------+
```
## Key DevOps Practices & Engineering Highlights
### 1. Multi-Stage & Secure Containerization
- **Optimized Image Footprint**: Utilized multi-stage Docker builds to keep images minimal and secure.
  - *Backend (Go)*: ~29MB image running on Alpine runtime.
  
  - *Frontend (Next.js)*: Leveraged Next.js standalone output mode to prune build-time node_modules, reducing image size from ~1GB to ~170MB.

- **Security Hardening**: Enforced non-root system users (nextjs:1001) inside production runtime containers to prevent privilege escalation vulnerabilities.

### 2. CI/CD Pipeline (GitHub Actions)

- **Automated Builds**: Listens on push triggers to main and manually via workflow_dispatch.

- **Layer Caching**: Integrates GitHub Actions Cache (type=gha) with docker/setup-buildx-action to drastically accelerate image build execution time.

- **Deterministic Versioning**: Tags Docker images using immutable GitHub commit SHAs (${{ github.sha }}) alongside latest tags.

### 3. Kubernetes Orchestration Strategy

- **Decoupled Architecture**: Frontend and Backend deployed as separate stateless Deployments with independent replica counts and resource constraints.

- **Path-Based Ingress Routing**: Single Nginx Ingress routes / traffic to Next.js and /api requests to Go API, eliminating cross-origin (CORS) overhead.

- **Health Probes & Liveness Checks**: Exposes /health endpoints on the backend API to enable Kubernetes livenessProbe and readinessProbe automated container restarts.🔄 CI/CD Automation FlowPlaintext
```
[ Git Push (main) ]

         │
         ▼
 ┌────────────────────────────────────────────────────────┐
 │                GitHub Actions Pipeline                 │
 ├────────────────────────────────────────────────────────┤
 │  1. Checkout & Setup Docker Buildx                     │
 │  2. Authenticate to Docker Hub Registry                │
 │  3. Build & Push Backend Container (Cache-Enabled)     │
 │  4. Build & Push Frontend Container (Standalone Output)│
 └──────────────────────────┬─────────────────────────────┘
                            │
                            ▼
         [ Published Images on Docker Hub ]
                            │
                            ▼
         [ Kubernetes Rolling Update Deployment ]

```
## Tech Stack & Tooling

| Domain | Technology |
|:---------|:--------:|
| Backend API |Golang, Echo Framework, GORM |
| Frontend UI|Next.js (App Router), Tailwind CSS, Framer Motion |
| Database | PostgreSQL |
| Containerization| Docker, Docker Compose, Multi-Stage Builds |
| CI/CD|GitHub Actions, GitHub Container Registry / Docker Hub |
| Orchestration|Kubernetes, Nginx Ingress Controller |


## GitOps, Observability & the Ops Canvas

```text
 git push ──► GitHub Actions ──► Docker Hub (image:<sha>)
                    │
                    └─► commits tag into charts/clicker/values.yaml
                                   │
                         Argo CD (polls main every 60s)
                                   │  helm template + sync
                                   ▼
   ┌──────────────────── kind cluster ─────────────────────┐
   │ clicker-app: frontend · backend · postgres         │
   │              ops canvas (dashboard.localhost)           │
   │ monitoring : Prometheus · Grafana · kube-state-metrics  │
   │ argocd     : Argo CD                                    │
   └────────────────────────────────────────────────────────┘
```

| Piece | Where |
|---|---|
| App Helm chart (frontend, backend, Postgres, dashboard, ingress, ServiceMonitor) | `charts/clicker/` |
| Argo CD Application (auto-sync, prune, self-heal) | `argocd/clicker-app.yaml` |
| Platform Helm values + installer | `infra/` |
| Backend metrics (`/metrics`): request rate/latency, clicks | `backend/metrics.go` |

### Ops Canvas (dashboard.localhost)
A Figma-style infinite canvas with a dotted grid. Every component of the system is a live card:
commits and Actions runs from GitHub, Argo CD sync/health and deployed image tags, every pod
(ready state, CPU, memory, restarts, req/s), p95 latency sparklines, Prometheus
target health. Edges animate when traffic actually flows. Pan by dragging,
zoom with the wheel, drag cards by their header (layout is remembered), press `F` to fit.
A "Traffic generator" card fires real requests through the stack so you can watch it light up.

## Run it locally

Prerequisites: Docker Desktop, `kind`, `kubectl`, `helm`.

```bash
kind create cluster --name clicker-cluster --config kind-config.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
bash infra/install.sh   # monitoring, Argo CD, then hands the app to Argo CD
```

| UI | URL | Login |
|---|---|---|
| App | http://localhost | — |
| Ops canvas | http://dashboard.localhost | — |
| Argo CD | http://argocd.localhost | admin / admin |
| Grafana | http://grafana.localhost | admin / admin |
| Prometheus | http://prometheus.localhost | — |

`*.localhost` resolves to 127.0.0.1 in Chrome/Edge/Firefox without touching the hosts file.
Default credentials are for a local demo only.
