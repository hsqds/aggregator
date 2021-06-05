# RSS аггрегатор

Настройки аггрегатора хранятся в `config.json`.
В корне находятся две секции: `feeds` и `db`.  
`db` содержит настройки подключения к БД.  
```json
{
    "host": "localhost",
    "port": "5432",
    "username": "postgres",
    "password": "pass",
    "dbname": "rss-reader"
}
```
В `feeds` хранится массив объектов:
```json
{
    "url": "https://blog.golang.org",
    "rules": {
        "postPath": "/feed/entry",
        "titlePath": "title",
        "linkPath": "link@href",
        "descriptionPath": "summary"
    }
}
```

В поле `url` находится адрес rss-фида или доменное имя.
Поле `rules` содержит объект с описанием правил парсинга фида.

`postPath` - XPath элемента с описанием поста  
`titlePath` - путь к элементу, содержащему заголовок поста  
`linkPath` -  путь к элементу, содержащему ссылку на пост  
`descriptionPath` - путь к элементу, содержащему краткое описание поста  

XPath может оканчиваться именем аттрибута, в котором хранится искомое значение: `@attrName`

## Зависимости

 * docker
 * docker-compose
 * go >= 1.16

## Запуск

* `docker-compose up -d` - запустит postgres, adminer, накатит миграции

## Миграции
Если вдруг надо запустить миграции руками
* `make install-tern` - установит tern - утилиту для организации миграций БД
* `make migrate` - создаст таблицу и индексы в базе

## Обновление базы постов
* `go run cmd/update` - запустит обновление базы постов
```
  -c string
        config file path (default "./config.json")
```

## Поиск 
`go run ./cmd/search`
```
  -c string
    	config file path (default "./config.json")
  -l int
    	results limit (default 10)
  -s string
    	search request
```