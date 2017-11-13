package bot

import (
	"fmt"
	"math/rand"
	"time"

	"../config"
	"github.com/bwmarrin/discordgo"
)

func waitComeBack(s *discordgo.Session, m *discordgo.GuildMemberRemove) {

	// Get guild
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Couldn't get " + m.User.Username + "'s guild ID.")
		fmt.Println(err.Error())
		return
	}

	// Get channel
	channel := config.GetWelcomeChannelByGuildID(guild.ID)
	if channel == "" {
		return
	}

	// Create an invite structure
	var invStruct discordgo.Invite
	invStruct.Temporary = true

	// Invite to WelcomeChannel
	var invite *discordgo.Invite
	invite, err = s.ChannelInviteCreate(channel, invStruct)
	if err != nil {
		fmt.Println("Couldn't create an invite in " + guild.Name + ".")
		fmt.Println(err.Error())
		return
	}

	// Bot?
	if !m.User.Bot {

		// Open channel
		privateChannel, err := s.UserChannelCreate(m.User.ID)
		if err != nil {
			fmt.Println("Couldn't create a private channel with " + m.User.Username + ".")
			fmt.Println(err.Error())
			return
		}

		// Typing!
		err = s.ChannelTyping(privateChannel.ID)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Send message
		_, err = s.ChannelMessageSend(privateChannel.ID, getPrivateByeMessage(invite.Code))
		if err != nil {
			fmt.Println("Couldn't say bye to " + m.User.Username + "!")
			fmt.Println(err.Error())
		}

	} else {

		// Typing!
		err = s.ChannelTyping(channel)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Send message
		_, err = s.ChannelMessageSend(channel, getByeBotMessage(m.User.ID))
		if err != nil {
			fmt.Println("Couldn't say bye to " + m.User.Username + "!")
			fmt.Println(err.Error())
		}
	}
}

func getPrivateByeMessage(inviteCode string) string {

	// Bye Messages
	var byeList []string

	// Messages
	byeList = append(byeList, "Oh, je suis triste de te voir partir! Si tu veux nous rejoindre à nouveau, j'ai créé une invitation pour toi : https://discord.gg/"+inviteCode)
	byeList = append(byeList, "Au revoir! Voici une invitation si tu changes d'idée : https://discord.gg/"+inviteCode)
	byeList = append(byeList, "Tu vas me manquer. Si tu veux me revoir, j'ai créé une invitation pour toi : https://discord.gg/"+inviteCode)

	// Random
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rand := rand.New(source)

	// Return
	return byeList[rand.Intn(len(byeList))]
}

func getByeBotMessage(userID string) string {

	// Bye Messages
	var byeBotList []string

	// Messages
	byeBotList = append(byeBotList, "Bon débarras, <@"+userID+">.")
	byeBotList = append(byeBotList, "Bien! Personne ne va s'ennuyer de <@"+userID+">.")
	byeBotList = append(byeBotList, "De toute façon, <@"+userID+"> n'avait aucun lien avec Eldarya.")

	// Random
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rand := rand.New(source)

	// Return
	return byeBotList[rand.Intn(len(byeBotList))]
}
