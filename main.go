package main
//Welcome to the horrible mess of my second go program.
// <shitpost>
// Rawr xD this verriw nice firwst progewm OwO
// Bye XD
// </shitpost>

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Main Settings.
//TODO: Settings.json and delete tempShit so we can track more channels then one.
var (
	Token string = "NTc1NjExMTA2ODcwMDM0NDQy.XNKdng.8XBWbmguV3F5Ff6ngpKbcAcurDw"
	Prefix string = "]"
	tempCate string
	tempText string
	tempVoice string
	tempSettings string
)

type ChannelJson struct {
	CategoryID string
	TextID	string
	VoiceID string
	GuildID string
	OwnerID string
}

func main() {
	//Todo: Rename it
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
	dg.AddHandler(MessageReactions)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

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

func MessageReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//fmt.Println(r.MessageReaction.Emoji.Name == "❌")
	if r.MessageReaction.Emoji.Name == "❌" && r.MessageReaction.UserID == "218310787289186304" && r.MessageID == tempSettings {
		s.ChannelDelete(tempVoice)
		s.ChannelDelete(tempText)
		s.ChannelDelete(tempCate)
	}
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
		//channelText, err 		:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildText)
		//channelVoice, err 	:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildVoice)

		channelText, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:                 strings.Trim(m.Content, Prefix + "cc "),
			Type:                 discordgo.ChannelTypeGuildText,
			Topic:                "Created by: " + m.Author.Username + ". This is a temporary channel.",
			Position:             0,
			ParentID:             channelCategory.ID,
			NSFW:                 false,
		})

		// Create the text channel
		channelVoice, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name: strings.Trim(m.Content, Prefix + "cc "),
			Type: discordgo.ChannelTypeGuildVoice,
			ParentID: channelCategory.ID,
		})

		if err != nil {
			log.Printf("Cannot create voice channel: %v\n", err)
		}

		tempCate = channelCategory.ID
		tempText = channelText.ID
		tempVoice = channelVoice.ID

		//_, err = s.ChannelEditComplex(channelText.ID, &discordgo.ChannelEdit{ParentID: channelCategory.ID, Topic: "Created by: " + m.Author.Username + ". This is a temporary channel."})
		//_, err = s.ChannelEditComplex(channelVoice.ID, &discordgo.ChannelEdit{ParentID: channelCategory.ID})

		channelSettings, err := s.ChannelMessageSend(channelText.ID, "```*** Channel Settings ***```\nThis is your newly created channel.\n" +
			"The channel owner is: " + m.Author.Mention() + "\n" +
			"Only that person* can perform options by reacting to the emoji's below.\n\n" +
			"❌ = Delete the channel\n" +
			"🔒 = Disable this Chat channel.\n\n" +
			"For more commands enter ]channel help\n\n\n" +
			"⚠ Note: Staff of this server is also able to change every setting.")
		tempSettings = channelSettings.ID;
		s.ChannelMessagePin(channelText.ID, channelSettings.ID)

		s.MessageReactionAdd(channelText.ID, channelSettings.ID, "❌")
		s.MessageReactionAdd(channelText.ID, channelSettings.ID, "🔒")


		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		} else {
			s.ChannelMessageDelete(m.ChannelID, tempchan.ID)
			s.ChannelMessageSend(m.ChannelID, "Channel created. please join the channel within 30 seconds, or it will be deleted")
		}
	}
}
