
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
sudo apt-get remove docker docker-engine docker.io containerd runc

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
sudo mkdir -p /etc/systemd/system/docker.service.d
sudo tee /etc/systemd/system/docker.service.d/override.conf <<EOF
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
