## How to start working with API?


1. Install [docker](https://docs.docker.com/engine/installation/) and [docker-compose](https://docs.docker.com/compose/install/)

2. From projects root build containers

`$ docker-compose build`

3. Start a project

`$ docker-compose up`

4. Bash into Golang container

`$ docker exec -it golang_echo bash`

5. Init SQL schema with fallowing:

`$ cd /db && go run init.go ./init.sql`


### You are ready to GO! Server works on port 8001


***

####__UA__

##Доступні методи

*POST* - "/login" - ендпоінт авторизації, очкує отримати в тілі запиту JSON наступного формату:
```json
{
  "login": "stepan",
  "password": "123456"
}
```
повертає в респонсі JWT токен з тривалістю дії в полі `exp`.


*POST* - "/auth" - ендпоінт для реєстрації нових юзерів. Очікує отримати в тілі запиту JSON наступного формату:
```json
{
	"username": "John Doe",
	"password": "q1w2e3r4",
	"login": "doe",
	"age": 32,
	"telephone": "+380936557609"
}
```
повертає JSON з полем success або з полем error.


*GET* - `/api/v1/film` - ендпоінт для вибору усіх достпних фільмів у базі данних. підтримує пагінацію параметрами URL
в форматі `limit`/`offset`. Наприклад, запит `/api/v1/film?limit=5&offset=10` поверне 5 фільмів з відступом у десять фільмів. 

Підтримує фільтрацію за роком та жанрами (genre_id розділеними через кому). Наприклад:

`api/v1/film?genre=1,2&year=2001&limit=10`

Поверне до 10 фільмів з жанрам під айді 1 і 2 рік випуску яких - 2001.

Також повертає поля count - всього фільмів та left - булове значення - чи залишились фільми в базі.
Приклад респонсу:
```json
{
    "left": false,
    "count": 3,
    "result": [
        {
            "id": 2,
            "name": "Commando",
            "year": 1992,
            "added_at": "2017-06-20T19:37:32.932759Z",
            "genres": [
                {
                    "id": 1,
                    "name": "Comedy",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                },
                {
                    "id": 2,
                    "name": "Horror",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                }
            ]
        },
        {
            "id": 3,
            "name": "Red Mist",
            "year": 2000,
            "added_at": "2017-06-20T19:37:53.15788Z",
            "genres": [
                {
                    "id": 3,
                    "name": "Drama",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                }
            ]
        }
    ]
}
```


*POST* - `/api/v1/rent` - ендпоінт оренди фільму. Доступний лише авторизованим юзерам. Під авторизованими юзерами мається 
на увазі той, у якого присутній HTTP Header у запиті з токеном в форматі:

`Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTgyNDgxNzAsIm5hbWUiOiJsZW1waXkiLCJ1c2VyX2lkIjoxfQ.8O_rCqLSBnYaHsZph2Yp8JeV4wtQ_MHv3D5c_5_WTw8`

В тілі запиту ендпоінт очікує отримати JSON з айді орендованого фільму:
```json
{
	"film_id": 1
}
```
повертає JSON з полем success або з полем error.


*POST* - `/api/v1/finish` - ендпоінт для завершення оренди фільму  Доступний лише авторизованим юзерам.
В тілі запиту ендпоінт очікує отримати JSON з айді орендованого фільму:
```json
{
	"film_id": 1
}
```
повертає JSON з полем success або з полем error.


*POST* - `/api/v1/film` - ендпоінт для додавання фільмів в базу. Він знадобиться для наповнення бази фільмами. Лише для Авторизованих юзерів. Оскільки спочатку вона буде пустою. Приймає він JSON в форматі:
```json
{
	"name": "Commando",
	"year": 1990,
	"genres": [1,2]
}
```


*GET* - `/api/v1/rented-film` - ендпоінт для вибору усіх всіх орендованих фільмів авторизованим юзером. підтримує пагінацію параметрами URL
в форматі `limit`/`offset`.

Підтримує фільтрацію за роком та жанрами (genre_id розділеними через кому). Наприклад:

`api/v1/rented-film?genre=1,2&year=2001&limit=10`

Поверне до 10 орендованих юзером фільмів з жанрам під айді 1 і 2 рік випуску яких - 2001.

 Також повертає поля count - всього фільмів та left - булове значення - чи залишились фільми в базі.
Приклад респонсу:
```json
{
    "left": false,
    "count": 3,
    "result": [
        {
            "id": 2,
            "name": "Commando",
            "year": 1992,
            "added_at": "2017-06-20T19:37:32.932759Z",
            "genres": [
                {
                    "id": 1,
                    "name": "Comedy",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                },
                {
                    "id": 2,
                    "name": "Horror",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                }
            ]
        },
        {
            "id": 3,
            "name": "Red Mist",
            "year": 2000,
            "added_at": "2017-06-20T19:37:53.15788Z",
            "genres": [
                {
                    "id": 3,
                    "name": "Drama",
                    "added_at": "2017-06-20T19:30:28.683215Z"
                }
            ]
        }
    ]
}
```
