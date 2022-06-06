# Добавление в todo-app функции кэширования в Redis

## Информация по реализации Redis
При GET запросах данные кэшируются в хэш-таблицу Redis.

Структура хэш-таблицы:

    - HKEYS: user:'userId' - ключи таблицы;
    - FIELD: lists (getAllLists), list:'id' (getListById), items:'listId' (getAllItems), item:'id' (getItemById) - поля ключей user:'userId' хэш-таблицы.
В скобках рядом с полями указаны handler функции GET, которST кэшируют данные в это поле.

Handler функция createList удаляет из хэш-таблицы поле lists.
Handler функции updateList, deleteList, updateItem, deleteItem удаляют из хэш-таблицы весь ключ user:'userId'.
