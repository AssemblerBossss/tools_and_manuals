# Proxmox

Команды ниже выполняются на хосте proxmox


## master-node
```bash
qm create 901 --name kube-master --memory 2048 --net0 virtio,bridge=vmbr0 --cores 2 --sockets 1 --ostype l26
qm importdisk 901  jammy-server-cloudimg-amd64.img.2 local-lvm
qm set 901 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-901-disk-0
qm set 901 --ide2 local-lvm:cloudinit
qm set 901 --boot c --bootdisk scsi0
qm set 901 --serial0 socket --vga serial0
qm set 901 --cipassword "yourpassword"  # задать пароль
qm set 901 --sshkey ~/.ssh/id_rsa.pub   # добавить SSH ключ
```


## kube-worker1
```bash
# kube-worker-1
qm create 902 --name kube-worker-1 --memory 2048 --net0 virtio,bridge=vmbr0 --cores 2 --sockets 1 --ostype l26
qm importdisk 902 jammy-server-cloudimg-amd64.img.2 local-lvm
qm set 902 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-902-disk-0
qm set 902 --ide2 local-lvm:cloudinit
qm set 902 --boot c --bootdisk scsi0
qm set 902 --serial0 socket --vga serial0
qm set 902 --cipassword "yourpassword"
qm set 902 --sshkey ~/.ssh/id_rsa.pub
qm start 902
```

## kube-worker1

```bash
# kube-worker-2
qm create 903 --name kube-worker-2 --memory 2048 --net0 virtio,bridge=vmbr0 --cores 2 --sockets 1 --ostype l26
qm importdisk 903 jammy-server-cloudimg-amd64.img.2 local-lvm
qm set 903 --scsihw virtio-scsi-pci --scsi0 local-lvm:vm-903-disk-0
qm set 903 --ide2 local-lvm:cloudinit
qm set 903 --boot c --bootdisk scsi0
qm set 903 --serial0 socket --vga serial0
qm set 903 --cipassword "yourpassword"
qm set 903 --sshkey ~/.ssh/id_rsa.pub
qm start 903
```

## 1. Подготовка всех узлов

1. Отключить управление сетью cloud-init
```bash
echo "network: {config: disabled}" | sudo tee /etc/cloud/cloud.cfg.d/99-disable-network-config.cfg
```
2. Создать свой netplan-файл `/etc/netplan/01-netcfg.yaml`

```yaml        
network:
  version: 2
  renderer: networkd
  ethernets:
    ens18:
      dhcp4: false
      addresses:
      # 101-master, 102-worker-1, 103-worker-2
        - 10.0.2.101/24 
      gateway4: 10.0.2.2
      nameservers:
        addresses: [8.8.8.8, 1.1.1.1]
```

3. Применить (статические IP-адреса сохранятся даже после перезагрузки)
```bash
sudo netplan apply
```
4. По базе

```bash
sudo apt update && sudo apt upgrade -y
```

5. Устанавливаем нужные пакеты
```bash   
sudo apt install -y apt-transport-https ca-certificates curl
```

6. Отключаем swap (Kubernetes требует)
   
```bash
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab
```

7. Включаем необходимые модули ядра
   
```bash
sudo modprobe overlay
sudo modprobe br_netfilter
```        

8. Настройка sysctl для Kubernetes

```bash
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system
```


## 2. Установка Docker (или containerd)


```bash
sudo apt install -y docker.io
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER
```

> Kubernetes больше не работает напрямую с Docker — нужен адаптер cri-dockerd (он делает Docker совместимым с Kubernetes CRI API).

```bash
wget https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.14/cri-dockerd_0.3.14.3-0.ubuntu-jammy_amd64.deb
sudo apt install -y ./cri-dockerd_0.3.14.3-0.ubuntu-jammy_amd64.deb
```

Запускаем и проверяем сервис

