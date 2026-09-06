# Deploying Terminal Portfolio (TUI) to cli.harshadnikam.me on Azure VM

A complete, production-ready guide to deploying your interactive Go + Bubble Tea + Wish terminal portfolio to **`cli.harshadnikam.me`** on an Azure Virtual Machine.

Visitors can access your interactive portfolio directly from any terminal:

```bash
ssh cli.harshadnikam.me
# or if running on custom port:
ssh -p 2222 cli.harshadnikam.me
```

---

## 1. Architecture Overview

Your portfolio application is composed of two components:

1. **CLI Binary (`cmd/cli`)**: Runs locally when launched via terminal (`./harshad-tui`).
2. **SSH Server Daemon (`cmd/ssh`)**: Built using Charm's [Wish](https://github.com/charmbracelet/wish) middleware on top of [Bubble Tea](https://github.com/charmbracelet/bubbletea). It listens for incoming SSH client connections, allocates a virtual PTY for each session, initializes the model, and cleanly terminates when the user presses `q`, `Esc`, `Ctrl+C`, or disconnects.

---

## 2. Azure VM Sizing & Budget ($100 Credit Plan)

To run a persistent SSH daemon, you need a Linux host with a dedicated public IPv4 address. Standard serverless platforms (Vercel, Netlify, Cloudflare Pages) cannot host long-running raw TCP/SSH sockets.

### Recommended Azure Specification:
* **VM Size:** `Standard_B1s` (1 vCPU, 1 GiB RAM) — approx. **$7.50 – $8.00 USD/month** (~₹650 INR).
* **OS:** `Ubuntu Server 24.04 LTS - x64 Gen2` (No OS licensing fees).
* **Disk:** `32 GB Standard SSD (locally-redundant storage)` (approx. **$2.40 USD/month**).
* **IP:** Static Public IPv4 (~**$2.90 USD/month**).
* **Total Estimated Cost:** **~$10.00 – $11.00 USD/month**.

> [!TIP]
> With your **$100 Azure Credit**, this VM configuration will run continuously for **9 to 10 months** without any out-of-pocket expense.

---

## 3. Step-by-Step Azure VM Provisioning

