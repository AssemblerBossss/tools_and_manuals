# Redis ACL (Access Control List) - Руководство

## 📋 Что такое ACL в Redis?

ACL (Access Control List) позволяет создавать пользователей с разными уровнями доступа:
- Разные пароли для разных пользователей
- Ограничение доступа к определенным ключам
- Ограничение доступа к определенным командам
- Ограничение Pub/Sub каналов

## 🚀 Быстрый старт

### 1. Запуск Redis с ACL

```bash
# Используем docker-compose с ACL конфигурацией
docker-compose -f docker-compose-acl.yml up -d

# Проверка логов
docker-compose -f docker-compose-acl.yml logs -f
```

### 2. Подключение под разными пользователями

```bash
# Admin пользователь (полные права)
docker exec -it redis-server redis-cli -u redis://admin:AdminStrongPass123!@localhost:6379

# Read-only пользователь
docker exec -it redis-server redis-cli --user readonly --pass ReadOnlyPass456!

# Cache writer
docker exec -it redis-server redis-cli --user cache_writer --pass CachePass789!
```

## 📝 Синтаксис ACL правил

### Формат команды:
```
ACL SETUSER <username> <flags> <passwords> <keys> <channels> <commands>
```

### Основные флаги:

| Флаг | Описание |
|------|----------|
| `on` | Включить пользователя |
| `off` | Отключить пользователя |
| `>password` | Установить пароль |
| `nopass` | Без пароля (небезопасно!) |
| `resetpass` | Удалить все пароли |

### Доступ к ключам:

| Правило | Описание | Пример |
|---------|----------|--------|
| `~*` | Все ключи | Доступ ко всему |
| `~cache:*` | Паттерн | Только ключи начинающиеся с `cache:` |
| `~app:* ~session:*` | Несколько паттернов | Доступ к `app:*` И `session:*` |
| `resetkeys` | Удалить все паттерны | |

### Доступ к Pub/Sub каналам:

| Правило | Описание |
|---------|----------|
| `&*` | Все каналы |
| `&channel:*` | Каналы по паттерну |
| `&events:* &logs:*` | Несколько паттернов |

### Доступ к командам:

#### Категории команд:
```bash
+@all          # Все команды
+@read         # Команды чтения (GET, KEYS, SCAN, etc)
+@write        # Команды записи (SET, DEL, INCR, etc)
+@admin        # Административные (CONFIG, SHUTDOWN, etc)
+@dangerous    # Опасные (FLUSHDB, FLUSHALL, etc)
+@connection   # Подключение (AUTH, PING, SELECT, etc)
+@pubsub       # Pub/Sub (PUBLISH, SUBSCRIBE, etc)
+@string       # Строковые команды
+@list         # Списки
+@set          # Множества
+@sortedset    # Отсортированные множества
+@hash         # Хеши
+@slow         # Медленные команды (KEYS, FLUSHDB, etc)
```

#### Отдельные команды:
```bash
+get           # Разрешить GET
+set           # Разрешить SET
-del           # Запретить DEL
+info          # Разрешить INFO
```

## 🎯 Примеры пользователей

### 1. Администратор (полный доступ)
```bash
ACL SETUSER admin on >StrongAdminPass123! ~* &* +@all
```

### 2. Read-Only пользователь
```bash
ACL SETUSER readonly on >ReadPass456! ~* &* +@read +@connection -@write -@dangerous
```
- Доступ ко всем ключам (`~*`)
- Все каналы Pub/Sub (`&*`)
- Только команды чтения
- Запрещены команды записи и опасные

### 3. Cache Worker (только кеш)
```bash
ACL SETUSER cache_worker on >CachePass789! ~cache:* &* +get +set +setex +del +expire +ttl
```
- Доступ только к ключам `cache:*`
- Только специфичные команды для кеша

### 4. Application User
```bash
ACL SETUSER app_user on >AppPass999! ~app:* ~session:* &* +@all -@admin -@dangerous
```
- Доступ к ключам `app:*` и `session:*`
- Все команды кроме административных

### 5. Monitoring User
```bash
ACL SETUSER monitor on >MonitorPass111! ~* &* +@read +@connection +info +ping +slowlog
```
- Чтение всех ключей
- Команды мониторинга (INFO, PING, SLOWLOG)

### 6. Queue Worker
```bash
ACL SETUSER queue_worker on >QueuePass333! ~queue:* ~job:* &* +@list +@connection +del +expire
```
- Работа с очередями (списки)
- Доступ только к `queue:*` и `job:*`

### 7. Pub/Sub User
```bash
ACL SETUSER pubsub_user on >PubSubPass222! &channel:* &events:* +@pubsub +@connection
```
- Только Pub/Sub команды
- Только каналы `channel:*` и `events:*`

