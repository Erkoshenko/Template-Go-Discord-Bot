package main

import (
	"bot/database"
	"bot/handlers"
	"bot/repository"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	dbConfig := database.PGConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.InitDB(dbConfig)
	if err != nil {
		log.Fatalln(err)
	}

	repos := repository.NewRepositories(db)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	repos.Users.CreateTableIfNotExists(ctx)

	h := handlers.NewHandler(repos)

	sess, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("Ошибка при создании сессий: " + err.Error())
		return
	}

	h.RegisterAll(sess)

	err = sess.Open()
	if err != nil {
		fmt.Println("Ошибка при открытии сессий: " + err.Error())
		return
	}

	defer sess.Close()

	fmt.Println("Нажмите ctrl+c чтобы убить")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
