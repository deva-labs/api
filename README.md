# 🛠 DevOps SaaS Platform - Database Schema

A system that helps users create software projects with CI/CD, deploy YAML (Kubernetes/Docker), import from GitHub, and manage user permissions based on subscription plans.

---

## 🧩 Feature Overview

- Create projects from templates or GitHub repositories.
- Generate files like `Dockerfile`, `docker-compose.yml`, `deployment.yaml`, CI pipelines, etc.
- Deploy applications to Kubernetes or Docker hosts.
- User permissions based on subscription plans: Free / Pro / Enterprise.
- Support for teams, billing, secrets, webhooks, and usage tracking.

---

## 📊 Database Schema (20 Tables)

### 1. `users`
Registered user information.
- `id`, `email`, `name`, `role`, `plan_id`, `created_at`

---

### 2. `plans`
Service plan definitions.
- `id`, `name`, `price`, `max_projects`, `max_deployments_per_day`

---

### 3. `billing_subscriptions`
Paid subscription status.
- `user_id`, `provider`, `subscription_id`, `status`, `started_at`, `renew_at`

---

### 4. `teams`
Teams for collaboration.
- `id`, `name`, `owner_id`, `created_at`

---

### 5. `team_members`
Assign users to teams with roles.
- `team_id`, `user_id`, `role`, `joined_at`

---

### 6. `projects`
User/team software projects.
- `team_id`, `name`, `repo_url`, `source_type`, `status`

---

### 7. `project_templates`
Templates for creating sample projects.
- `name`, `description`, `language`, `is_official`

---

### 8. `project_configs`
Environment configurations for generating deployment files.
- `project_id`, `language`, `framework`, `port`, `env_vars`, `ci_tool`, `deploy_target`

---

### 9. `project_files`
Stored files like Dockerfile, deployment.yaml, etc.
- `project_id`, `path`, `content`, `is_generated`

---

### 10. `ci_pipelines`
CI/CD pipeline status for projects.
- `project_id`, `provider`, `status`, `commit_sha`, `log_url`

---

### 11. `pipeline_steps`
Detailed steps in each pipeline.
- `pipeline_id`, `name`, `status`, `log`

---

### 12. `deployments`
Track application deployments.
- `project_id`, `platform`, `status`, `log`, `triggered_by`

---

### 13. `deployment_targets`
Define deployment targets (Kubernetes, Docker hosts).
- `team_id`, `name`, `type`, `host`, `auth`

---

### 14. `github_integrations`
Store GitHub OAuth tokens for repo imports and CI.
- `user_id`, `github_user_id`, `access_token`

---

### 15. `activity_logs`
Record actions like project creation, deployments, etc.
- `user_id`, `project_id`, `action`, `metadata`

---

### 16. `api_keys`
API tokens for user automation.
- `user_id`, `name`, `token_hash`, `scopes`

---

### 17. `notifications`
CI/CD status, billing, deployment notifications.
- `user_id`, `type`, `message`, `is_read`

---

### 18. `webhooks`
User-registered webhooks for events.
- `project_id`, `url`, `event_type`, `secret`

---

### 19. `secrets`
Environment variables and sensitive project tokens.
- `project_id`, `key`, `value_encrypted`, `scope`

---

### 20. `usage_metrics`
Track build/deploy counts for plan usage calculations.
- `user_id`, `project_id`, `metric_type`, `count`, `period_start`, `period_end`

---

## 🧠 Extensions

- Additional support: Google OAuth, GitLab, Bitbucket.
- Additional tables: `audit_logs`, `cluster_metrics`, `custom_plugins`.
- Public API gateway with `api_keys` table.

---

## 📌 Contact

Contact [yourteam@yourdomain.com] for feedback or architecture expansion requests.

# Docker Remote API Setup Guide

This guide explains how to configure the Docker Remote API with TLS on a Linux server. The setup uses OpenSSL to generate certificates for secure communication between the client and Docker daemon.

## Prerequisites

- Linux system with Docker installed (tested on Docker 28.1.1)
- Root/sudo access
- OpenSSL installed

## Step 1: Install Docker

### For Ubuntu/Debian:

```bash
# 1. Remove any old Docker versions
sudo apt-get remove docker docker-runtime docker.io containerd runc

# 2. Install prerequisites
sudo apt-get update
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# 3. Add Docker's official GPG key
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# 4. Set up the repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 5. Install Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 6. Verify installation
sudo docker run hello-world
```

