#!/usr/bin/env bash
# Installs the platform layer (monitoring, Argo CD) and hands the app to Argo CD.
# Re-runnable: every step is `helm upgrade --install` / `kubectl apply`.
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

# Monitoring first: it ships the ServiceMonitor CRD the other charts use.
helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
  --version 92.10.1 -n monitoring --create-namespace -f monitoring-values.yaml --wait --timeout 10m

helm upgrade --install argocd argo/argo-cd \
  --version 10.10.5 -n argocd --create-namespace -f argocd-values.yaml --wait --timeout 10m

kubectl apply -f ../argocd/clicker-app.yaml
