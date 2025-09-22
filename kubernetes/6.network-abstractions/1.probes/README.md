### Resources and Probes

**Probes (зонды)** - это механизм в Kubernetes для проверки состояния контейнеров. Они помогают Kubernetes понимать, когда контейнер готов к работе, жив и здоров.

**Liveness Probe**

- Проверяет, жив ли контейнер. Если проверка fails, контейнер будет перезапущен.

**Readiness Probe**

- Проверяет, готов ли контейнер принимать трафик. Если проверка fails, контейнер будет удален из эндпоинтов сервиса.

**Startup Probe**

- Используется для проверки, что приложение запустилось. После успешной проверки остальные probe продолжают работу.


Создаем deployment с ресурсами и пробами

```bash
kubectl apply -f ~/school-dev-k8s/practice/6.network-abstractions/1.probes/deployment-with-stuff.yaml
```

Смотрим что получилось

```bash
kubectl get pod
```

Должны увидеть что-то типа такого:

```bash
NAME                             READY   STATUS    RESTARTS   AGE
my-deployment-69695544f6-v97jr   1/1     Running   0          20s
my-deployment-69695544f6-xcpq9   1/1     Running   0          20s
```

Поменяем специально Readiness Probe на заведомо неверную в манифесте и применим его снова

```bash
vim ~/school-dev-k8s/practice/6.network-abstractions/1.probes/deployment-with-stuff.yaml
kubectl apply -f ~/school-dev-k8s/practice/6.network-abstractions/1.probes/deployment-with-stuff.yaml
```

Смотрим что получилось

```bash
kubectl get pod
```

Видим что pod'ы не переходят в статус `1/1 Running`. Смотрим describe pod'а:
```bash
kubectl describe po my-deployment-845d88fdcf-9bd29
```
Чистим за собой кластер
```bash
kubectl delete deployment --all
```
