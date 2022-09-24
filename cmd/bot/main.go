package main

import (
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/handler"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/repository"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/services"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error load config file: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error load env file: %s", err.Error())
	}

	cfg := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbName"),
		SSLMode:  viper.GetString("db.sslMode"),
	}

	logrus.Infof("db config: %v", cfg)

	db, err := repository.NewPostgresDB(cfg)

	if err != nil {
		logrus.Fatalf("error load database: %s", err.Error())
	}

	db.SetMaxIdleConns(50)

	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Fatalf("error close db: %s", err.Error())
		}
	}()
	// injection dependency
	svs := services.New()
	repos := repository.New(db)
	usecaseCnf := usecases.UseCaseConfig{
		Repository: repos,
		Service:    svs,
	}
	usecases := usecases.NewUseCase(usecaseCnf)
	handlers := handler.NewHandler(usecases, svs)
	handlers.InitHandlers()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		logrus.Fatalf("error created bot: %s", err.Error())
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, handlers)
	if err := telegramBot.Start(); err != nil {
		logrus.Debugf("error send msg: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