### Step 3.1: Basics Tab
1. Open the [Azure Portal](https://portal.azure.com/).
2. Search for **Virtual machines** and click **Create** > **Azure virtual machine**.
3. Fill in the parameters:
   * **Subscription:** Select the subscription with your credit.
   * **Resource group:** Click *Create new* -> `rg-harshad-tui`.
   * **Virtual machine name:** `vm-harshad-tui`.
   * **Region:** `Central India` (or closest to your primary audience; `East US` is cheapest).
   * **Availability options:** *No infrastructure redundancy required*.
   * **Security type:** *Standard*.
   * **Image:** `Ubuntu Server 24.04 LTS - x64 Gen2`.
   * **VM architecture:** `x64`.
   * **Size:** `Standard_B1s` (1 vcpu, 1 GiB memory).
   * **Authentication type:** `SSH public key`.
   * **Username:** `azureuser`.
   * **SSH public key source:** `Generate new key pair` (Key name: `vm-harshad-tui_key`).
   * **Public inbound ports:** Select `Allow selected ports` -> choose `SSH (22)`.

### Step 3.2: Disks Tab
* Change **OS disk type** from *Premium SSD* to **Standard SSD (locally-redundant storage)** to maximize credit lifespan.
* Check **Delete with VM**.

### Step 3.3: Networking Tab
* **Virtual network:** Allow default (`rg-harshad-tui-vnet`).
* **Subnet:** Default (`10.0.0.0/24`).
* **Public IP:** Click *Create new* -> ensure assignment is set to **Static** so the IP never changes on reboot.
* **NIC network security group:** Choose **Basic** (we configure rules next).
* Check **Delete public IP and NIC when VM is deleted**.

### Step 3.4: Review + Create
1. Click **Review + create** and verify the configuration.
2. Click **Create** and download the private key file (`vm-harshad-tui_key.pem`).
3. Save the key locally in `~/.ssh/vm-harshad-tui_key.pem`.

---

## 4. SSH Port Strategy & Azure NSG Configuration

Every Linux server runs OpenSSH for administrative login. You have two port routing strategies for the TUI:

### Strategy Comparison:

| Feature | Option A: TUI on Port 2222 | Option B (Recommended): TUI on Port 22 |
| :--- | :--- | :--- |
| **Visitor Command** | `ssh -p 2222 cli.harshadnikam.me` | `ssh cli.harshadnikam.me` (Cleanest!) |
| **Admin Login** | `ssh azureuser@<IP>` (Port 22) | `ssh -p 22022 azureuser@<IP>` (Port 22022) |
| **Setup Complexity** | Simplest (no OpenSSH config changes) | Requires changing admin SSH port to 22022 |

---

### Configuring Azure Network Security Group (NSG) Rules

In Azure Portal, navigate to **vm-harshad-tui** -> **Networking** (under Settings) -> **Inbound port rules**:

#### If Using Option B (TUI on Port 22 — Cleanest):
1. Click **Add inbound port rule** for admin SSH:
   * **Destination port ranges:** `22022`
   * **Protocol:** `TCP`
   * **Action:** `Allow`
   * **Priority:** `1010`
   * **Name:** `Allow_Admin_SSH_22022`
2. Ensure the existing rule for Port `22` is allowed with destination port `22` and priority `1020` (Name: `Allow_TUI_Port_22`).

#### If Using Option A (TUI on Port 2222):
1. Keep the default Port `22` rule for administrative SSH.
2. Click **Add inbound port rule**:
   * **Destination port ranges:** `2222`
   * **Protocol:** `TCP`
   * **Action:** `Allow`
   * **Priority:** `1010`
   * **Name:** `Allow_TUI_Port_2222`

---

## 5. Domain & DNS Setup (`cli.harshadnikam.me`)

Point your subdomain to your Azure VM's Static Public IPv4 address.

### Step 5.1: Retrieve Azure Public IP
In the Azure Portal, open **vm-harshad-tui** -> on the **Overview** page, copy the **Public IP address** (e.g. `20.198.xxx.xxx`).

### Step 5.2: Create DNS Record in Cloudflare / DNS Provider
Log in to your DNS provider for `harshadnikam.me`:

| Type | Name / Host | Content / IPv4 Value | TTL | Proxy Status (Cloudflare) |
| :--- | :--- | :--- | :--- | :--- |
| **A** | `cli` | `<YOUR_AZURE_PUBLIC_IP>` | Auto / 300s | **DNS only (Grey Cloud)** ⚠️ |

> [!CAUTION]
> **CRITICAL FOR CLOUDFLARE USERS:**
> Cloudflare's standard CDN/proxy (Orange Cloud) only routes HTTP/HTTPS (ports 80 & 443).
> Because SSH uses raw TCP traffic on port 22 or 2222, you **MUST set the record to DNS Only (Grey Cloud)**. Otherwise, client SSH connections will hang or fail.

### Step 5.3: Verify DNS Propagation
From your local terminal:
```bash
dig +short cli.harshadnikam.me
# Should return your Azure VM public IPv4 address

ping -c 1 cli.harshadnikam.me
```

---

## 6. Initial Server Setup & Hardening

On your local machine, set proper file permissions on your downloaded SSH key:
```bash
# Linux / macOS / WSL
chmod 400 ~/.ssh/vm-harshad-tui_key.pem
```

Connect to the Azure VM:
```bash
ssh -i ~/.ssh/vm-harshad-tui_key.pem azureuser@<AZURE_PUBLIC_IP>
```

Update packages and install core utilities:
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y ufw curl git build-essential libcap2-bin
```

---

## 7. Reconfiguring OpenSSH Port (Required for Option B)

To allow the portfolio TUI to bind directly to port 22 so visitors only need to run `ssh cli.harshadnikam.me`, move your administrative OpenSSH server to port **22022**:

1. Edit OpenSSH server configuration:
   ```bash
   sudo nano /etc/ssh/sshd_config
   ```

2. Locate the line `#Port 22` and edit it to:
   ```text
   Port 22022
   ```

3. Save and exit (`Ctrl+O`, `Enter`, `Ctrl+X`).

4. Restart the OpenSSH service:
   ```bash
   sudo systemctl restart ssh || sudo systemctl restart sshd
   ```

5. > [!IMPORTANT]
   > **DO NOT CLOSE YOUR CURRENT TERMINAL WINDOW!**
   > Open a second terminal window on your local machine and verify administrative login on port 22022:
   > ```bash
   > ssh -i ~/.ssh/vm-harshad-tui_key.pem -p 22022 azureuser@<AZURE_PUBLIC_IP>
   > ```
   > Only proceed once you have verified you can log in on the new port.

---

## 8. Firewall (UFW) Configuration

Configure Ubuntu's Uncomplicated Firewall (UFW):

```bash
# If using Option B (Admin on 22022, TUI on 22):
sudo ufw allow 22022/tcp comment "Admin OpenSSH"
sudo ufw allow 22/tcp comment "Portfolio TUI Wish"

# If using Option A (Admin on 22, TUI on 2222):
sudo ufw allow 22/tcp comment "Admin OpenSSH"
sudo ufw allow 2222/tcp comment "Portfolio TUI Wish"

# Enable the firewall
sudo ufw enable
sudo ufw status verbose
```

---

## 9. Deployment Method 1: Native Binary via Systemd (Recommended)

This method provides the lowest resource footprint (< 20MB RAM) and zero virtualization overhead.

### Step 9.1: Build Linux Binary on Your Local Machine
From the root of the project directory on your local machine:

```bash
# Cross-compile static Linux binary
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o harshad-ssh ./cmd/ssh
```

### Step 9.2: Upload Binary to Azure VM
```bash
# If using Option B (Admin port 22022):
scp -i ~/.ssh/vm-harshad-tui_key.pem -P 22022 harshad-ssh azureuser@<AZURE_PUBLIC_IP>:/home/azureuser/

# If using Option A (Admin port 22):
scp -i ~/.ssh/vm-harshad-tui_key.pem -P 22 harshad-ssh azureuser@<AZURE_PUBLIC_IP>:/home/azureuser/
```

### Step 9.3: Move to System Path & Set Capabilities
On the Azure VM:
```bash
sudo mv /home/azureuser/harshad-ssh /usr/local/bin/harshad-ssh
sudo chmod +x /usr/local/bin/harshad-ssh

# Allow non-root binary to bind to privileged ports (like port 22):
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/harshad-ssh
```

### Step 9.4: Create Unprivileged Service User & Directories
```bash
sudo useradd -r -s /bin/false -d /var/lib/harshad-tui -m harshad-tui
sudo mkdir -p /var/lib/harshad-tui/.ssh
sudo chown -R harshad-tui:harshad-tui /var/lib/harshad-tui
```

### Step 9.5: Create Systemd Service Unit
Create the service configuration:
```bash
sudo nano /etc/systemd/system/harshad-tui.service
```

Add the following configuration:

```ini
[Unit]
Description=Harshad Nikam SSH Terminal Portfolio Service
After=network.target

[Service]
Type=simple
User=harshad-tui
Group=harshad-tui
WorkingDirectory=/var/lib/harshad-tui
ExecStart=/usr/local/bin/harshad-ssh
Restart=always
RestartSec=3s
LimitNOFILE=65535

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/harshad-tui
PrivateTmp=true
AmbientCapabilities=CAP_NET_BIND_SERVICE

# Environment overrides
Environment=SSH_HOST=0.0.0.0
Environment=SSH_PORT=22
Environment=SSH_KEY_PATH=/var/lib/harshad-tui/.ssh/harshad_ed25519

[Install]
WantedBy=multi-user.target
```

*(Note: Change `Environment=SSH_PORT=22` to `2222` if you chose Option A).*

### Step 9.6: Enable and Start Service
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now harshad-tui
sudo systemctl status harshad-tui
```

---

## 10. Deployment Method 2: Docker & Docker Compose

If you prefer containerized deployment:

### Step 10.1: Install Docker on Azure VM
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker azureuser
# Log out and log back in to refresh group membership
```

### Step 10.2: Deploy with Docker Compose
Clone the repository and run:
```bash
git clone https://github.com/dev-harshhh19/Portfolio-TUI.git
cd Portfolio-TUI
docker compose -f deployments/docker-compose.yml up -d --build
```

---

## 11. Wish Production Security Checklist

The application includes the following production safeguards:

1. **Persistent ED25519 Host Key:** Host keys are saved to disk (`/var/lib/harshad-tui/.ssh/harshad_ed25519`). This prevents returning visitors from receiving frightening `WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED` alerts when the daemon restarts.
2. **PTY Session Enforcement:** All incoming sessions must request a PTY (`ssh -t` or standard interactive terminal). Raw non-interactive command injection is rejected.
3. **Automated Session Cleanup:** Idle timeouts (15 minutes) and max session timeouts (2 hours) automatically reclaim dead client sockets.
4. **Sandboxed VFS:** All `cd`, `cat`, and `ls` commands execute against the in-memory virtual filesystem with zero access to the host operating system.

---

## 12. Verification & Testing

### 12.1: Check Listening Sockets on the VM
```bash
sudo ss -tulpn | grep -E '22|22022|2222'
```
You should see:
* Port `22022` listened to by `sshd`
* Port `22` listened to by `harshad-ssh`

### 12.2: Test Connection from Your Local Computer
From your laptop/desktop terminal:

```bash
# For standard Option B:
ssh cli.harshadnikam.me

# Or if using custom port:
ssh -p 2222 cli.harshadnikam.me
```

### 12.3: Monitor Live Connection Logs
On the Azure VM:
```bash
journalctl -u harshad-tui -f
```

---

## 13. Showcase on GitHub and Website

Add the terminal command snippet to your web portfolio and GitHub profile README:

```markdown
### 🖥️ Interactive Terminal Portfolio

Explore my interactive terminal portfolio directly from your command line:

```bash
ssh cli.harshadnikam.me
```
*(No installation or dependencies required)*
```
