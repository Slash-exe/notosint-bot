# notOSINT Bot

Telegram-бот для быстрого доступа к популярным OSINT-инструментам с системой баланса и оплатой через CryptoBot.

## Функции

- 📦 Каталог OSINT-инструментов с описаниями и файлами
- 💳 Пополнение баланса через [CryptoBot](https://t.me/CryptoBot) (USDT)
- 🗄️ PostgreSQL база данных пользователей
- 🔒 Защита от двойного зачисления средств
- ⚙️ Команда администратора для выдачи баланса

## Инструменты в боте

| Инструмент | Описание |
|---|---|
| Sherlock | Поиск аккаунтов по нику на 300+ сайтах |
| Blackbird | Поиск по нику и email |
| GHunt | Анализ Google-аккаунта по Gmail |
| Nexus | Мультитул: поиск, пробив номера, SMS Bomber |
| FSociety | Хакерский органайзер и генератор словарей |
| GettingTool | Пробив, сканирование, deepfake |
| theHarvester | OSINT и сетевое сканирование |
| Arkada New | Лаунчер для пробива по номеру, нику, почте |

## Стек

- **Go** + [telebot v3](https://github.com/tucnak/telebot)
- **PostgreSQL** через [lib/pq](https://github.com/lib/pq)
- **CryptoBot API** через [resty](https://github.com/go-resty/resty)

## Установка

### 1. Клонируй репозиторий

```bash
git clone https://github.com/твой-ник/notosint-bot
cd notosint-bot
```

### 2. Создай базу данных PostgreSQL

```bash
psql -U postgres
CREATE DATABASE notosint;
\q
```

### 3. Создай `.env` файл

```env
BOT_TOKEN=токен_от_BotFather
CRYPTO_TOKEN=токен_от_CryptoBot
TELEGRAM_ID=твой_telegram_id
DATABASE_URL=postgres://postgres:пароль@localhost:5432/notosint
```

### 4. Добавь файлы

Положи инструменты в папку `files/`, картинки в `files/images/`.

### 5. Запусти

```bash
go run main.go
```

## Структура проекта

```
notosint-bot/
├── main.go
├── .env
├── .gitignore
└── files/
    ├── images/
    │   ├── welcome.jpg
    │   ├── sherlock.png
    │   └── ...
    ├── sherlock-master.rar
    └── ...
```

## Команды

| Команда | Описание |
|---|---|
| `/start` | Главное меню |
| `/topup 5` | Создать счёт на 5 USDT |
| `/addbalance 10` | Выдать себе 10 USDT (только для админа) |

## Лицензия

MIT