```bash 
sudo systemctl enable cri-docker.service
sudo systemctl enable --now cri-docker.socket
sudo systemctl start cri-docker.service
sudo systemctl status cri-docker.service
```

## 3. Установка kubeadm, kubelet и kubectl

These instructions are for Kubernetes v1.33.

1. Update the apt package index and install packages needed to use the Kubernetes apt repository:

```bash
sudo apt-get update
# apt-transport-https may be a dummy package; if so, you can skip that package
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
```

2. Download the public signing key for the Kubernetes package repositories. The same signing key is used for all repositories so you can disregard the version in the URL:

```bash
# If the directory `/etc/apt/keyrings` does not exist, it should be created before the curl command, read the note below.
# sudo mkdir -p -m 755 /etc/apt/keyrings
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.33/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
```

> Note:
In releases older than Debian 12 and Ubuntu 22.04, directory /etc/apt/keyrings does not exist by default, and it should be created before the curl command.


3. Add the appropriate Kubernetes apt repository. Please note that this repository have packages only for Kubernetes 1.33; for other Kubernetes minor versions, you need to change the Kubernetes minor version in the URL to match your desired minor version (you should also check that you are reading the documentation for the version of Kubernetes that you plan to install).

```bash
# This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.33/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list
```

4. Update the apt package index, install kubelet, kubeadm and kubectl, and pin their version:

```bash
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
(Optional) Enable the kubelet service before running kubeadm:

sudo systemctl enable --now kubelet
```

## На мастере

Инициализируем кластер с `Docker runtime`:

```bash
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --cri-socket unix:///var/run/cri-dockerd.sock
```

```bash
mkdir -p $HOME/.kube
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

## На каждом воркере

```bash
sudo kubeadm join 10.0.2.101:6443 --token <TOKEN> --discovery-token-ca-cert-hash sha256:<HASH>

```

## Скрипт `reset.sh`
```bash
#!/bin/bash
set -euo pipefail

echo "[INFO] Starting Kubernetes node cleanup on $(hostname)..."

# 1. Сброс конфигурации кластера
echo "[INFO] Running kubeadm reset..."
sudo kubeadm reset -f || true

# 2. Остановка служб
echo "[INFO] Stopping kubelet and container runtime..."
sudo systemctl stop kubelet 2>/dev/null || true
sudo systemctl stop containerd 2>/dev/null || true
sudo systemctl stop docker 2>/dev/null || true

# 3. Размонтирование всех kubelet-томов
echo "[INFO] Unmounting leftover kubelet mounts..."
sudo umount $(mount | grep '/var/lib/kubelet' | awk '{print $3}') 2>/dev/null || true

# 4. Очистка каталогов Kubernetes и CNI
echo "[INFO] Removing Kubernetes and CNI directories..."
sudo rm -rf /etc/kubernetes \
             /var/lib/etcd \
             /var/lib/kubelet/* \
             /var/lib/kube-proxy \
             /etc/cni/net.d \
             /var/lib/cni \
             /opt/cni \
             /var/run/kubernetes \
             /var/lib/dockershim 2>/dev/null || true

# 5. Очистка сетевых интерфейсов CNI
echo "[INFO] Cleaning up network interfaces..."
for iface in cni0 flannel.1 weave br0 docker0; do
    sudo ip link delete $iface 2>/dev/null || true
done

# 6. Очистка iptables
echo "[INFO] Flushing iptables rules..."
sudo iptables -F || true
sudo iptables -t nat -F || true
sudo iptables -t mangle -F || true
sudo iptables -X || true

# 7. Удаление kubeconfig пользователя
echo "[INFO] Removing kubeconfig for current user..."
rm -rf $HOME/.kube 2>/dev/null || true

# 8. Перезапуск системных демонов
echo "[INFO] Reloading systemd..."
sudo systemctl daemon-reload
sudo systemctl restart containerd 2>/dev/null || true
sudo systemctl restart docker 2>/dev/null || true

echo "[OK] Cleanup completed on $(hostname)"
```