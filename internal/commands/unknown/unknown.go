/*
 * Copyright Daniel Hawton
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package unknown

import (
	"slices"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/bwmarrin/discordgo"

	"github.com/vpaza/bot/internal/facility"
	"github.com/vpaza/bot/pkg/interactions"
	"github.com/vpaza/bot/pkg/utils"
)

var log = logger.Logger.WithField("component", "commands/unknown")

func Register() {
	interactions.AddCommand(&interactions.AppInteraction{
		Name:    "unknown",
		Handler: Handler,
		AppCommand: &discordgo.ApplicationCommand{
			Name:         "unknown",
			Description:  "Toggle display of unknown controllers on position",
			DMPermission: utils.PointerOf(false),
		},
	})
}

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	f, err := facility.FindFacility(
		&facility.Facility{
			DiscordID: i.GuildID,
		},
	)
	if err != nil {
		log.Errorf("Error finding facility: %s", err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Error finding facility",
			},
		})
		if err != nil {
			log.Errorf("Error responding to unknown command: %s", err)
		}
	}

	if f.UnknownControllers.RoleID == "" || !slices.Contains(i.Member.Roles, f.UnknownControllers.RoleID) {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You do not have permission to use this command",
			},
		})
		if err != nil {
			log.Errorf("Error responding to unknown command: %s", err)
		}
		return
	}

	if !f.UnknownControllers.Enabled {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown controller notifications are not enabled. Please contact your WM",
			},
		})
		if err != nil {
			log.Errorf("Error responding to unknown command: %s", err)
		}
		return
	}

	f.UnknownControllers.TempDisabled = !f.UnknownControllers.TempDisabled
	state := "off"
	if f.UnknownControllers.TempDisabled {
		state = "on"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Unknown controller notifications has been turned " + state,
		},
	})
	if err != nil {
		log.Errorf("Error responding to unknown command: %s", err)
	}
}
