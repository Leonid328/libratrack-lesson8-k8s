# Theme 8: Docker, Compose, K3s

## VM (70 pts): Docker + Floci + K3s + demos 1 and 2

On a Linux VM (Ubuntu/Debian):

```bash
# Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker "$USER"
newgrp docker

# K3s (single-node)
curl -sfL https://get.k3s.io | sh -
sudo kubectl get nodes

# Floci — run as in the course presentation (same VM, Docker already up)
```

Presentation demos (app 1 and 2):

```bash
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
kubectl apply -f deploy/k3s-demos/app1-nginx.yaml
kubectl apply -f deploy/k3s-demos/app2-whoami.yaml
kubectl get pods,svc
curl -s http://127.0.0.1:30081/ | head
curl -s http://127.0.0.1:30082/
```

Screenshots: `kubectl get nodes`, both demos running, curl output.

## Docker Compose (+30 pts)

From repo root:

```bash
docker compose up --build -d
curl -s http://localhost:8080/health
curl -s http://localhost:8080/api/v1/books | head
open http://localhost:8080/catalog
```

## Kubernetes (LibraTrack)

Build image on the VM and load into K3s:

```bash
docker build -t libratrack:latest .
docker save libratrack:latest | sudo k3s ctr images import -
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/postgres.yaml
kubectl wait -n libratrack --for=condition=ready pod -l app=postgres --timeout=120s
kubectl apply -f deploy/k8s/api.yaml
kubectl get all -n libratrack
curl -s http://127.0.0.1:30080/health
```

## GitHub for LMS

Push this folder (or monorepo path `homework/libratrack`) and submit:

- link to `Dockerfile`, `docker-compose.yml`
- link to `deploy/k8s/` and `deploy/k3s-demos/`
- video or screenshots: Compose health + catalog, K3s `kubectl get pods`, demo 1/2
