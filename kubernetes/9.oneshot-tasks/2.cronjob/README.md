# CronJob

**CronJob** — это объект Kubernetes, который позволяет запускать задачи (Jobs) по расписанию, аналогично cron в Linux. Он автоматически создаёт и управляет Job'ами в заданное время или с заданной периодичностью.

**Важные параметры**

- `startingDeadlineSeconds` - Максимальное время (в секундах) после запланированного времени запуска, в течение которого CronJob может быть запущен.
  - Если по каким-то причинам запуск Job был пропущен (например, кластер был недоступен), то Job будет запущен, если задержка не превышает startingDeadlineSeconds.
  - Если задержка больше, запуск пропускается.

- `concurrencyPolicy` - Определяет, как CronJob обрабатывает одновременные запуски, если предыдущий Job ещё не завершился.
  - `Allow` - Разрешить запускать несколько Job одновременно (по умолчанию).
  - `Forbid` - Запретить запуск нового Job, если предыдущий ещё выполняется.
  - `Replace` - Прервать (удалить) текущий выполняющийся Job и запустить новый.
- `successfulJobsHistoryLimit` - Максимальное количество успешных Job, которые Kubernetes будет хранить в истории.
  - По умолчанию: 3
- `failedJobsHistoryLimit` - Максимальное количество неудачных Job, которые Kubernetes будет хранить в истории.
  - По умолчанию: 1



1) Создаем CronJob:

```bash
kubectl apply -f ~/school-dev-k8s/practice/9.oneshot-tasks/2.cronjob/cronjob.yaml
```

2) Проверяем

```bash
kubectl get cronjob
```

Видим:

```bash
NAME    SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
hello   */1 * * * *   False     0        <none>          14s
```

3) Через минуту пробуем посмотреть на Job'ы

```bash
kubectl get job
```

Видим созданный Job

```bash
NAME               COMPLETIONS   DURATION   AGE
hello-1552924260   1/1           2s         49s
```

4) Смотрим на Pod'ы

```bash
kubectl get pod
```

Видим Pod

```bash
NAME                     READY   STATUS      RESTARTS   AGE
hello-1552924260-gp7pk   0/1     Completed   0          80s
```

5) Если мы подождем 5-10 минут, то увидим что старые Job'ы и Pod'ы удаляются по мере появления новых

```bash
kubectl get job,pod
```

6) Удаляем CronJob

```bash
kubectl delete -f ~/school-dev-k8s/practice/9.oneshot-tasks/2.cronjob/cronjob.yaml
```
