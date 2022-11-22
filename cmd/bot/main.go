package main

import (
	"context"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/handler"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/logger"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/repository"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/services"
	"github.com/evgeniy-krivenko/particius-vpn-bot/internal/usecases"
	"github.com/evgeniy-krivenko/particius-vpn-bot/pkg/telegram"
	"github.com/evgeniy-krivenko/vpn-api/gen/conn_service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	lgr := logrus.New()
	lgr.SetFormatter(new(logrus.JSONFormatter))

	log := logger.NewLogrusLogger(lgr)

	if err := initConfig(); err != nil {
		log.Fatalf("error load config file: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error load env file: %s", err.Error())
	}

	cfg := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbName"),
		SSLMode:  viper.GetString("db.sslMode"),
	}

	log.Infof("db config: %v", cfg)

	db, err := repository.NewPostgresDB(cfg)

	if err != nil {
		log.Fatalf("error load database: %s", err.Error())
	}

	db.SetMaxIdleConns(50)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("error close db: %s", err.Error())
		}
	}()
	// injection dependency
	// получаем список серверов
	// в цикле создаем подключения и клиенты и добавляем их в мапу по ключу, например NL: client
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	conn, err := grpc.DialContext(ctxWithTimeout, "localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("error connect to grpc server %s", err.Error())
	}
	defer conn.Close()

	cs := conn_service.NewConnectionServiceClient(conn)

	// inject grpc service to services
	svs := services.New()
	repos := repository.New(db)
	usecaseCnf := usecases.UseCaseConfig{
		Repository: repos,
		Grpc:       cs,
		Log:        log,
	}
	uc := usecases.NewUseCase(usecaseCnf)
	handlers := handler.NewHandler(uc, svs, log)
	handlers.InitHandlers()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_API_KEY"))
	if err != nil {
		log.Fatalf("error created bot: %s", err.Error())
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot, handlers)
	if err := telegramBot.Start(ctx); err != nil {
		log.Debugf("error send msg: %s", err.Error())
	}

	<-ctx.Done()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