## 🔧 Управление пользователями

### Просмотр пользователей
```bash
# Список всех пользователей с правами
ACL LIST

# Информация о конкретном пользователе
ACL GETUSER admin

# Текущий пользователь
ACL WHOAMI

# Список всех категорий команд
ACL CAT

# Команды в категории
ACL CAT read
```

### Изменение пользователя
```bash
# Изменить пароль
ACL SETUSER readonly >NewPassword123!

# Добавить доступ к ключам
ACL SETUSER app_user ~newpattern:*

# Добавить команду
ACL SETUSER cache_worker +incr +decr

# Отключить пользователя
ACL SETUSER temp_user off
```

### Удаление пользователя
```bash
ACL DELUSER temp_user
```

### Сохранение и загрузка
```bash
# Сохранить текущие ACL в файл
ACL SAVE

# Загрузить ACL из файла (перезаписывает текущие)
ACL LOAD
```

## 🧪 Тестирование прав

### Проверить права команды
```bash
# Проверить, может ли пользователь выполнить команду
ACL DRYRUN readonly SET key value
# Вернет OK или ошибку
```

### Попытка выполнить запрещенную команду
```bash
# Подключиться как readonly
AUTH readonly ReadPass456!

# Попробовать записать (должна быть ошибка)
SET test "value"
# Error: NOPERM this user has no permissions to run the 'set' command
```

## 💻 Подключение из приложений

### Python
```python
import redis

# Admin подключение
admin_client = redis.Redis(
    host='localhost',
    port=6379,
    username='admin',
    password='AdminStrongPass123!',
    decode_responses=True
)

# Read-only подключение
readonly_client = redis.Redis(
    host='localhost',
    port=6379,
    username='readonly',
    password='ReadOnlyPass456!',
    decode_responses=True
)

# Cache worker подключение
cache_client = redis.Redis(
    host='localhost',
    port=6379,
    username='cache_writer',
    password='CachePass789!',
    decode_responses=True
)

# Использование
cache_client.setex('cache:user:123', 3600, 'user_data')
value = cache_client.get('cache:user:123')
```
## ⚠️ Важные замечания

### 1. Default пользователь
По умолчанию есть пользователь `default` с полными правами. Для безопасности его нужно отключить:
```bash
user default off
```

### 2. Backwards compatibility
Если используете `requirepass` в конфиге, это устанавливает пароль для `default` пользователя:
```conf
requirepass mypassword
# Эквивалентно:
user default on >mypassword ~* &* +@all
```

### 3. Сохранение ACL
- Изменения через `ACL SETUSER` временные (до перезапуска)
- Используйте `ACL SAVE` для сохранения в файл
- Или добавьте пользователей в `redis.conf` или `users.acl`

### 4. Безопасность паролей
- Используйте длинные случайные пароли (20+ символов)
- Избегайте спецсимволов, которые нужно экранировать в shell
- Храните пароли в переменных окружения или секретах

## 🔐 Best Practices

1. **Принцип наименьших привилегий**: давайте только необходимые права
2. **Отключите default пользователя**: `user default off`
3. **Используйте паттерны ключей**: ограничивайте доступ к префиксам
4. **Отключите опасные команды**: `-@dangerous` для большинства пользователей
5. **Мониторинг**: создайте специального пользователя для мониторинга
6. **Ротация паролей**: периодически меняйте пароли
7. **Логирование**: включите логирование ACL событий

## 📊 Пример архитектуры

```
┌─────────────────┐
│  Web App        │ → app_user (доступ к app:*, session:*)
└─────────────────┘

┌─────────────────┐
│  Cache Service  │ → cache_writer (доступ к cache:*)
└─────────────────┘

┌─────────────────┐
│  Queue Worker   │ → queue_worker (доступ к queue:*, job:*)
└─────────────────┘

┌─────────────────┐
│  Monitoring     │ → monitor (read-only + INFO)
└─────────────────┘

┌─────────────────┐
│  Admin Tools    │ → admin (полный доступ)
└─────────────────┘
```

## 🐛 Отладка проблем

### Проверить права пользователя
```bash
ACL GETUSER username
```

### Проверить, почему команда запрещена
```bash
ACL DRYRUN username COMMAND key
```

### Логи ACL событий
В Redis 7.0+ можно включить логирование ACL:
```conf
acllog-max-len 128
```

Просмотр лога:
```bash
ACL LOG
ACL LOG RESET  # Очистить лог
```

## 📚 Дополнительные ресурсы

- Официальная документация: https://redis.io/docs/management/security/acl/
- Список категорий команд: `ACL CAT`
- Команды в категории: `ACL CAT @read`
