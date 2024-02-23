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

package commands

import (
	"github.com/adh-partnership/api/pkg/logger"
	"github.com/bwmarrin/discordgo"

	"github.com/vpaza/bot/internal/commands/events"
	"github.com/vpaza/bot/internal/commands/ping"
	"github.com/vpaza/bot/internal/commands/unknown"
	"github.com/vpaza/bot/pkg/interactions"
)

var log = logger.Logger.WithField("component", "commands")

func SetupCommands() {
	ping.Register()
	events.Register()
	unknown.Register()
}

func RegisterCommands(s *discordgo.Session, guildid string) error {
	// Clean all commands associated with us
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, guildid)
	if err != nil {
		return err
	}
	for _, registeredCommand := range registeredCommands {
		log.Warnf("Unregistering command %s, we didn't register", registeredCommand.Name)
		err := s.ApplicationCommandDelete(s.State.User.ID, guildid, registeredCommand.ID)
		if err != nil {
			return err
		}
	}

	for _, command := range interactions.GetCommands() {
		log.Infof("Registering command %s in %s", command.Name, guildid)
		appCommand, err := s.ApplicationCommandCreate(s.State.User.ID, guildid, command.AppCommand)
		if err != nil {
			return err
		}
		command.ID = appCommand.ID
	}

	return nil
}

func Unregister(s *discordgo.Session, guildid string) error {
	for _, command := range interactions.GetCommands() {
		if command.ID != "" {
			log.Infof("Unregistering command %s", command.Name)
			err := s.ApplicationCommandDelete(s.State.User.ID, guildid, command.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
