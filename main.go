package main
//Welcome to the horrible mess of my second go program.
// <shitpost>
// Rawr xD this verriw nice firwst progewm OwO
// Bye XD
// </shitpost>

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
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
	tempCate string //Deprecated
	tempText string //Deprecated
	tempVoice string //Deprecated
	tempSettings string //Deprecated
)

type ChannelJson struct {
	CategoryID string
	TextID	string
	VoiceID string
	GuildID string
	OwnerID string
}

func main() {
	// Logo
	//goland:noinspection GoPrintFunctions
	fmt.Println("  _______                    _____ _                   \n" +
		" |__   __|                  / ____| |                 \n" +
		"    | | ___ _ __ ___  _ __ | |    | |__   __ _ _ __   \n" +
		"    | |/ _ \\ '_ ` _ \\| '_ \\| |    | '_ \\ / _` | '_ \\  \n" +
		"    | |  __/ | | | | | |_) | |____| | | | (_| | | | | \n" +
		"    |_|\\___|_| |_| |_| .__/ \\_____|_| |_|\\__,_|_| |_| \n" +
		"                     | |                              \n" +
		"                     |_|                                \n")

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

	fmt.Println("TempChan is running. Logged in as: " + dg.State.User.Username)

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

	if r.MessageReaction.Emoji.Name == "🔞" && r.MessageReaction.UserID != s.State.User.ID && r.MessageID == tempSettings {
		s.ChannelEditComplex(tempText, &discordgo.ChannelEdit{NSFW: true})
		s.MessageReactionRemove(tempText, tempSettings, "🔞", r.UserID);
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "iOS")  {
		_, err := s.ChannelMessageSend(m.ChannelID, "> " + m.Content + "\n" + m.Author.Mention() + " Its `Ios` Not `iOS` :)")
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Show a simple help list.
	if m.Content == Prefix + "help" {
		s.ChannelMessageSend(m.ChannelID, 	"To create Channels: ]cc <name>\n" +
													"Setting a Limit ]climit <number of people>")
	}

	if m.Content == Prefix + "exit" && m.Author.ID == "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Shutting down.")
		fmt.Println("Got exit call from discord command.")
		s.Close()
		os.Exit(0)
	} else if m.Content == Prefix + "exit" && m.Author.ID != "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Watch this: <https://www.youtube.com/watch?v=dQw4w9WgXcQ>")
	}

	if m.Content == Prefix + "cc" {
		s.ChannelMessageSend(m.ChannelID, "Error missing channel name!")
	}

	//todo: make a command framework?
	if strings.HasPrefix(m.Content, Prefix + "cc ") {

		tempchan, err := s.ChannelMessageSend(m.ChannelID, "Creating temporay channels for you...")

		//Create the first category first.
		channelCategory, err 	:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Prefix + "cc "), discordgo.ChannelTypeGuildCategory)

		channelText, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:                 strings.Trim(m.Content, Prefix + "cc "),
			Type:                 discordgo.ChannelTypeGuildText,
			Topic:                "Created by: " + m.Author.Username + ". This is a temporary channel.",
			Position:             0,
			ParentID:             channelCategory.ID,
			NSFW:                 false,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the text channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		}

		// Create the text channel
		channelVoice, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name: strings.Trim(m.Content, Prefix + "cc "),
			Type: discordgo.ChannelTypeGuildVoice,
			ParentID: channelCategory.ID,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the voice channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		}

		//todo: make a manager something like that to keep track of channels.
		tempCate = channelCategory.ID
		tempText = channelText.ID
		tempVoice = channelVoice.ID

		channelSettings, err := s.ChannelMessageSend(channelText.ID, "```*** Channel Settings ***```\nThis is your newly created channel.\n" +
			"The channel owner is: " + m.Author.Mention() + "\n" +
			"Only that person* can perform options by reacting to the emoji's below.\n\n" +
			"❌ = Delete the channel\n" +
			"🔒 = Disable this Chat channel.\n" +
			"🔞 = Set Chat as NSFW\n\n"+
			"For more commands enter ]channel help\n\n\n" +
			"⚠ Note: Staff of this server is also able to change every setting.")
		tempSettings = channelSettings.ID;
		s.ChannelMessagePin(channelText.ID, channelSettings.ID)

		s.MessageReactionAdd(channelText.ID, channelSettings.ID, "❌")
		s.MessageReactionAdd(channelText.ID, channelSettings.ID, "🔒")
		s.MessageReactionAdd(channelText.ID, channelSettings.ID, "🔞")


		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		} else {
			s.ChannelMessageDelete(m.ChannelID, tempchan.ID)
			s.ChannelMessageSend(m.ChannelID, "Channel created. please join the channel within 30 seconds, or it will be deleted")
			fmt.Println("Created a new temporary channel to watch on... Guild ID: " + m.GuildID + " Category ID: " + tempCate + " TextChannel ID: " + tempText + " VoiceChannel ID: " + tempVoice)
		}
	}
}