### For CentOS/RHEL:

```bash
sudo yum install docker-ce docker-ce-cli containerd.io
```

### Start and Enable Docker:

```bash
sudo systemctl enable --now docker
```

## Step 2: Configure Docker Remote API with TLS

### 2.1 Create Certificates Directory

```bash
sudo mkdir -p /etc/docker/certs
cd /etc/docker/certs
```

### 2.2 Generate CA Certificate

```bash
sudo openssl genrsa -out ca-key.pem 4096
sudo openssl req -new -x509 -days 3650 -key ca-key.pem -sha256 -out ca.pem -subj "/CN=docker-ca"
```

### 2.3 Generate Server Certificate with Proper SANs

- Create OpenSSL config file for server certificate

```bash
# First get your hostname
HOSTNAME=$(hostname)

# Then create the openssl.cnf file with the actual hostname
sudo tee openssl.cnf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name

[req_distinguished_name]

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = $HOSTNAME
DNS.2 = localhost
IP.1 = 127.0.0.1
IP.2 = 192.168.237.116
EOF
```

- Generate server key and CSR

```bash
sudo openssl genrsa -out server-key.pem 4096
sudo openssl req -new -key server-key.pem -out server.csr -subj "/CN=$(hostname)" -config openssl.cnf
```

- Sign the certificate

```bash
sudo openssl x509 -req -days 3650 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extensions v3_req -extfile openssl.cnf
```

### 2.4 Generate Client Certificate

```bash
sudo openssl genrsa -out key.pem 4096
sudo openssl req -new -key key.pem -out client.csr -subj "/CN=client"
sudo openssl x509 -req -days 3650 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem
```

### 2.5 Set Proper Permissions

```bash
sudo chmod 644 /etc/docker/certs/*.pem
sudo chmod 600 /etc/docker/certs/*-key.pem
```

## Step 3: Configure Docker Daemon

### 3.1 Create `daemon.json`

```bash
sudo tee /etc/docker/daemon.json <<EOF
{
  "hosts": ["tcp://0.0.0.0:2376", "unix:///var/run/docker.sock"],
  "tlsverify": true,
  "tlscacert": "/etc/docker/certs/ca.pem",
  "tlscert": "/etc/docker/certs/server-cert.pem",
  "tlskey": "/etc/docker/certs/server-key.pem",
  "features": {
    "buildkit": true
  }
}
EOF
```

### 3.2 Create Systemd Override

```bash
sudo mkdir -p /etc/systemd/system/docker.services.d
sudo tee /etc/systemd/system/docker.services.d/override.conf <<EOF
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd
EOF
```

### 3.3 Reload and Restart Docker

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

## Step 4: Verify the Setup

### 4.1 Check Docker Status

```bash
sudo systemctl status docker
```

### 4.2 Test Connection

```bash
curl --cert /etc/docker/certs/cert.pem      --key /etc/docker/certs/key.pem      --cacert /etc/docker/certs/ca.pem      https://$(hostname):2376/version
```

## Step 5: Using Remote API from Client

### 5.1 Copy These Files from Server to Client

- `ca.pem`
- `cert.pem`
- `key.pem`
```bash
scp ties@192.168.237.116:/etc/docker/certs/{ca.pem,cert.pem,key.pem} "/mnt/e/Source Code/dockerwizard/api/store/secrets/"
chmod 644 "/mnt/e/Source Code/dockerwizard/api/store/secrets/ca.pem"
chmod 644 "/mnt/e/Source Code/dockerwizard/api/store/secrets/cert.pem"
chmod 600 "/mnt/e/Source Code/dockerwizard/api/store/secrets/key.pem"
```
### 5.2 Set Environment Variables on Client

```bash
export DOCKER_HOST=tcp://<server-ip>:2376
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=/path/to/certs
```

### 5.3 Test from Client
Create Docker context with TLS
```bash
docker context create myremote \
  --docker "host=tcp://192.168.237.116:2376,ca=/app/store/secrets/ca.pem,cert=/app/store/secrets/cert.pem,key=/app/store/secrets/key.pem"
docker context use myremote

```

---

This setup ensures secure communication between Docker clients and the Docker daemon using TLS. Make sure to replace `<server-ip>` with the actual server IP in your environment.
