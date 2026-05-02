# Sorted set

**Sorted set** — аналогичен обычному сету с уникальными значением за
исключением того что в sorted set с каждым значением ассоциируется score
благодаря которому set сортирует значения.

## ОСОБЕННОСТИ

- Максимальная длина хеш сета 2^32 - 1 элементов
- Может хранить список изображений (Base64)
- Отличное решение для построения рейтингов

### Базовые операции

- `zadd <key> <score> <member>` — добавить элемент с указанным score  
  👉 `zadd players 100 "Alice"`  

- `zrem <key> <member>` — удалить элемент  
  👉 `zrem players "Alice"`  

- `zscore <key> <member>` — получить score элемента  
  👉 `zscore players "Alice"` → `100`  

- `zcard <key>` — получить количество элементов  
  👉 `zcard players` → `3`  

- `zcount <key> <min> <max>` — посчитать количество элементов в диапазоне  
  👉 `zcount players 50 150` → `2`  

### Получение элементов

- `zrange <key> <start> <stop>` — элементы по индексу (возрастание score)  
  👉 `zrange players 0 -1` → `["Alice","Bob","Eve"]`  

- `zrevrange <key> <start> <stop>` — элементы по индексу (убывание score)  
  👉 `zrevrange players 0 1` → `["Eve","Bob"]`  

- `zrangebyscore <key> <min> <max>` — элементы по диапазону score (возрастание)  
  👉 `zrangebyscore players 50 200` → `["Alice","Bob"]`  

- `zrevrangebyscore <key> <max> <min>` — элементы по диапазону score (убывание)  
  👉 `zrevrangebyscore players 200 50` → `["Bob","Alice"]`  

- `zrangebylex <key> <min> <max>` — элементы в лексикографическом порядке  
  👉 `zrangebylex players - +`  

- `zrevrangebylex <key> <max> <min>` — элементы в обратном лексикографическом порядке  
  👉 `zrevrangebylex players + -`  

### Ранги и позиции

- `zrank <key> <member>` — индекс элемента (возрастание score)  
  👉 `zrank players "Alice"` → `0`  

- `zrevrank <key> <member>` — индекс элемента (убывание score)  
  👉 `zrevrank players "Alice"` → `2`  

### Изменение значений

- `zincrby <key> <increment> <member>` — увеличить score элемента  
  👉 `zincrby players 50 "Alice"` → `150`  

### Удаление диапазонов

- `zremrangebyscore <key> <min> <max>` — удалить элементы по score  
  👉 `zremrangebyscore players 0 100`  

- `zremrangebyrank <key> <start> <stop>` — удалить элементы по индексу  
  👉 `zremrangebyrank players 0 1`  

- `zremrangebylex <key> <min> <max>` — удалить элементы по диапазону (лексикографически)  
  👉 `zremrangebylex players [A [C`  

### Извлечение с удалением

- `zpopmin <key>` — извлечь минимальный score  
  👉 `zpopmin players` → `("Alice",100)`  

- `zpopmax <key>` — извлечь максимальный score  
  👉 `zpopmax players` → `("Eve",300)`  

- `bzpopmin <key> <timeout>` — блокирующая версия zpopmin  
  👉 `bzpopmin players 5`  

- `bzpopmax <key> <timeout>` — блокирующая версия zpopmax  
  👉 `bzpopmax players 5`  

### Агрегация множеств

- `zunionstore <dest> <numkeys> <key ...>` — объединение множеств  
  👉 `zunionstore all_players 2 players1 players2`  

- `zinterstore <dest> <numkeys> <key ...>` — пересечение множеств  
  👉 `zinterstore common 2 players1 players2`  

### Итерация

- `zscan <key> <cursor>` — итерация по множеству  
  👉 `zscan players 0` → `(0, ["Alice","100","Bob","200"])`  
