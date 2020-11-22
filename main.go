package main
//Welcome to the horrible mess of my second go program.
// <shitpost>
// Rawr xD this verriw nice firwst progewm OwO
// Bye XD
// </shitpost>

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Config struct {
	Token 				string `json:"Token"`
	Prefix 				string `json:"Prefix"`
	Owner_id 			string `json:"Owner_id"`
	database_host 		string `json:"database_host"`
	database_port 		string `json:"database_port"`
	database_connect	string `json:"database_connect"`
	database_socket		string `json:"database_socket"`
	database_user 		string `json:"database_user"`
	database_password 	string `json:"database_password"`
	database_table 		string `json:"database_table"`
}

type Channel struct {
	id 					int
	name 				string
	owner_id 			string
	category_id 		string
	voicechannel_id 	string
	textchannel_id 		string
	settingmessage_id 	string
	options 			string
}

var DB *sql.DB

func main() {
	// Logo
	//goland:noinspection GoPrintFunctions to make my IDE not shit itself
	fmt.Println("  _______                    _____ _                   \n" +
		" |__   __|                  / ____| |                 \n" +
		"    | | ___ _ __ ___  _ __ | |    | |__   __ _ _ __   \n" +
		"    | |/ _ \\ '_ ` _ \\| '_ \\| |    | '_ \\ / _` | '_ \\  \n" +
		"    | |  __/ | | | | | |_) | |____| | | | (_| | | | | \n" +
		"    |_|\\___|_| |_| |_| .__/ \\_____|_| |_|\\__,_|_| |_| \n" +
		"                     | |                              \n" +
		"                     |_|                                \n")

	//Read config file
	fmt.Println("INIT: Config file")
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Unable to read the config file!")
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&Config); err != nil {
		fmt.Println("Unable to parse Json.")
		panic(err)
	}
	fmt.Println("Config read successfully.")

	//Init Database
	//TODO: implement Unix socket connections
	fmt.Println("INIT: Database in TCP mode")

	db, err := sql.Open("mysql", Config.database_user + ":" + Config.database_password + "@(" + Config.database_host + ":" + Config.database_port + " )/" + Config.database_table)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Databast init success!")
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	//Setting global variable...
	DB = db

	//New Discord session
	fmt.Println("INIT: Discord Connection")
	dg, err := discordgo.New("Bot " + Config.Token)
	if err != nil {
		fmt.Println("Error creating discord session. Discord down Lulz", err)
		return
	} else {
		fmt.Println("Discord session created.")
	}

	//Register Handlers
	fmt.Println("Registering handlers...")
	dg.AddHandler(MessageCreate)
	dg.AddHandler(MessageReactions)

	fmt.Println("Setting Intents...")
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection. ", err)
		return
	}

	fmt.Println("Done. TempChan is running. Logged in as: " + dg.State.User.Username + " Prefix set to: " + Config.Prefix)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	db.Close()
}

