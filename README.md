Сервис для автоматического создания и размещения рекламных постов в вк группе, а так же может быть легко модифицирован для мультипостинга в разные вк группы.
При желании может быть настроен для размещения постов в сообществе, так же настраивается для создание постов из готовых фотографий (не по ссылке в табличке, а из фотографий на жестком диске)

Создает пост используя готовые крео, а так же размещает его в группе и оставляет ссылки в комментариях
Работает с гугл таблицей (в будущем SQLdb, но пока что таблица используется для того, что бы рекламодатель мог оставить нужные данные в удобном формате) Таблица: https://docs.google.com/spreadsheets/d/18yE-2O7boF0cWboMpkNNlDQsXNp2_l9DQo6HUiSqPLY/edit?usp=sharing

Как пользоваться таблицей:
1.   Лист Data служит для удобного просмотра всех данных и не несет ни какого функционала
2.   Лист Creo служит для хранения информации о Крео. Содержит следующие поля:
   1.   Номер крео - ID, по которому крео выбирается
   2.   Ссылка на фото 1 и 2 - ссылки, по которым качаются фотографии
   3.   Текст - текст самого поста
   4.   Авг охват - средний охват у конкретного крео за все время
   5.   Последнее использование - дата, в которую данное крео было в последний раз использовано
3.   Листы для рекламодателей, которые содержат следующие поля:
   1.   Дата - Дата и время поста
   2. 	Номер поста - Номер поста в серии (1, если единичный пост)
   3. 	PostID - ID поста в ВК
   4. 	Тгк - Ссылка на рекламируемый ТГК
   5. 	Ссылка на пост - Ссылка на пост, которая появляется после публикации
   6. 	Крео - ID крео, по которому автоматически будет создоваться пост
   7. 	Охват через сутки - Статистика поста по просмотрам через сутки
После заполнения полей "Номер пост", "ТГК" и "Крео" создается горутина для каждого поста (строчка листа рекламодателя, которая содержит номер поста, ссылку на пост и ID крео) на сегодняшний (workerToday) и вчерашний день (workerYesterday).
workerToday создает пост, используя данные о кре, а так же публикует его на стене сообщества в нужное время. Для этого он использует токен фейковой страницы и токен администратора сообщества (тоже может быть фейковой страницей с ролью "модератор"). Существует два варианта создания поста:
1. Если в сообществе не требуется 2fa для размещения записей или есть токен с подключенной 2fa - workerToday пригласит фейковый аккаунт в сообщество -> выдаст роли для создания поста -> сделает пост -> закроет комментарии -> выгонит фейковый аккаунт из сообщества
2. Если в сообществе требуется 2fa, а токенов таких аккаунтов нет - workerToday пригласит пригласит фейковый аккаунт в сообщество -> сделает пост от имении сообщества, используя токен администратора (не банится за просто создание поста без ссылок) -> напишет комментарий от имени фейкового пользователя -> закроет комментарии -> выгонит фейкового юзера из сообщества
После создания поста workerToday добавит в таблицу ссылку на пост и PostID
workerYesterday получает статистику по постам ровно через сутки после их публикации и добавляет в таблицу

