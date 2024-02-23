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

package events

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"

	"github.com/vpaza/bot/internal/facility"
)

func EventPostSelectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Infof("Received event-select component from %s in %s #%s", i.Member.User.ID, i.GuildID, i.ChannelID)
	f, err := facility.FindFacility(
		&facility.Facility{
			DiscordID: i.GuildID,
		},
	)
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error finding facility",
			},
		})
		if err != nil {
			log.Warnf("Error responding to channel: %s", err)
		}
		return
	}

	val := &EventSelectData{}
	err = json.Unmarshal([]byte(i.MessageComponentData().Values[0]), &val)
	if err != nil {
		log.Errorf("Failed to parse event data: %s", err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error finding event",
			},
		})
		if err != nil {
			log.Warnf("Error responding to channel: %s", err)
		}
		return
	}

	event := f.GetEvent(val.EventID)
	if event == nil {
		log.Errorf("Failed to find event with id %d", val.EventID)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error finding event",
			},
		})
		if err != nil {
			log.Warnf("Error responding to channel: %s", err)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Creating post for event " + event.Title,
		},
	})
	if err != nil {
		log.Warnf("Error responding to channel: %s", err)
	}

	switch val.Type {
	case "announce":
		postAnnouncement(s, i, event, f)
	case "position":
		postPositions(s, i, event, f)
	}
}
