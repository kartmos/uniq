В `README.md` для программы, реализующей функционал команды `uniq`, стоит подробно описать назначение программы, ее функции, установку, примеры использования и возможные флаги. Вот шаблон, который вы можете использовать для своего `README.md`.

---

# `uniq` - Утилита для работы с уникальными строками

`uniq` — это утилита командной строки, предназначенная для фильтрации дублирующихся строк из текстовых данных. Программа читает входные данные, сравнивает последовательные строки и выводит либо уникальные строки, либо только дубликаты, в зависимости от заданных флагов.

## Основные возможности

- **Фильтрация дублирующихся строк**: программа может выводить только уникальные строки или, наоборот, только дубликаты.
- **Подсчет количества повторений**: можно выводить количество вхождений каждой уникальной строки.
- **Регистронезависимое сравнение строк**: опция для игнорирования регистра при сравнении строк.
- **Пропуск начальных полей**: возможность игнорировать указанное количество начальных слов при сравнении строк.

## Установка

Для запуска программы, клонируйте репозиторий и соберите исполняемый файл:

```bash
git clone <URL вашего репозитория>
cd uniq
go build -o uniq .
```

## Использование

Запустите программу с помощью командной строки, передав входной файл или ввод через `stdin`:

```bash
./uniq [флаги] [входной_файл] > [выходной_файл]
```

### Доступные флаги

- `-c`: Подсчитать количество вхождений каждой строки.
- `-d`: Вывести только дубликаты (повторяющиеся строки).
- `-u`: Вывести только уникальные строки (которые не повторяются).
- `-i`: Игнорировать регистр при сравнении строк.
- `-s N`: Пропустить первые `N` слов (полей) при сравнении строк.

### Примеры

1. **Вывести только дублирующиеся строки**:

    ```bash
    ./uniq -d input.txt
    ```

2. **Подсчитать количество вхождений строк с игнорированием регистра**:

    ```bash
    ./uniq -c -i input.txt
    ```

3. **Пропустить первые два слова и вывести уникальные строки**:

    ```bash
    ./uniq -u -s 2 input.txt
    ```

## Примеры входных данных

Создайте файл `input.txt` со следующим содержимым для тестирования:

```
apple
banana
apple
orange
banana
banana
apple
grape
grape
kiwi
orange
```

## Пример вывода

Использование `uniq -d input.txt` с указанным выше файлом даст следующий результат:

```
apple
banana
grape
```

Этот вывод покажет только те строки, которые повторяются.

---