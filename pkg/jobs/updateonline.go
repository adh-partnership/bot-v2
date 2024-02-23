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

package jobs

import (
	"fmt"
	"strings"
	"time"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/adh-partnership/api/pkg/network/vatsim"
	"github.com/bwmarrin/discordgo"

	"github.com/vpaza/bot/internal/bot"
	"github.com/vpaza/bot/internal/facility"
	"github.com/vpaza/bot/pkg/utils"
)

var uolog = logger.Logger.WithField("component", "jobs/updateonline")

var unknownControllers map[string]map[string]bool

func init() {
	unknownControllers = make(map[string]map[string]bool)
}

func UpdateOnline() {
	if discord == nil || !discord.DataReady {
		log.Infof("Discord data not ready, skipping UpdateGuilds job")
		return
	}

	data, err := vatsim.GetData()
	if err != nil {
		uolog.Errorf("Error getting data from VATSIM: %s", err)
		return
	}
	for _, f := range facility.FacCfg {
		if f.OnlineChannel == "" {
			continue
		}
		go updateOnline(f, data)
	}
}

func updateOnline(f *facility.Facility, data *vatsim.VATSIMData) {
	uolog.Debugf("Starting UpdateOnline job for %s", f.Facility)
	start := time.Now()

	online := make(map[string][]string)
	newUnknown := make(map[string]bool)

	for _, p := range f.Positions {
		online[p.Name] = []string{}
	}

	for _, c := range data.Controllers {
		if f.UnknownControllers.Enabled && !f.UnknownControllers.TempDisabled && !isRostered(f, fmt.Sprint(c.CID)) {
			if _, ok := unknownControllers[f.Facility][genIndex(c)]; !ok {
				_, err := bot.GetSession().ChannelMessageSend(f.UnknownControllers.Channel, fmt.Sprintf("Unrostered controller on %s, %s (%d)", c.Callsign, c.Name, c.CID))
				if err != nil {
					uolog.Errorf("Error sending message: %s", err)
				}
				newUnknown[genIndex(c)] = true
			}
		}

		group := findGroup(f, c.Callsign)
		if group != "" {
			online[group] = append(
				online[group],
				fmt.Sprintf(
					"%s - %s - %s - <t:%d:R>",
					c.Callsign,
					getOI(f, fmt.Sprint(c.CID)),
					c.Frequency,
					c.LogonTime.Unix(),
				),
			)
		}
	}

	// a little expensive ops wise, but ensures we clear out dropped connections
	unknownControllers[f.Facility] = newUnknown

	message := &discordgo.MessageEmbed{
		Author:    &discordgo.MessageEmbedAuthor{},
		Color:     0x00ff00,
		Fields:    []*discordgo.MessageEmbedField{},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "Online Controllers",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "END OF LINE.",
		},
	}

	for _, p := range f.Positions {
		if len(online[p.Name]) == 0 {
			online[p.Name] = []string{f.GetNoControllersOnlineMessage()}
		}
		message.Fields = append(
			message.Fields,
			&discordgo.MessageEmbedField{
				Name:   p.Name,
				Value:  strings.Join(online[p.Name], "\n"),
				Inline: false,
			},
		)
	}

	post, err := bot.FindBotMessageIn(f.OnlineChannel)
	if err != nil {
		uolog.Errorf("Error finding bot message in %s: %s", f.OnlineChannel, err)
		return
	}

	if post == nil {
		_, err := bot.GetSession().ChannelMessageSendEmbed(f.OnlineChannel, message)
		if err != nil {
			log.Errorf("Error sending message: %s", err)
		}
	} else {
		_, err := bot.GetSession().ChannelMessageEditEmbed(f.OnlineChannel, post.ID, message)
		if err != nil {
			log.Errorf("Error sending message: %s", err)
		}
	}

	uolog.Debugf("Finished UpdateOnline job for %s (took: %s)", f.Facility, time.Since(start))
}

func findGroup(f *facility.Facility, callsign string) string {
	parts := strings.Split(callsign, "_")
	for _, p := range f.Positions {
		if utils.Contains(p.Callsigns.Prefix, parts[0]) && utils.Contains(p.Callsigns.Suffix, parts[len(parts)-1]) {
			return p.Name
		}
	}
	return ""
}

func genIndex(c *vatsim.VATSIMController) string {
	return fmt.Sprintf("%d/%s", c.CID, c.Callsign)
}

func isRostered(f *facility.Facility, cid string) bool {
	_, err := f.FindUserByCID(cid)
	return err == nil
}

func getOI(f *facility.Facility, cid string) string {
	u, err := f.FindUserByCID(cid)
	if err != nil {
		return cid
	}

	return u.OperatingInitials
}
