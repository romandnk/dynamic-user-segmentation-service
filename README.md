## Запуск
Перед запуском необходимо настроить конфиги в папке configs. 

Находясь в папке с проектом введите `make run` или `docker compose -f ./deployments/docker-compose.yaml up -d --build
`

## Методы HTTP

### 1) Создание сегмента

- **HTTP метод**: POST
- **Путь**: `api/v1/segments`

**Curl запрос**:

```bash
curl --location 'http://172.26.0.3:8080/api/v1/segments' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "AVITO",
    "auto_add_percentage": "100%"
}'
```
Коды ответов:

- 201 (успешно)
- 400
- 500

Ограничения:

- Название должно состоять из больших букв.
- Строка с процентами должна быть в формате 10%, 5%, 100% (от 1 до 100, целые числа).

### 2) Удаление сегмента

- **HTTP метод**: DELETE
- **Путь**: `api/v1/segments`

**Curl запрос**:

```bash
curl --location --request DELETE 'http://172.26.0.3:8080/api/v1/segments' \
--header 'Content-Type: application/json' \
--data '{
    "slug": "AVITO"
}'
```
Коды ответов:

- 200 (успешно)
- 400
- 500

Ограничения:

- Название должно состоять из больших букв.
- Название не может быть пустым.

### 3) Добавление и удаление сегментов пользователя

- **HTTP метод**: POST
- **Путь**: `api/v1/users`

**Curl запрос**:

```bash
curl --location 'http://172.26.0.3:8080/api/v1/users' \
--header 'Content-Type: application/json' \
--data '{
    "segments_to_add": ["AVITO"],
    "segments_to_delete": ["TEST"],
    "user_id": 1
}'
```
Коды ответов:

- 201 (успешно)
- 400
- 500

Ограничения:

- Название сегментов должно состоять из больших букв.
- Массив сегментов для добавления и удаления оба не должны быть пустыми.
- Идентификатор пользователя должен быть больше нуля.

### 4) Получение активных сегментов пользователя

- **HTTP метод**: POST
- **Путь**: `api/v1/users/active_segments`

**Curl запрос**:

```bash
curl --location 'http://172.26.0.3:8080/api/v1/users/active_segments' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1
}'
```
Коды ответов:

- 200 (успешно)
- 400
- 500

**JSON ответ**

```JSON
{
  "segments": ["AVITO"]
}
```

Ограничения:

- Идентификатор пользователя должен быть больше нуля.

### 5) Получение ссылки на отчет по пользователям в течении какого-то месяца

- **HTTP метод**: POST
- **Путь**: `api/v1/users/report`

**Curl запрос**:

```bash
curl --location 'http://172.26.0.3:8080/api/v1/users/report' \
--header 'Content-Type: application/json' \
--data '{
    "date": "2023-08"
}'
```
Коды ответов:

- 200 (успешно)
- 400
- 500

**JSON ответ**

```JSON
{
  "report_url": "http://172.26.0.3:8080/api/v1/users/report/1f674039-d035-4b1a-ac8b-51b67ab350e1"
}
```

Ограничения:

- Дата в формате год-месяц (например, 2023-08)

### 6) Получение файла отчета в формате CSV

- **HTTP метод**: GET
- **Путь**: `api/v1/users/report/{id}`

**Curl запрос**:

```bash
curl --location 'http://172.26.0.3:8080/api/v1/users/report/1f674039-d035-4b1a-ac8b-51b67ab350e1'
```
Коды ответов:

- 200 (успешно)
- 400
- 500

## Дополнительные задания

### Отчет по пользователям
При запросе на получении ссылки генерируется CSV файл локально с уникальным ID, который содержит информацию из таблицы operations в БД PostgreSQL. При открытии ссылки происходит скачивание файла в формате CSV, который ищется локально. Место, где хранятся отчеты, можно конфигурировать в файле конфигурации.

### Автоматическое добавление пользователя в сегмент
В главной горутине создается горутина с тикером, которая раз в определенное время (конфигурируется в файле конфигурации) сканирует БД и автоматически добавляет пользователей в сегмент.

Алгоритм:
- Считаем количество пользователей (amount).
- Выбираем сегменты, у которых стоит процент добавления.
- Проходимся по каждому сегменту, считаем количество пользователей, которые исходя из процента сегмента и количества пользователей, должны быть в сегменте (amount * percent / 100).
- Выбираем количество пользователей, которые уже находятся в этом сегменте (n).
- Если их количество меньше нужного числа, то выбираем (amount - n) пользователей, которые должны быть в сегменте и добавляем им этот сегмент.