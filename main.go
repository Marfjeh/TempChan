package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Main Settings
var (
	Token string = "NTc1NjExMTA2ODcwMDM0NDQy.XNKdng.8XBWbmguV3F5Ff6ngpKbcAcurDw"
	Prefix string = "]"
)

type ChannelJson struct {
	CategoryID string
	TextID	string
	VoiceID string
	GuildID string
	OwnerID string
}

func main() {
	fmt.Println("MarfBOT:Go Kernel Initialized." +
					 "\n------------------------------")

	//New Discord session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating discord session. Discord down Lulz", err)
		return
	}

	//Register Handlers
	dg.AddHandler(MessageCreate)


	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection. ", err)
		return
	}

	fmt.Println("MarfBOT is running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == Prefix + "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == Prefix + "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == Prefix + "help" {
		s.ChannelMessageSend(m.ChannelID, "To create Channels: ]cc <name>")
		s.ChannelMessageSend(m.ChannelID, "Setting a Limit ]climit <number of people>")
	}

	if m.Content == Prefix + "cc" {
		s.ChannelMessageSend(m.ChannelID, "Error missing channel name!")
	}

	if m.Content == Prefix + "exit" && m.Author.ID == "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Shutting down.")
		os.Exit(0)
	} else if m.Content == Prefix + "exit" && m.Author.ID != "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Really?")
	}

	if strings.HasPrefix(m.Content, Prefix + "cc ") {

		tempchan, err := s.ChannelMessageSend(m.ChannelID, "Creating temporay channels for you...")

		channelCategory, err 	:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildCategory)
		channelText, err 		:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildText)
		channelVoice, err 		:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildVoice)

		_, err = s.ChannelEditComplex(channelText.ID, &discordgo.ChannelEdit{ParentID: channelCategory.ID, Topic: "Created by: " + m.Author.Username + ". This is a temporary channel."})
		_, err = s.ChannelEditComplex(channelVoice.ID, &discordgo.ChannelEdit{ParentID: channelCategory.ID})

		_, err = s.ChannelMessageSend(channelText.ID, m.Author.Mention() + " This is your newly created channel. " +
			"Join Witin 30 seconds, or this channel will be deleted.")

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Sorry, I was unable to create the channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		} else {
			s.ChannelMessageDelete(m.ChannelID, tempchan.ID)
			s.ChannelMessageSend(m.ChannelID, "Channel created. please join the channel within 30 seconds, or it will be deleted")
		}
	}
}
