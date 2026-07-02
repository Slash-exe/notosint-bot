package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v3"
)

const (
	picError  = "⚠️ Ошибка отправки фото. Попробуйте позже или обратитесь в техподдержку."
	fileError = "⚠️ Ошибка отправки файла. Попробуйте позже или обратитесь в техподдержку."
)

type Tool struct {
	ImagePath string
	Caption   string
	FilePath  string
	FileName  string
}

func main() {
	godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id       BIGINT PRIMARY KEY,
		username TEXT,
		balance  NUMERIC DEFAULT 0
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS paid_invoices (
		invoice_id TEXT PRIMARY KEY
	)`)

	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	cryptoClient := resty.New().
		SetBaseURL("https://pay.crypt.bot/api").
		SetHeader("Crypto-Pay-API-Token", os.Getenv("CRYPTO_TOKEN"))

	// --- Хелперы ---

	getOrCreateUser := func(id int64, username string) {
		db.Exec(`INSERT INTO users (id, username, balance) VALUES ($1, $2, 0) ON CONFLICT DO NOTHING`, id, username)
	}

	getBalance := func(id int64) float64 {
		var balance float64
		db.QueryRow(`SELECT balance FROM users WHERE id = $1`, id).Scan(&balance)
		return balance
	}

	mainMenu := func() *tele.ReplyMarkup {
		menu := &tele.ReplyMarkup{}
		menu.Inline(
			menu.Row(
				menu.Data("🛠 Инструменты", "tools"),
				menu.URL("💬 Техподдержка", "t.me/Heirlom"),
			),
			menu.Row(menu.Data("💳 Пополнить баланс", "topup_info")),
		)
		return menu
	}

	sendMainMenu := func(c tele.Context) error {
		balance := getBalance(c.Sender().ID)
		photo := &tele.Photo{
			File: tele.FromDisk("files/images/well_cum.jpg"),
			Caption: fmt.Sprintf(
				"<b>notOSINT</b> — тулка в один клик\n\nПривет, @%s\n\n💰 <b>Баланс:</b> %.2f USDT\n\n🛠 Здесь ты найдёшь множество популярных инструментов.",
				c.Sender().Username, balance,
			),
		}
		return c.Send(photo, mainMenu(), tele.ModeHTML)
	}

	sendTool := func(c tele.Context, tool Tool) error {
		menu := &tele.ReplyMarkup{}
		menu.Inline(menu.Row(menu.Data("⬅️ Назад", "tools")))

		photo := &tele.Photo{
			File:    tele.FromDisk(tool.ImagePath),
			Caption: tool.Caption,
		}
		if err := c.Send(photo, menu, tele.ModeHTML); err != nil {
			log.Printf("Ошибка отправки фото (%s): %v", tool.ImagePath, err)
			c.Send(picError)
		}

		doc := &tele.Document{
			File:     tele.FromDisk(tool.FilePath),
			FileName: tool.FileName,
		}
		if err := c.Send(doc); err != nil {
			log.Printf("Ошибка отправки файла (%s): %v", tool.FilePath, err)
			c.Send(fileError)
		}
		return nil
	}

	// --- Хендлеры ---

	b.Handle("/start", func(c tele.Context) error {
		getOrCreateUser(c.Sender().ID, c.Sender().Username)
		return sendMainMenu(c)
	})

	b.Handle("\fback", func(c tele.Context) error {
		return sendMainMenu(c)
	})

	b.Handle("\ftools", func(c tele.Context) error {
		menu := &tele.ReplyMarkup{}
		menu.Inline(
			menu.Row(menu.Data("🔍 Sherlock", "sherlock"), menu.Data("👥 FSociety", "fsociety")),
			menu.Row(menu.Data("🐦 BlackBird", "blackbird"), menu.Data("🛠 Getting Tool", "gettool")),
			menu.Row(menu.Data("🎯 GHunt", "ghunt"), menu.Data("⛏️ The Harvester", "harvester")),
			menu.Row(menu.Data("🌌 Nexus", "nexus"), menu.Data("🚀 Arkada", "arkada")),
			menu.Row(menu.Data("⬅️ Назад", "back")),
		)
		return c.Send("Выбери тулку:", menu)
	})

	b.Handle("\fsherlock", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/sherlock.png",
			Caption:   "<b>Sherlock</b> — популярный open-source инструмент для OSINT. Ищет аккаунты человека по нику сразу на сотнях сайтов.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python sherlock.py</code>",
			FilePath:  "files/sherlock-master.rar",
			FileName:  "Sherlock Master.rar",
		})
	})

	b.Handle("\fblackbird", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/bird.jpg",
			Caption:   "<b>Blackbird</b> — продвинутый OSINT-инструмент. Ищет цифровые следы по нику <b>и</b> email одновременно.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python blackbird.py</code>",
			FilePath:  "files/blackbird-main.rar",
			FileName:  "BlackBird Main.rar",
		})
	})

	b.Handle("\fghunt", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/ghunt.png",
			Caption:   "<b>GHunt</b> — OSINT-инструмент для глубокого анализа цифрового следа через экосистему Google по Gmail или Gaia ID.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python main.py</code>",
			FilePath:  "files/GHunt-master.rar",
			FileName:  "GHunt Master.rar",
		})
	})

	b.Handle("\fnexus", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/Nexus.png",
			Caption:   "<b>Nexus</b> — консольный мультитул для кибербезопасности и поиска информации.\n\nФункции:\n• Username Search\n• Phone Lookup\n• Webhook Spammer\n• SMS Bomber\n\n<b>Инструкция</b>\nТребуется Python 3.10+. Запустите <b>setup.bat</b>.",
			FilePath:  "files/Nexus.rar",
			FileName:  "Nexus multitool.rar",
		})
	})

	b.Handle("\ffsociety", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/fsociety.jpg",
			Caption:   "<b>FSociety</b> — хакерский органайзер, который сам скачивает и запускает утилиты для пентеста.\n\nФункции:\n• Поиск SQL-инъекций\n• Генератор словарей для брутфорса\n• Bing-парсер уязвимостей\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python fsociety.py</code>",
			FilePath:  "files/fsociety-master.rar",
			FileName:  "FSociety.rar",
		})
	})

	b.Handle("\fgettool", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/GettingTool.png",
			Caption:   "<b>GettingTool</b> — мультитул для пробива информации, сетевого сканирования и работы с медиафайлами.\n\nФункции: leak, iplogger, whois, stealer, deepfake и другие.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python main.py</code>\n\n⚠️ <b>База данных не входит в комплект!</b>",
			FilePath:  "files/GETTINGTOOL.rar",
			FileName:  "Getting tool.rar",
		})
	})

	b.Handle("\fharvester", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/harvester.png",
			Caption:   "<b>theHarvester</b> — мультитул для OSINT и сетевого сканирования.\n\nФункции: iplogger, number, osint, stealer, deepfake, DDoS и другие.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python theharvester.py</code>",
			FilePath:  "files/theHarvester-master.rar",
			FileName:  "theHarvester.rar",
		})
	})

	b.Handle("\farkada", func(c tele.Context) error {
		return sendTool(c, Tool{
			ImagePath: "files/images/arkada.png",
			Caption:   "<b>Arkada New</b> — лаунчер для OSINT-пробива и деанона. Принимает номер, ник или почту и запускает соответствующие модули.\n\nФункции: dadata, pasport, auto, mail, nick и другие.\n\n<b>Инструкция</b>\nТребуется Python 3.10+.\n1. <code>cd c:/путь/к/папке/</code>\n2. <code>pip install -r requirements.txt</code>\n3. <code>python main.py</code>",
			FilePath:  "files/arkadanew.rar",
			FileName:  "Arkada.rar",
		})
	})

	b.Handle("\ftopup_info", func(c tele.Context) error {
		menu := &tele.ReplyMarkup{}
		menu.Inline(menu.Row(menu.Data("⬅️ Назад", "back")))
		return c.Send(`💳 <b>Пополнение баланса</b>

Для пополнения напиши команду:
<code>/topup 5</code>

Где <b>5</b> — сумма в USDT.

Минимальная сумма: <b>0.1 USDT</b>
Оплата через: <b>CryptoBot</b>`, menu, tele.ModeHTML)
	})

	b.Handle("/topup", func(c tele.Context) error {
		amountStr := c.Message().Payload
		if amountStr == "" {
			return c.Send("Использование: /topup 5\nВведи сумму в USDT после команды.")
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount < 0.1 {
			return c.Send("Введи корректную сумму (минимум 0.1), например: /topup 5")
		}

		var result map[string]interface{}
		_, err = cryptoClient.R().
			SetResult(&result).
			SetQueryParams(map[string]string{
				"asset":       "USDT",
				"amount":      amountStr,
				"description": "Пополнение баланса notOSINT",
				"payload":     fmt.Sprintf("%d", c.Sender().ID),
			}).
			Get("/createInvoice")

		if err != nil {
			log.Printf("Ошибка createInvoice: %v", err)
			return c.Send("Ошибка создания счёта, попробуй позже.")
		}

		if ok, _ := result["ok"].(bool); !ok {
			return c.Send("Ошибка CryptoBot. Попробуй позже.")
		}

		invoiceData := result["result"].(map[string]interface{})
		payURL := invoiceData["pay_url"].(string)
		invoiceID := fmt.Sprintf("%.0f", invoiceData["invoice_id"].(float64))

		menu := &tele.ReplyMarkup{}
		menu.Inline(
			menu.Row(menu.URL("💳 Оплатить", payURL)),
			menu.Row(menu.Data("✅ Проверить оплату", "check_"+invoiceID)),
			menu.Row(menu.Data("⬅️ Назад", "back")),
		)
		return c.Send(fmt.Sprintf("Счёт на <b>%.2f USDT</b> создан!\nПосле оплаты нажми ✅ Проверить оплату.", amount), menu, tele.ModeHTML)
	})

	b.Handle(tele.OnCallback, func(c tele.Context) error {
		data := strings.TrimPrefix(c.Callback().Data, "\f")
		if !strings.HasPrefix(data, "check_") {
			return nil
		}

		invoiceID := strings.TrimPrefix(data, "check_")

		var result map[string]interface{}
		cryptoClient.R().
			SetResult(&result).
			SetQueryParams(map[string]string{
				"invoice_ids": invoiceID,
				"status":      "paid",
			}).
			Get("/getInvoices")

		invoices, _ := result["result"].(map[string]interface{})
		items, _ := invoices["items"].([]interface{})

		if len(items) == 0 {
			return c.Respond(&tele.CallbackResponse{Text: "❌ Оплата не найдена. Попробуй чуть позже."})
		}

		invoice := items[0].(map[string]interface{})
		amount, _ := strconv.ParseFloat(invoice["amount"].(string), 64)

		var exists int
		db.QueryRow(`SELECT COUNT(*) FROM paid_invoices WHERE invoice_id = $1`, invoiceID).Scan(&exists)
		if exists > 0 {
			return c.Respond(&tele.CallbackResponse{Text: "❌ Этот счёт уже был зачислен."})
		}

		db.Exec(`INSERT INTO paid_invoices (invoice_id) VALUES ($1)`, invoiceID)
		db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, c.Sender().ID)

		return c.Respond(&tele.CallbackResponse{Text: fmt.Sprintf("✅ Баланс пополнен на %.2f USDT!", amount)})
	})

	b.Handle("/addbalance", func(c tele.Context) error {
		adminID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_ID"), 10, 64)
		if c.Sender().ID != adminID {
			return nil
		}

		amountStr := c.Message().Payload
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			return c.Send("Использование: /addbalance 10")
		}

		db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, c.Sender().ID)
		return c.Send(fmt.Sprintf("✅ Добавлено %.2f USDT", amount))
	})

	log.Println("Бот запущен...")
	b.Start()
}
