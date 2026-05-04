# Deployment

## 1. Создаем deployment

Для этого выполним команду:

```bash
kubectl apply -f ~/school-dev-k8s/practice/3.application-abstractions/3.deployment/
```

Проверяем список pods, для этого выполним команду:

```bash
kubectl get pod
```

Результат должен быть примерно таким:

```bash
NAME                             READY     STATUS              RESTARTS   AGE
my-deployment-7c768c95c4-47jxz   0/1       ContainerCreating   0          2s
my-deployment-7c768c95c4-lx9bm   0/1       ContainerCreating   0          2s
```

Проверяем список replicaset, для этого выполним команду:

```bash
kubectl get replicaset
```

Результат должен быть примерно таким:

```bash
NAME                       DESIRED   CURRENT   READY     AGE
my-deployment-7c768c95c4   2         2         2         1m
```

## 2. Обновляем версию image

Обновляем версию image для container в deployment my-deployment.
Для этого выполним команду:

```bash
kubectl set image deployment my-deployment nginx=nginx:1.13
```

Проверяем результат, для этого выполним команду:

```bash
kubectl get pod
```

Результат должен быть примерно таким:

```bash
NAME                             READY     STATUS              RESTARTS   AGE
my-deployment-685879478f-7t6ws   0/1       ContainerCreating   0          1s
my-deployment-685879478f-gw7sq   0/1       ContainerCreating   0          1s
my-deployment-7c768c95c4-47jxz   0/1       Terminating         0          5m
my-deployment-7c768c95c4-lx9bm   1/1       Running             0          5m
```

И через какое-то время вывод этой команды станет:

```bash
NAME                             READY     STATUS    RESTARTS   AGE
my-deployment-685879478f-7t6ws   1/1       Running   0          33s
my-deployment-685879478f-gw7sq   1/1       Running   0          33s
```

Проверяем что в новых pod новый image. Для этого выполним команду,
заменив имя pod на имя вашего pod:

```bash
kubectl describe pod my-deployment-685879478f-7t6ws
```

Результат должен быть примерно таким:

```bash
    Image:          nginx:1.13
```

Проверяем что стало с replicaset, для этого выполним команду:

```bash
kubectl get replicaset
```

Результат должен быть примерно таким:

```bash
NAME                       DESIRED   CURRENT   READY     AGE
my-deployment-685879478f   2         2         2         2m
my-deployment-7c768c95c4   0         0         0         7m
```

## 3. Чистим за собой кластер

```bash
kubectl delete deployment --all
```

## 4. Самостоятельная работа

1. Отредактируйте deployment из предыдущего задания таким образом, чтобы он запускался с одной репликой и rolling update проходил без downtime. То есть при обновлении образа всегда должен оставаться один рабочий pod.

    Используйте поля maxSurge и maxUnavailable со значениями, приведенными в видео для деплойментов с одной репликой.

2. Обновите образ с `nginx:1.20` на `nginx:1.21`. Посмотрите в каком порядке запускаются и тушатся  pod'ы в момент обновления.

3. Выполните команду (скопировать именно всё):

```bash
kubectl get deployment my-deployment -o custom-columns='NAME:.metadata.name,MAXSURGE:.spec.strategy.rollingUpdate.maxSurge,MAXUNAVAILABLE:.spec.strategy.rollingUpdate.maxUnavailable'
```

> Это еще одна возможность ключа -o. Она позволяет вывести описание объекта с пользовательским набором полей.

4. Отправьте результат ее выполнения как ответ на этот шаг.

5. Удалите deployment после выполнения задания.



## Сопоставьте названия типов объектов в Kubernetes с их назначением.

- [x] **Pod** - Объединяет несколько контейнеров в одну минимальную логическую единицу
- [x] **Replicaset** - Контролирует количество реплик приложения  
- [x] **Deployment** - Позволяет выполнять обновление версий образов приложений

## Выберите все верные утверждения из списка ниже.

- [x] В реальной жизни для запуска приложений в кластере Kubernetes в подавляющем большинстве случаев используются Deployment'ы, а Pod и Replicaset нужны в основном только как служебные абстракции
- [ ] Replicaset может самостоятельно обновить версии образов в своих pod'ах
- [ ] Контейнеры одного pod'а могут быть запущены на разных host'ах
- [x] Deployment по умолчанию осуществляет обновление приложений с помощью создания нового replicaset'а и постепенного изменения количества реплик в старом и новом replicaset'ах
- [x] Минимальной единицей Kubernetes является pod
- [ ] Единственной возможной стратегией обновления в Deployment является rollingUpdateate