func MessageReactions(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//Prevent bot reactions to trigger itself.
	if r.UserID == s.State.User.ID {
		return
	}

	//Check the message where the reaction is made, if someone reacted to a message not from the bot, just return it and dont do a SQL query.
	OrigMessage, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if OrigMessage.Author.ID != s.State.User.ID {
		return
	}

	var channel = Channel{}
	fmt.Println(channel)
	_ = DB.QueryRow("SELECT * from channels WHERE settingmessage_id = "+r.MessageID).Scan(&channel.id, &channel.name, &channel.owner_id, &channel.category_id, &channel.voicechannel_id, &channel.textchannel_id, &channel.settingmessage_id, &channel.options)
	fmt.Println(channel.settingmessage_id)

	if r.MessageReaction.Emoji.Name == "❌" && r.MessageReaction.UserID == "218310787289186304" && r.MessageID == channel.settingmessage_id {
		s.ChannelDelete(channel.voicechannel_id)
		s.ChannelDelete(channel.textchannel_id)
		s.ChannelDelete(channel.category_id)

		stmt, _ := DB.Prepare("DELETE from channels where settingmessage_id = ?")
		_, _ = stmt.Exec(r.MessageID)
	}

	if r.MessageReaction.Emoji.Name == "🔞" && r.MessageReaction.UserID != s.State.User.ID && r.MessageID == channel.settingmessage_id  {
		s.ChannelEditComplex(channel.textchannel_id, &discordgo.ChannelEdit{NSFW: true})
		s.MessageReactionRemove(channel.textchannel_id, channel.settingmessage_id, "🔞", r.UserID);
		s.MessageReactionRemove(channel.textchannel_id, channel.settingmessage_id, "🔞", s.State.User.ID)
	}

	if err != nil {
		fmt.Println(err)
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
	if m.Content == Config.Prefix + "help" {
		s.ChannelMessageSend(m.ChannelID, 	"To create Channels: ]cc <name>\n" +
													"Setting a Limit ]climit <number of people>")
	}

	if m.Content == Config.Prefix + "channel delete all" && m.Author.ID == "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Deleting all channels currently in database.")
		channelrows, err := DB.Query("SELECT * from channels")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error! Was unable to get data from database.")
			fmt.Println(err)
		}

		var channel = Channel{}
		for channelrows.Next() {
			_ = channelrows.Scan(&channel.id, &channel.name, &channel.owner_id, &channel.category_id, &channel.voicechannel_id, &channel.textchannel_id, &channel.settingmessage_id, &channel.options)
			fmt.Println("Deleting: " + channel.name)
			s.ChannelDelete(channel.voicechannel_id)
			s.ChannelDelete(channel.textchannel_id)
			s.ChannelDelete(channel.category_id)

			stmt, _ := DB.Prepare("DELETE from channels where id = ?")
			_, _ = stmt.Exec(channel.id)
		}
	}

	if m.Content == Config.Prefix + "exit" && m.Author.ID == "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Shutting down.")
		fmt.Println("Got exit call from discord command.")
		s.Close()
		DB.Close()
		os.Exit(0)
	} else if m.Content == Config.Prefix + "exit" && m.Author.ID != "218310787289186304" {
		s.ChannelMessageSend(m.ChannelID, "Watch this: <https://www.youtube.com/watch?v=dQw4w9WgXcQ>")
	}

	if m.Content == Config.Prefix + "cc" {
		s.ChannelMessageSend(m.ChannelID, "Error missing channel name!")
	}

	//todo: make a command framework?
	if strings.HasPrefix(m.Content, Config.Prefix + "cc ") {

		tempchan, err := s.ChannelMessageSend(m.ChannelID, "Creating temporay channels for you...")

		//Create the first category first.
		channelCategory, err 	:= s.GuildChannelCreate(m.GuildID, strings.Trim(m.Content, Config.Prefix + "cc "), discordgo.ChannelTypeGuildCategory)

		channelText, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name:                 strings.Trim(m.Content, Config.Prefix + "cc "),
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
			Name: strings.Trim(m.Content, Config.Prefix + "cc "),
			Type: discordgo.ChannelTypeGuildVoice,
			ParentID: channelCategory.ID,
		})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the voice channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		}

		channelSettings, err := s.ChannelMessageSend(channelText.ID, "```*** Channel Settings ***```\nThis is your newly created channel.\n" +
			"The channel owner is: " + m.Author.Mention() + "\n" +
			"The Channel owner and staff are able to change the settings via the emojis. You are also still required to obey the server rules!\n\n" +
			"❌ = Delete the channel\n" +
			"🔒 = Disable this Chat channel.\n" +
			"🔞 = Set Chat as NSFW\n\n" +
			"For more commands enter ]channel help\n\n\n" +
			"⚠ Note: Staff of this server is also able to change every setting.")

		//Saving TO DB
		fmt.Println("Saving to DB...")
		stmt, err := DB.Prepare("INSERT INTO `channels` (`name`, `owner_id`, `category_id`, `voicechannel_id`, `textchannel_id`, `settingmessage_id`, `options`)  VALUES (?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(strings.Trim(m.Content, Config.Prefix + "cc "), m.Author.ID, channelCategory.ID, channelVoice.ID, channelText.ID, channelSettings.ID, nil)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "⚠ Sorry, I was unable to create the channel. <@218310787289186304> Logged error to console.")
			fmt.Println("Error: ", err)
		} else {
			s.ChannelMessageDelete(m.ChannelID, tempchan.ID)
			s.ChannelMessageSend(m.ChannelID, "Channel created. please join the channel within 30 seconds, or it will be deleted")
			fmt.Println("Created a new temporary channel to watch on... ")
		}

		fmt.Println("Adding reactions...")
		err = s.ChannelMessagePin(channelText.ID, channelSettings.ID)
		err = s.MessageReactionAdd(channelText.ID, channelSettings.ID, "❌")
		err = s.MessageReactionAdd(channelText.ID, channelSettings.ID, "🔒")
		err = s.MessageReactionAdd(channelText.ID, channelSettings.ID, "🔞")
		if err != nil {
			fmt.Println("Unable to add a reaction...")
		}

		fmt.Println("Channel creation Done.")
	}
}
