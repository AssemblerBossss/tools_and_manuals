1. Добавить сетевой интерфейс Host Only Adapter


iface enp0s8 inet manual

auto vmbr1
iface vmbr1 inet static
    address 192.168.56.10X/24
    bridge-ports enp0s8
    bridge-stp off
    bridge-fd 0



## Настройка PROXMOX для общей сети
```bash
vboxmanage modifyvm "Proxmox-VE" --nic2 hostonly --hostonlyadapter2 vboxnet0 --nicpromisc2 allow-all --cableconnected2 on
```
Эта команда настраивает виртуальную машину “Proxmox-VE” в VirtualBox, изменяя параметры её второго сетевого адаптера (NIC2). Вкратце она делает следующее:

- Устанавливает тип сети NIC2 как Host-Only — VM получает сетевое соединение только с хостом и другими VM в этой сети.
- Назначает адаптеру Host-Only сеть vboxnet0.
- Включает режим promisc (allow-all) — адаптер сможет принимать весь сетевой трафик в сети, что полезно для бриджинга или мониторинга.
- Включает "подключён кабель" — адаптер считается активным

Настройка времени

# На ВСЕХ узлах выполнить:
apt update && apt install chrony -y

# Редактировать конфиг chrony
nano /etc/chrony/chrony.conf

закомментировать строчку с pool

server 192.168.56.101 iburst

systemctl restart chrony