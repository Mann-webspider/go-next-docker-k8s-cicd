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


## Local Development Setup1. Run using Docker Compose

```Bash
# Clone the repository
git clone https://github.com/manndalsania/k8s-devops-pipeline.git
cd k8s-devops-pipeline

# Spin up PostgreSQL, Go Backend, and Next.js Frontend
docker-compose up --build -d

```
2. Deploy to Kubernetes Cluster (Minikube / EKS)Bash
```
# Apply ConfigMaps, Secrets, and Persistent Volumes
kubectl apply -f k8s/config/

# Deploy Services & Deployments
kubectl apply -f k8s/backend/
kubectl apply -f k8s/frontend/

# Apply Ingress Controller
kubectl apply -f k8s/ingress.yaml
```
---

### Why this README will make a great recruiter impression:
1. **Badges at the top:** Instantly signals your tech stack visually.
2. **ASCII Flowcharts:** Easy for hiring managers to scan in under 10 seconds.
3. **Quantifiable Metrics:** Emphasizing image sizes (**~1GB to 170MB** for Next.js, **~29MB** for Go) proves real-world Docker optimization skills.