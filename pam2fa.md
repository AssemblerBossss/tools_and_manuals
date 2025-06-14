# Лабораторная работа на PAM

## Часть 1: Создание базового образа Ubuntu вручную

1. Установим необходимые утилиты

```bash
sudo apt install debootstrap schroot
```

2. Создадим директорию для rootfs
```bash
sudo mkdir -p /srv/my-ubuntu
```

3. Запустим debootstrap
```bash
sudo debootstrap --variant=minbase jammy /srv/my-ubuntu http://archive.ubuntu.com/ubuntu
```

`minbase` - минимальный набор необходимых утилит

`jammy` - Ubuntu 22.04 

4. Подключим необходимые файловые системы (если работаем не в докере, а на хосте)
```bash
sudo mount --bind /dev /srv/my-ubuntu/dev
sudo mount --bind /proc /srv/my-ubuntu/proc
sudo mount --bind /sys /srv/my-ubuntu/sys
```

5. Войдем в chroot
```bash
sudo chroot /srv/my-ubuntu /bin/bash
```

6. Для наших пакетов добавим секции
```bash
echo 'deb http://archive.ubuntu.com/ubuntu jammy main universe restricted multiverse' > /etc/apt/sources.list
echo 'deb http://archive.ubuntu.com/ubuntu jammy-updates main universe restricted multiverse' >> /etc/apt/sources.list
echo 'deb http://security.ubuntu.com/ubuntu jammy-security main universe restricted multiverse' >> /etc/apt/sources.list
```
и обновим списки пакетов
```
apt update
```

7. Установим необходимые для нас пакеты
```bash
apt install nano sudo vim curl net-tools iproute2 systemd ssh otpw-bin libpam-otpw libpam-oath oathtool
```

8. Создадим пользователя и зададим ему пароль
```bash
useradd -m -s /bin/bash user
passwd user
```

9. Выйдем и размонтируем
```bash
exit
sudo umount /srv/my-ubuntu/{dev,proc,sys}
```

10. Упакуем `chroot`
```bash
cd /srv/my-ubuntu
sudo tar --numeric-owner -c . | gzip > /tmp/my-ubuntu.tar.gz
```

11. Импортируем образ в Docker
```bash
cat /tmp/my-ubuntu.tar.gz | docker import - my-ubuntu:jammy
```

12. Запустим Docker контейнер
```bash
sudo docker run -it --hostname=pam2fa --name=pam-2fa-container my-ubuntu:jammy bash
```

```bash
sudo docker run -d -p 2222:22 --hostname=pam2fa --name=pam-2fa my-ubuntu:jammy bash -c "mkdir -p /run/sshd && /usr/sbin/sshd -D"
```

## Часть 2: Выполнение лабораторной работы

1.	Установить программу otpw-bin  и библиотеку libpam-otpw
```bash
sudo apt install otpw-bin libpam-otpw
```
Можно не устанавливать, так как они уже были установлены при создании контейнера

2.	Сгенерировать достаточное число одноразовых паролей для пользователя user
```bash
sudo su - user
otpw-gen
```

3.	Настроить PAM таким образом, чтобы при аутентификации в рамках программы sudo аутентификация велась только с использованием одноразового пароля
```bash
sudo nano /etc/pam.d/sudo
```
Удаляем в нем или комментируем все содержимое и добавляем строчку
```
auth       required   pam_otpw.so
```

4.	Проверить корректность аутентификации
Сбрасываем кеш `sudo`
```bash
sudo -k
```
и проверяем
```bash
sudo ls
```
В появившемся окне вводим префикс и две части по 4 символа без пробела пароля по списку

5.	Настроить PAM таким образом, чтобы при аутентификации в рамках программы sudo аутентификация велась сначала по обычному паролю пользователя, а потом  с использованием одноразового пароля
Опять редактируем файл
```bash
sudo nano /etc/pam.d/sudo
```
и пишем в нем 
```
auth       required   pam_unix.so
auth       required   pam_otpw.so
```

6.	Проверить корректность аутентификации
Сбрасываем кеш `sudo`
```bash
sudo -k
```
и проверяем
```bash
sudo ls
```

7.	 Вернуть настройку Pam в исходное состояние
Нужно удалить добавленные записи и раскоментировать исходное содержание файла, либо восстановить его из копии

8.	Установить программу oathtool и библиотеку libpam-oath
```bash
sudo apt install oathtool libpam-oath
```
Можно не устанавливать, в контейнере они есть

9.	Используя программу oathtool  показать пароли привязанные ко времени
```bash
oathtool --totp -b "JBSWY3DPEHPK3PXP"
```

10.	Настроить Pam на аутентификацию  программы sudo  через libpam-oath с паролями привязанными ко времени

Создаем файл
```bash
sudo nano /etc/users.oath
```
заполняем следующим
```
HOTP/T30/6 user - JBSWY3DPEHPK3PXP
```
устанавливаем на него права
```
sudo chown root:root /etc/users.oath
sudo chmod 600 /etc/users.oath
```
настраиваем PAM для sudo
```bash
sudo nano /etc/pam.d/sudo
```
добавляем в начало
```
auth required pam_oath.so usersfile=/etc/users.oath window=30 digits=6
```
Объяснение:

`HOTP` — основа алгоритма (в данном случае TOTP, который на базе HOTP) HMAC-Based One-Time Password

`T30` — шаг времени 30 секунд

`6` — длина одноразового кода

`user` — имя пользователя

`-` — поле счетчика (актуально для HOTP). Для TOTP остается пустым (-)

`JBSWY...` — секрет в base32, используемый для генерации пароля

файл /etc/users.oath не изменяется, потому что используется алгоритм TOTP и каждые 30 секунд генерируется новый пароль. Файл будет изменяться при использовании алгоритма HOTP, т.к. он использует счетчик для генерации паролей.

Для проверки генерируем пароль
```bash
oathtool --totp -b JBSWY3DPEHPK3PXP
```
и проверяем
```bash
sudo -k
sudo ls
```

12.	Установить сервер SSH и настроить для него двух факторную аутентификацию с одноразовым паролем в качестве дополнительного фактора аутентификации. Вход пользователей без одноразовых паролей должен быть запрещен

Настраиваем PAM для SSH
```bash
sudo nano /etc/pam.d/sshd
```
и добавляем в начало строку
```
auth required pam_oath.so usersfile=/etc/users.oath window=30 digits=6
```
изменяем файл
```bash
sudo nano /etc/ssh/sshd_config
```
и добавляем/исправляем строки
```
KbdInteractiveAuthentication yes
ChallengeResponseAuthentication yes
UsePAM yes
AuthenticationMethods keyboard-interactive
```
перезагружаем ssh
```bash
sudo kill -HUP $(pidof sshd)
```
далее в другом терминале с хоста подключаемся к нашему контейнеру по ssh
```bash
ssh user@0.0.0.0 -p 2222
```
у нас запросят OTP пароль, после генерируем пароль в самом контейнере, вводим его на хосте, вводим пароль от пользователя и подключаемся
```bash
oathtool --totp -b JBSWY3DPEHPK3PXP
```