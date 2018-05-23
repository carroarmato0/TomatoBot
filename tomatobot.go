package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ini/ini"
)

func main() {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + cfg.Section("").Key("TOKEN").String())

	if err != nil {
		log.Println(err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// Plugin interface
type Plugin interface {
	process()
}

// Ping plugin
type Ping struct{}

func (p Ping) process(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

// Time plugin
type Time struct{}

func (t Time) process(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "time" || m.Content == "tijd" {
		t := time.Now()
		time := t.Format(time.RFC1123)
		s.ChannelMessageSend(m.ChannelID, time)
	}
}

// Contains - Check if string is in slice
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Greet plugin
type Greet struct{}

func (g Greet) process(s *discordgo.Session, m *discordgo.MessageCreate) {
	dutch := []string{
		"hallo",
		"dag",
		"dag iedereen",
	}
	english := []string{
		"hello",
		"hi",
		"greetings",
		"evening",
	}

	message := strings.ToLower(m.Content)

	var replyDutch map[string]string
	replyDutch = make(map[string]string)
	replyDutch["morning"] = "Goeie morgen"
	replyDutch["noon"] = "Goeie middag"
	replyDutch["evening"] = "Goeie namiddag"
	replyDutch["night"] = "Goeie avond"
	var replyEnglish map[string]string
	replyEnglish = make(map[string]string)
	replyEnglish["morning"] = "Good morning"
	replyEnglish["noon"] = "Good afternoon"
	replyEnglish["evening"] = "Good evening"
	replyEnglish["night"] = "Good night"

	currentHour := time.Now().Hour()

	var language string
	if contains(dutch, message) {
		language = "dutch"
	} else if contains(english, message) {
		language = "english"
	} else {
		return
	}

	var timeOfDay string
	if currentHour >= 0 && currentHour < 12 {
		timeOfDay = "morning"
	} else if currentHour >= 12 && currentHour < 14 {
		timeOfDay = "noon"
	} else if currentHour >= 14 && currentHour < 18 {
		timeOfDay = "evening"
	} else {
		timeOfDay = "night"
	}

	var reply string
	if language == "dutch" {
		reply = replyDutch[timeOfDay]
	} else if language == "english" {
		reply = replyEnglish[timeOfDay]
	}

	s.ChannelMessageSend(m.ChannelID, reply+" <@"+m.Author.ID+">!")

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	new(Ping).process(s, m)
	new(Time).process(s, m)
	new(Greet).process(s, m)
}
