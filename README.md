notOSINT Bot

Telegram-бот для быстрого доступа к популярным OSINT-инструментам с системой баланса и оплатой через CryptoBot.

Функции


📦 Каталог OSINT-инструментов с описаниями и файлами
💳 Пополнение баланса через CryptoBot (USDT)
🗄️ PostgreSQL база данных пользователей
🔒 Защита от двойного зачисления средств
⚙️ Команда администратора для выдачи баланса


Инструменты в боте

ИнструментОписаниеSherlockПоиск аккаунтов по нику на 300+ сайтахBlackbirdПоиск по нику и emailGHuntАнализ Google-аккаунта по GmailNexusМультитул: поиск, пробив номера, SMS BomberFSocietyХакерский органайзер и генератор словарейGettingToolПробив, сканирование, deepfaketheHarvesterOSINT и сетевое сканированиеArkada NewЛаунчер для пробива по номеру, нику, почте

Стек


Go + telebot v3
PostgreSQL через lib/pq
CryptoBot API через resty


Установка

1. Клонируй репозиторий

bashgit clone https://github.com/твой-ник/notosint-bot
cd notosint-bot

2. Создай базу данных PostgreSQL

bashpsql -U postgres
CREATE DATABASE notosint;
\q

3. Создай .env файл

envBOT_TOKEN=токен_от_BotFather
CRYPTO_TOKEN=токен_от_CryptoBot
TELEGRAM_ID=твой_telegram_id
DATABASE_URL=postgres://postgres:пароль@localhost:5432/notosint

4. Добавь файлы

Положи инструменты в папку files/, картинки в files/images/.

5. Запусти

bashgo run main.go

Структура проекта

notosint-bot/
├── main.go
├── .env
├── .gitignore
└── files/
    ├── images/
    │   ├── well_cum.jpg
    │   ├── sherlock.png
    │   └── ...
    ├── sherlock-master.rar
    └── ...

Команды

КомандаОписание/startГлавное меню/topup 5Создать счёт на 5 USDT/addbalance 10Выдать себе 10 USDT (только для админа)
