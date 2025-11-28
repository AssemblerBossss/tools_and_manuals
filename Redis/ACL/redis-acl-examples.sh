#!/bin/bash

# Подключение к Redis
# docker exec -it redis-server redis-cli

# ============================================
# СОЗДАНИЕ ПОЛЬЗОВАТЕЛЕЙ
# ============================================

# 1. Создать пользователя с полными правами
ACL SETUSER admin on >StrongPassword123 ~* &* +@all

# 2. Создать пользователя только для чтения
ACL SETUSER readonly on >ReadPass456 ~* &* +@read +@connection -@write -@dangerous

# 3. Создать пользователя для записи в кеш (только определенные команды)
ACL SETUSER cache_writer on >CachePass789 ~cache:* &* +get +set +del +expire +ttl

# 4. Создать пользователя для работы только с определенными ключами
ACL SETUSER app_user on >AppPass999 ~app:* ~session:* &* +@all

# 5. Создать пользователя для Pub/Sub
ACL SETUSER pubsub_user on >PubSubPass111 &* +@pubsub +@connection

# ============================================
# ПРОСМОТР ПОЛЬЗОВАТЕЛЕЙ
# ============================================

# Список всех пользователей
ACL LIST

# Детальная информация о пользователе
ACL GETUSER admin

# Показать текущего пользователя
ACL WHOAMI

# ============================================
# УПРАВЛЕНИЕ ПОЛЬЗОВАТЕЛЯМИ
# ============================================

# Отключить пользователя (но не удалять)
ACL SETUSER readonly off

# Включить обратно
ACL SETUSER readonly on

# Удалить пользователя
ACL DELUSER cache_writer

# Сохранить ACL в файл
ACL SAVE

# Загрузить ACL из файла
ACL LOAD

# ============================================
# ТЕСТИРОВАНИЕ ПРАВ
# ============================================

# Авторизоваться под пользователем
AUTH readonly ReadPass456

# Проверить, какие команды доступны
ACL CAT

# Проверить права на конкретную команду
ACL DRYRUN readonly SET key value
