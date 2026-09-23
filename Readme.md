### Система бронирования

#### Урл запущенного приложения - http://localhost:8081

Запуск приложения:

Сборка:
```bash
make docker-build
```

Пересоздание контейнеров:
```bash
make docker-up
```

Применение миграций:
```bash
make migrate-up
```

Откат миграций:
```bash
make migrate-down
```
---

Ручки:

События:

<details><summary>GET http://localhost:8081/events - получить список событий</summary>

```bash
curl --location 'http://localhost:8081/events'
```

</details>

<details><summary>POST http://localhost:8081/events - добавить новое событие</summary>

```bash
curl --location 'http://localhost:8081/events' \
--header 'Content-Type: application/json' \
--data '{
  "name": "test event",
  "starts_at" : "2026-11-12 11:12:13"
}'
```

</details>

<details><summary>GET http://localhost:8081/events/{id} - получить конкретное событие по его id</summary>

```bash
curl --location 'http://localhost:8081/events/1'
```

</details>

---

Места на событии:

<details><summary>POST http://localhost:8081/events/{eventID}/seats - добавить места на событие по его id</summary>

```bash
curl --location 'http://localhost:8081/events/1/seats' \
--header 'Content-Type: application/json' \
--data '{
  "seats": [
    {
      "number": "a1",
      "price": 1500
    },
     {
      "number": "a2",
      "price": 1500
    }
  ]
}'
```

</details>

<details><summary>GET http://localhost:8081/events/{eventID}/seats - получить список мест по конкретному событию</summary>

```bash
curl --location 'http://localhost:8081/events/4/seats'
```

</details>

---

Бронирование мест:

<details><summary>POST http://localhost:8081/events/{eventID}/seats/{seatID}/reserve - бронирование  места, статус reserved</summary>

```bash
curl --location 'http://localhost:8081/events/4/seats/9/reserve' \
--header 'Content-Type: application/json' \
--data '{
    "userId": 1
}'
```

</details>

<details><summary>GET http://localhost:8081/users/{userID}/bookings - список бронирований юзера</summary>

```bash
curl --location 'http://localhost:8081/users/1/bookings'
```

</details>

<details><summary>POST http://localhost:8081/bookings/{bookingID}/confirm - подтверждение бронирования, статус confirmed</summary>

```bash
curl --location 'http://localhost:8081/bookings/2/confirm' \
--header 'Content-Type: application/json' \
--data '{
  "userId": 1
}'
```

</details>

<details><summary>DELETE http://localhost:8081/bookings/{bookingID} - отмена бронирования, статус canceled</summary>

```bash
curl --location 'http://localhost:8081/users/1/bookings'
```

</details>

<details><summary>DELETE http://localhost:8081/bookings/{bookingID} - отмена бронирования, статус canceled</summary>

```bash
curl --location --request DELETE 'http://localhost:8081/bookings/3' \
--header 'Content-Type: application/json' \
--data '{
  "userId": 1
}'
```

</details>


