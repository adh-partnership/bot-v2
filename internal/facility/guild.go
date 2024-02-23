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

package facility

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (f *Facility) GetOwnerID(s *discordgo.Session) string {
	b, err := f.cache.Get(
		fmt.Sprintf("/%s/owner", f.Facility),
	)
	if err == nil {
		return b.(string)
	}

	g, err := s.Guild(f.DiscordID)
	if err != nil {
		log.Errorf("Failed to get guild %s: %s", f.DiscordID, err)
		return ""
	}

	return g.OwnerID
}
