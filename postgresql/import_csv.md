# Импортируем CSV на практике
В этой лекции выполним практический пример импорта CSV-файла в PostgreSQL с помощью команды COPY. В качестве источника данных возьмём открытый датасет о пациентах с коронавирусом с сайта Kaggle.

## 1. Получаем данные
На сайте Kaggle можно найти множество датасетов в разных форматах — JSON, CSV, SQLite и других. Для примера используем CSV-файл patient.csv из датасета о коронавирусе.

Файл содержит данные о пациентах:

```
patient_id — уникальный идентификатор пациента
sex — пол
birth_year — год рождения
country, region — страна и регион
infection_reason — возможная причина заражения
infected_by — кем был заражён (если известно)
confirmed_date — дата подтверждения диагноза
released_date — дата выписки
deceased_date — дата смерти
state — текущее состояние пациента
```
## 2. Создаём таблицу
Создадим таблицу patients в базе данных testdb, чтобы импортировать в неё данные.

```sql
CREATE TABLE patients (
    patient_id     INTEGER PRIMARY KEY,
    sex            TEXT,
    birth_year     INTEGER,
    country        TEXT,
    region         TEXT,
    infection_reason TEXT,
    infected_by    INTEGER,
    contact_number INTEGER,
    confirmed_date DATE,
    released_date  DATE,
    deceased_date  DATE,
    state          TEXT
);
```

                  
## 3. Импортируем данные через psql
Откройте консоль psql и подключитесь к базе данных:

```bash
psql -U postgres -d testdb
```
                  
После подключения выполните команду `COPY`, указав путь к файлу:

```sql
COPY patients
FROM 'C:\Users\User\Downloads\coronavirus_dataset\patient.csv'
DELIMITER ',' 
CSV HEADER;
```
                  
Если всё сделано правильно, PostgreSQL выведет количество импортированных строк:
```sql
COPY 4812
```
                  
## 4. Проверяем результат

```sql
SELECT * FROM patients LIMIT 10;
```
                  
После импорта можно использовать SQL для анализа данных, например:

```sql
-- Посчитаем количество умерших пациентов
SELECT COUNT(*) 
FROM patients
WHERE deceased_date IS NOT NULL;

                  
-- Посмотрим возраст умерших
SELECT birth_year 
FROM patients
WHERE deceased_date IS NOT NULL
ORDER BY birth_year;
```
                  
## 5. Пример анализа
Используя SQL-запросы, можно самостоятельно анализировать такие датасеты:
- вычислять долю выздоровевших
- строить выборки по возрасту, полу, региону и т.д.

Это отличный пример того, как SQL позволяет работать с реальными данными и проверять гипотезы без внешних инструментов.

## 6. Другие варианты импорта

- `JSON` — импорт с помощью функций `json_populate_recordset()` или `COPY` в колонку типа `json/jsonb.`
- `Excel (XLSX)` — предварительно сохранить как CSV и импортировать тем же способом.
- `SQL-скрипты` — выполняются напрямую через `\i путь_к_файлу.sql.`
- `psql \copy` — клиентская версия `COPY`, когда сервер не имеет доступа к локальным файлам.

## Итог
Импорт CSV в PostgreSQL выполняется через команду COPY.
Важно заранее создать таблицу с подходящей структурой.
psql позволяет напрямую импортировать и сразу проверять данные с помощью SQL-запросов.
CSV — удобный формат для первых шагов в анализе данных и интеграции внешних источников.