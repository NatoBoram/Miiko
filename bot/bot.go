package bot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/NatoBoram/Go-Miiko/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	// DB : Connection to the database.
	DB *sql.DB

	// Me : The bot itself.
	Me *discordgo.User

	// Master : UserID of the bot's master.
	Master *discordgo.User
)

// Start : Starts the bot.
func Start(db *sql.DB, session *discordgo.Session, master string) error {

	// Database
	DB = db

	// Myself
	user, err := session.User("@me")
	if err != nil {
		fmt.Println("Couldn't get myself.")
		return err
	}
	Me = user

	// Master
	user, err = session.User(master)
	if err != nil {
		fmt.Println("Couldn't recognize my master.")
		return err
	}
	Master = user

	// Hey, listen!
	//session.AddHandler(messageHandler)
	session.AddHandler(reactHandler)
	session.AddHandler(leaveHandler)
	session.AddHandler(joinHandler)

	// It's alive!
	fmt.Println("Hi, master " + Master.Username + ". I am " + Me.Username + ", and everything's all right!")

	// Everything is fine!
	return nil
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Myself?
	if m.Author.ID == Me.ID {
		return
	}

	// Get channel structure
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't get a channel structure.")
		fmt.Println("Author : " + m.Author.Username)
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
		return
	}

	// Get guild structure
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Couldn't get a guild structure.")
		fmt.Println("Channel : " + channel.Name)
		fmt.Println("Author : " + m.Author.Username)
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
		return
	}

	// Get guild member
	member, err := s.GuildMember(channel.GuildID, m.Author.ID)
	if err != nil {
		fmt.Println("Couldn't get a member structure.")
		fmt.Println("Guild : " + guild.Name)
		fmt.Println("Channel : " + channel.Name)
		fmt.Println("Author : " + m.Author.Username)
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
		return
	}

	// DM?
	if channel.Type == discordgo.ChannelTypeDM {

		// Popcorn?
		commands.Popcorn(s, m)

		if Master.ID == "" {

			// No BotMaster
			fmt.Println(m.Author.Username + " : " + m.Content)

		} else if m.Author.ID == Master.ID {

			// Talking to Master

		} else {

			// Get Master's User
			user, err := s.User(Master.ID)
			if err != nil {
				fmt.Println("Couldn't get Master's User!")
				fmt.Println(err.Error())
				return
			}

			// Create channel with Master
			masterChannel, err := s.UserChannelCreate(Master.ID)
			if err != nil {
				fmt.Println("Couldn't create a private channel with " + user.Username + ".")
				fmt.Println(err.Error())
				return
			}

			// Forward the message to BotMaster!
			s.ChannelTyping(masterChannel.ID)
			_, err = s.ChannelMessageSend(masterChannel.ID, "<@"+m.Author.ID+"> : "+m.Content)
			if err != nil {
				fmt.Println("Couldn't forward a message to " + user.Username + ".")
				fmt.Println("Author : " + m.Author.Username)
				fmt.Println("Message : " + m.Content)
				fmt.Println(err.Error())
				return
			}
		}
		return
	}

	/*
		// Update welcome channel
		if m.Type == discordgo.MessageTypeGuildMemberJoin {
			config.UpdateWelcomeChannel(s, m)
			return
		}
	*/

	// Bot?
	if m.Author.Bot {
		return
	}

	// Place in a guard
	commands.PlaceInAGuard(s, guild, channel, member, m.Message)

	// Mentionned someone?
	if len(m.Mentions) > 0 {
		for x := 0; x < len(m.Mentions); x++ {

			// Mentionned me?
			if m.Mentions[x].ID == Me.ID {

				// Split
				command := strings.Split(m.Content, " ")

				// Commands with 2 words
				if len(command) == 2 {
					if command[1] == "prune" {
						commands.Prune(s, m)
					}
				}

				// Commands with 3 words
				if len(command) == 3 {
					if command[1] == "get" {
						if command[2] == "points" {
							commands.GetPoints(s, m)
						}
					}
				}

				// Commands with 4 words
				if len(command) == 4 {
					if command[1] == "set" {
						if command[2] == "welcome" {
							if command[3] == "channel" && m.Author.ID == guild.OwnerID {
								commands.SetWelcomeChannel(s, m)
							}
						}
					}
					if command[1] == "get" {
						if command[2] == "welcome" {
							if command[3] == "channel" {
								commands.GetWelcomeChannel(s, m)
							}
						}
					}
				}
			}
		}
	}

	// Reactions
	commands.Popcorn(s, m)
	commands.Nani(s, m)
}

func reactHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	// Get the message structure
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println("Couldn't get the message structure of a MessageReactionAdd!")
		fmt.Println("ChannelID : " + m.ChannelID)
		fmt.Println(err.Error())
		return
	}

	// Get channel structure
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't get the channel structure of a MessageReactionAdd!")
		fmt.Println("ChannelID : " + m.ChannelID)
		fmt.Println("Author : " + message.Author.Username)
		fmt.Println("Message : " + message.Content)
		fmt.Println(err.Error())
		return
	}

	// Get the guild structure
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Couldn't get the guild structure of a MessageReactionAdd!")
		fmt.Println("Channel : " + channel.Name)
		fmt.Println("Author : " + message.Author.Username)
		fmt.Println("Message : " + message.Content)
		fmt.Println(err.Error())
		return
	}

	// Pin popular message
	pin(s, guild, channel, message)
}

func leaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// Get guild
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Couldn't get the guild of " + m.User.Username + "!")
		fmt.Println(err.Error())
		return
	}

	// Invite people who leave
	waitComeBack(s, guild, m.Member)
}

func joinHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	// Get guild
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Couldn't get the guild", m.User.Username, "just joined.")
		fmt.Println(err.Error())
		return
	}

	// Ask for guard
	askForGuard(s, guild, m.Member)
}
