package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/NatoBoram/Go-Miiko/wheel"
)

func pin(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.Message) {

	// DM?
	if c.Type == discordgo.ChannelTypeDM {
		return
	}

	// Get the reactions
	var singleReactionCount int
	for _, reaction := range m.Reactions {
		singleReactionCount = wheel.MaxInt(singleReactionCount, reaction.Count)
	}

	// Minimum reactions
	minReactions, err := getMinimumReactions(g, c)
	if err != nil {
		printDiscordError("Couldn't get the minimum reactions for a channel.", g, c, m, nil, err)
	}

	// Check the reactions
	if singleReactionCount >= minReactions {

		// Pin it!
		err = s.ChannelMessagePin(c.ID, m.ID)
		if err != nil {
			printDiscordError("Couldn't pin a popular message!", g, c, m, nil, err)

			// Check the amount of pins in that channel
			messages, err := s.ChannelMessagesPinned(c.ID)
			if err != nil {
				printDiscordError("Couldn't obtain the amount of pins in a channel.", g, c, m, nil, err)
				return
			}

			// Upgrade the minimum
			if len(messages) >= 50 {
				err = addMinimumReactions(c)
				if err != nil {
					printDiscordError("Couldn't add to the minimum reactions of a channel", g, c, m, nil, err)
					return
				}

				purgePin(s, g, c, m, messages)
			}

			return
		}

		// Check if already in the database
		_, err := selectPin(m)
		if err == sql.ErrNoRows {

			// Status
			err = setStatus(s, "épingler "+m.Author.Username)
			if err != nil {
				printDiscordError("Couldn't set the status to pinning someone.", g, c, m, nil, err)
			}

			// Throw it in the hall of fame
			savePin(s, g, m)
			if err == sql.ErrNoRows {
				// There's no hall of fame in this server
			} else if err != nil {
				printDiscordError("Couldn't throw a message in the hall of fame.", g, c, m, nil, err)
			}

			// Not previously pinned, time to insert it!
			_, err = insertPin(g, m)
			if err != nil {
				printDiscordError("Couldn't insert a pin.", g, c, m, nil, err)
			}
		} else if err != nil {
			printDiscordError("Couldn't select a pin.", g, c, m, nil, err)
		}
	}
}

func purgePin(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.Message, messages []*discordgo.Message) {
	err := setStatus(s, "augmenter la difficulté de "+c.Name)
	if err != nil {
		printDiscordError("Couldn't set the status to upping difficulty for a channel.", g, c, m, nil, err)
	}

	// Check minimum
	channelMin, err := selectMinimumReactions(c)
	if err != nil {
		printDiscordError("Couldn't add to the minimum reactions of a channel", g, c, m, nil, err)
	}

	// For each messages
	for _, message := range messages {

		// Check pins
		var singleReactionCount int
		for _, reaction := range message.Reactions {
			singleReactionCount = wheel.MaxInt(singleReactionCount, reaction.Count)
		}

		// Unpin
		if channelMin > singleReactionCount {
			err := s.ChannelMessageUnpin(c.ID, message.ID)
			if err != nil {
				printDiscordError("Couldn't unpin a previously popular message", g, c, message, nil, err)
				continue
			}
		}

		// Delete pin
		_, err := deletePin(message)
		if err != nil {
			printDiscordError("Couldn't remove a pin from the database.", g, c, message, nil, err)
		}
	}
}

func savePin(s *discordgo.Session, g *discordgo.Guild, m *discordgo.Message) (saved bool) {

	// Only text messages are transferred. Delete empty messages.
	if m.Type != discordgo.MessageTypeDefault || m.Content == "" {
		return
	}

	// Get the hall of fame
	halloffame, err := getFameChannel(s, g)
	if err != nil || m.ChannelID == halloffame.ID {
		printDiscordError("Couldn't get the hall of fame.", g, nil, m, nil, err)
		return
	}

	// Get Member
	member, err := s.GuildMember(g.ID, m.Author.ID)
	if err != nil {
		printDiscordError("Couldn't get a pinned member.", g, nil, m, nil, err)
		return
	}

	// Get colour
	colour, _ := getColour(s, g, member)

	// Get name
	var author string
	if member.Nick == "" {
		author = member.User.Username
	} else {
		author = member.Nick
	}

	// Create Embed
	embed := &discordgo.MessageEmbed{
		Color: colour,
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://canary.discordapp.com/channels/" + g.ID + "/" + m.ChannelID + "/" + m.ID + "/",
			Name:    author,
			IconURL: m.Author.AvatarURL(""),
		},
		URL:         "https://canary.discordapp.com/channels/" + g.ID + "/" + m.ChannelID + "/" + m.ID + "/",
		Title:       "Message",
		Description: m.Content,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Salon",
				Value:  "<#" + m.ChannelID + ">",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Auteur",
				Value:  "<@" + m.Author.ID + ">",
				Inline: true,
			},
		},
	}

	if len(m.Attachments) > 0 {

		// For all attachments
		for _, attachment := range m.Attachments {

			// Thumbnail
			if embed.Thumbnail == nil && attachment.Width >= 80 && attachment.Height >= 80 {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
					URL:    attachment.URL,
					Width:  attachment.Width,
					Height: attachment.Height,
				}
			}

			// Fields
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "Attachement",
				Value:  "[" + attachment.Filename + "](" + attachment.URL + ")",
				Inline: true,
			})

			// Image
			if embed.Image == nil && attachment.Width >= 520 {
				embed.Image = &discordgo.MessageEmbedImage{
					URL:    attachment.URL,
					Width:  attachment.Width,
					Height: attachment.Height,
				}
			}
		}
	}

	// Footer
	footer := &discordgo.MessageEmbedFooter{
		IconURL: discordgo.EndpointGuildIcon(g.ID, g.Icon),
	}
	for _, reaction := range m.Reactions {
		footer.Text += "<:" + reaction.Emoji.APIName() + ">"
	}

	// Send embed
	_, err = s.ChannelMessageSendEmbed(halloffame.ID, embed)
	if err != nil {
		printDiscordError("Couldn't send an embed.", g, nil, m, nil, err)
		return
	}

	// Save it in the database
	_, err = insertMessagesFamed(g, m)
	if err != nil {
		printDiscordError("Couldn't insert a famed message in the hall of fame!", g, nil, m, nil, err)
		return
	}

	return true
}
