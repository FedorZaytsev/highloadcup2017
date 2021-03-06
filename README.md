# Solution for highload cup 2017

## Результат
* Score: 249.87749
* Место: 52 из 295
* Место среди решений на GO: 12 из 55

## Библиотеки
* fasthttp в качестве сервера
* easyjson для парсинга

## Как хранил данные?
В конкурсе id объектов шли по порядку, поэтому данные хранились следующим образом: некоторое число объектов первых объектов каждого типа хранятся в обычном массиве, все что не помещается в массив хранится в хэш-мапе. Для первого тура количество объектов было следующим:
* Пользователей - 1 млн.
* Локаций - 100 тысяч
* Посещений - 10 млн.

Так как сервер тестируется на 64 битной системе, где размер указателя - 8 байт, то размер всех трех массивов с указателями на объекты будет 11.100.000*8 ~= 85 Mb, что не так и много, а скорость работы гораздо выше чем у хэш маппы.

Для случая если какой-то id не влезет в ограничения, он сохраняется в обычном map[int]*Type.

Почему не sync.Map? В условии задачи это не написано, но при тестировании не было ситуации одновременного чтения или записи - этиоперации были разграниченны.

## Индексы
К сожалению я сделал только один индекс, который конечно сильно мне помог в скорости, но не достаточно чтобы пройти выше. Для объектов User и Location был добавлено поле Visits, содержащие id всех посещений этого пользователя/локации.

[AterCattus](https://github.com/AterCattus/bicycle-mrhlc) Пошел дальше и сделал индексы для стран - насколько это помогает я не знаю.

## Другие оптимизации
1. Route запросов по одному символу.
2. Исключение аллокаций памяти в большинстве запросов.
2. Отключение сборщика мусора. На этапе загрузки сборщик мусора необорот настраивается более агрессивно (debug.SetGCPercent(50)) для того чтобы OOM не убил сервер, а после загрузки полное отключение сборщика мусора.

## Благодарность
Большое спасибо организатором за проведенный конкурс, он позволил разобраться во многих технологиях за очень короткий промежуток времени.