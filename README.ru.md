# go-memorycache
[![en](https://img.shields.io/badge/lang-en-green.svg)](README.md)
[![ru](https://img.shields.io/badge/lang-ru-red.svg)](README.ru.md)

In-memory cache with expiration and eviction.

## Описание
Потокобезопасный кэш данных в памяти.

Память дорогостоящий ресурс и поэтому кэш снабжен механизмом очистки устаревших записей, а также механизмом вытеснения записей при достижении лимита. Механизм очистки запускается при указании CleanupInterval, а вытеснения при указании LimitEnties. Вытеснение работает по принципу - сначала вытесняются устаревшие записи, а за тем в порядке стойкости и методу FIFO.

## Установка
```
go get github.com/a-projects/go-memorycache@latest
```

## Использование
```

import (
	"fmt"
	"time"

	"github.com/a-projects/go-memorycache"
)

func main() {
	// создаём экземпляр кэша
	cache := memorycache.New(memorycache.MemoryCacheOptions{
		// выполнение очистки записей с периодичностью 15 минут
		CleanupInterval: time.Minute * 15,
		// лимит записей в кэше, после которого начинает работать вытеснение
		LimitEntries: 65_536,
		// файл для восстановления кэша при перезапуске приложения
		StoreFile: "cache.bin",
	})

	// пробуем получить значение записи по ключу "foo"
	res, ok := cache.Get("foo")

	// если не удаётся получить значение записи
	// причины:
	//   - никто не добавлял запись с этим ключом
	//   - запись была добавлена, но устарела
	//   - запись была добавлена, но была очищена
	//   - запись была добавлена, но была вытеснена
	if !ok {
		// получаем значение из внешних источников
		res = "bar"

		// добавляем отсутствующую запись в кэш в виде ключ, значение
		cache.Set("foo", res, memorycache.MemoryCacheEntryOptions{
			// время жизни записи
			Expiration: time.Now().Add(time.Minute * 5),
			// стойкость к вытеснению
			Durability: memorycache.Normal,
		})
	}

	// приводим результат к нужному типу и выводим в консоль
	fmt.Printf(res.(string))

	// останавливаем экземпляр, чтобы данные кэша были сохранены в файл
	cache.Close()
}
```