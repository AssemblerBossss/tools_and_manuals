# Hash

**Hash** — представляет собой отображение (map) из ключей и значений, аналогичное `Map` из JavaScript.

## ОСОБЕННОСТИ

- Максимальная длина хеш-сета 2^32 - 1 элементов
- Может хранить список изображений (Base64)
- Подходит для сохранения сериализованных JS, Ruby и др. объектов

## КОМАНДЫ

### Базовые операции

- `hset <key> <field> <value>` — установить значение для поля в хеше  
  Пример: `hset user:1 name "Alice"`

- `hget <key> <field>` — получить значение поля из хеша  
  Пример: `hget user:1 name` → `"Alice"`

- `hdel <key> <field>` — удалить поле из хеша  
  Пример: `hdel user:1 name`

- `hexists <key> <field>` — проверить наличие поля в хеше  
  Пример: `hexists user:1 name` → `1`

- `hgetall <key>` — получить все поля и значения хеша  
  Пример: `hgetall user:1` → `["name","Alice","age","25"]`

- `hmset <key> <field1> <value1> <field2> <value2> ...` — установить несколько полей и значений  
  Пример: `hmset user:1 name "Alice" age 25`

- `hmget <key> <field1> <field2> ...` — получить значения нескольких полей  
  Пример: `hmget user:1 name age` → `["Alice","25"]`

- `hsetnx <key> <field> <value>` — установить значение в хеше, если поля не существует  
  Пример: `hsetnx user:1 country "USA"`


### Операции с числовыми значениями

- `hincrby <key> <field> <increment>` — увеличить значение поля на заданное число  
  Пример: `hincrby user:1 age 1` → `26`

- `hincrbyfloat <key> <field> <increment>` — увеличить значение поля на число с плавающей точкой  
  Пример: `hincrbyfloat user:1 rating 0.5` → `4.5`


### Информация о хеше

- `hlen <key>` — получить количество полей в хеше  
  Пример: `hlen user:1` → `3`

- `hkeys <key>` — получить все поля хеша  
  Пример: `hkeys user:1` → `["name","age","country"]`

- `hvals <key>` — получить все значения хеша  
  Пример: `hvals user:1` → `["Alice","26","USA"]`

- `hstrlen <key> <field>` — получить длину значения поля  
  Пример: `hstrlen user:1 name` → `5`


### Итерация

- `hscan <key> <cursor>` — итерация по полям и значениям хеша  
  Пример: `hscan user:1 0` → `(0, ["name","Alice","age","26"])`
