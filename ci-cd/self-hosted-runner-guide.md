# Настройка self-hosted runner для GitHub Actions

## Что такое self-hosted runner

Self-hosted runner — это агент, который запускается на вашей машине и выполняет задачи из GitHub Actions пайплайна локально. В отличие от облачных runner'ов GitHub (`ubuntu-latest`), ваша машина сама выполняет джобы.

Типичный сценарий:
- Джобы `lint`, `test`, `build` выполняются на облачных runner'ах GitHub
- Джоб `deploy` выполняется на вашей машине через self-hosted runner

## Требования

- Установленный Docker и Docker Compose
- Доступ к интернету
- Права администратора или владельца репозитория на GitHub

## Шаг 1. Получение токена и команд установки

Перейдите в репозиторий на GitHub:

```
Settings -> Actions -> Runners -> New self-hosted runner
```

Выберите вашу ОС. GitHub сгенерирует команды с уникальным токеном.

Важно: токен одноразовый и истекает через 1 час. Если не успели — получите новый токен на той же странице.

## Шаг 2. Установка runner'а (Linux)

```bash
# Создайте папку БЕЗ пробелов в пути
mkdir ~/actions-runner && cd ~/actions-runner

# Скачайте runner
curl -o actions-runner-linux-x64-2.333.1.tar.gz -L \
  https://github.com/actions/runner/releases/download/v2.333.1/actions-runner-linux-x64-2.333.1.tar.gz

# Распакуйте
tar xzf ./actions-runner-linux-x64-2.333.1.tar.gz
```

Папка обязательно должна быть без пробелов в пути. Путь вида `/home/user/Рабочий стол/actions-runner` вызовет ошибки при выполнении джобов.

## Шаг 3. Регистрация runner'а

```bash
./config.sh --url https://github.com/ВАШ_АККАУНТ/ВАШ_РЕПОЗИТОРИЙ --token ВАШ_ТОКЕН
```

В процессе настройки будут заданы вопросы:

| Вопрос | Рекомендация |
|--------|-------------|
| Runner group | Enter (оставить default) |
| Runner name | Любое имя, например `my-local` |
| Labels | Enter (оставить self-hosted) |
| Work folder | Enter (оставить _work) |

## Шаг 4. Запуск runner'а

```bash
./run.sh
```

Успешный запуск выглядит так:

```
Connected to GitHub
Current runner version: '2.333.1'
2026-04-21 08:37:31Z: Listening for Jobs
```

Runner должен быть запущен в момент выполнения пайплайна. Если runner не запущен — джоб зависнет в очереди и упадёт с ошибкой `No runner available`.

## Шаг 5. Настройка workflow файла

Добавьте джоб `deploy` с `runs-on: self-hosted` в ваш `.github/workflows/ci.yml`:

```yaml
name: CI/CD

on:
  push:
    branches: [master, main]
  pull_request:
    branches: [master, main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - name: Install ruff
        run: pip install ruff
      - name: Run ruff
        run: ruff check .

  test:
    name: Tests
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
          cache: "pip"
      - name: Install dependencies
        run: pip install -r requirements.txt
      - name: Run tests
        run: pytest tests -v --tb=short

  build:
    name: Build & Push Docker image
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/main'

    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v4

      - name: Set image name
        run: echo "IMAGE=ghcr.io/$(echo '${{ github.repository }}' | tr '[:upper:]' '[:lower:]')" >> $GITHUB_ENV

      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ${{ env.IMAGE }}:latest
            ${{ env.IMAGE }}:${{ github.sha }}

  deploy:
    name: Deploy on self-hosted
    runs-on: self-hosted
    needs: build
    if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v4

      - name: Set image name
        run: echo "IMAGE=ghcr.io/$(echo '${{ github.repository }}' | tr '[:upper:]' '[:lower:]')" >> $GITHUB_ENV

      - name: Create .env file
        run: |
          cat > .env << EOF
          POSTGRES_HOST=db
          POSTGRES_PORT=5432
          POSTGRES_DB=your_db
          POSTGRES_USER=postgres
          POSTGRES_PASSWORD=postgres
          SECRET_KEY=${{ secrets.SECRET_KEY }}
          EOF

      - name: Log in to GHCR
        run: echo "${{ secrets.GITHUB_TOKEN }}" | docker login ghcr.io -u ${{ github.actor }} --password-stdin

      - name: Stop old containers
        run: docker compose down || true

      - name: Pull latest image
        run: docker pull ${{ env.IMAGE }}:latest

      - name: Start containers
        run: docker compose up -d
```

## Шаг 6. Переменные окружения и секреты

Файл `.env` не хранится в репозитории (добавлен в `.gitignore`). Runner каждый раз делает чистый checkout, поэтому `.env` нужно создавать в джобе.

Секретные значения выносятся в GitHub Secrets:

```
Settings -> Secrets and variables -> Actions -> New repository secret
```

| Имя секрета | Описание |
|-------------|----------|
| SECRET_KEY | Секретный ключ JWT или другие чувствительные данные |

Несекретные значения (хосты, порты, имена БД) можно хардкодить прямо в workflow.

## Шаг 7. Безопасность

Self-hosted runner на публичном репозитории опасен — любой может открыть Pull Request и выполнить произвольный код на вашей машине.

Минимальная защита для публичного репозитория:

```
Settings -> Actions -> General -> Fork pull request workflows
-> Require approval for all external contributors
```

Оптимальный вариант — использовать приватный репозиторий с self-hosted runner'ом.

## Шаг 8. Проверка

После пуша в `main` или `master` пайплайн запустится автоматически. Успешное выполнение всей цепочки выглядит так:

```
Lint (7s) -> Tests (1m 39s) -> Build & Push Docker image (1m 28s) -> Deploy on self-hosted
```

Для ручного триггера без изменений кода:

```bash
git commit --allow-empty -m "trigger CI"
git push
```

## Типичные ошибки

| Ошибка | Причина | Решение |
|--------|---------|---------|
| `No such file or directory` в пути к GITHUB_ENV | Пробел в пути к папке runner'а | Перенести `actions-runner` в путь без пробелов |
| `No runner available` | Runner не запущен | Запустить `./run.sh` перед пушем |
| `.env not found` | Runner делает чистый checkout | Создавать `.env` в джобе через `cat > .env` |
| `Unexpected symbol: \|` в выражении `github.repository \| lower` | Фильтр `\| lower` не поддерживается в GitHub Actions | Использовать `tr '[:upper:]' '[:lower:]'` через bash |
| Токен истёк при `config.sh` | Токен одноразовый, живёт 1 час | Получить новый токен в Settings -> Actions -> Runners |
