package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// SetWelcomeChannel sets the welcome channel
func SetWelcomeChannel(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Get channel structure
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Couldn't get the channel structure of a welcome message.")
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

	s.ChannelTyping(channel.ID)
	//config.UpdateWelcomeChannel(s, m)
	_, err = s.ChannelMessageSend(channel.ID, "D'accord! <#"+channel.ID+"> est maintenant le salon de bienvenue.")
	if err != nil {
		fmt.Println("Couldn't send a message.")
		fmt.Println("Guild : " + guild.Name)
		fmt.Println("Channel : " + channel.Name)
		fmt.Println("Author : " + m.Author.Username)
		fmt.Println("Message : " + m.Content)
		fmt.Println(err.Error())
	}
}

// GetWelcomeChannel gets the welcome channel
func GetWelcomeChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*
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

		// Does it exists?


			welcomeChannelID, exists := config.Database.WelcomeChannels[guild.ID]

			if exists {

				// Get channel structure
				welcomeChannel, err := s.State.Channel(welcomeChannelID)
				if err != nil {
					fmt.Println(guild.Name + "'s welcome channel doesn't exist anymore.")
					fmt.Println(err.Error())

					// Set this channel as the WelcomeChannel
					config.UpdateWelcomeChannel(s, m)
					welcomeChannel = channel
				}

				// Send the welcome channel
				s.ChannelTyping(channel.ID)
				_, err = s.ChannelMessageSend(channel.ID, "Le salon de bienvenue est <#"+welcomeChannel.ID+">.")
				if err != nil {
					fmt.Println("Couldn't send a message.")
					fmt.Println("Guild : " + guild.Name)
					fmt.Println("Channel : " + channel.Name)
					fmt.Println("Author : " + m.Author.Username)
					fmt.Println("Message : " + m.Content)
					fmt.Println(err.Error())
					return
				}

			} else {

				// Let it be this one
				config.UpdateWelcomeChannel(s, m)

				// Send the welcome channel
				s.ChannelTyping(channel.ID)
				_, err = s.ChannelMessageSend(channel.ID, "Le salon de bienvenue est <#"+channel.ID+">.")
				if err != nil {
					fmt.Println("Couldn't send a message.")
					fmt.Println("Guild : " + guild.Name)
					fmt.Println("Channel : " + channel.Name)
					fmt.Println("Author : " + m.Author.Username)
					fmt.Println("Message : " + m.Content)
					fmt.Println(err.Error())
					return
				}
			}
	*/
}
