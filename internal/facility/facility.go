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
	"errors"
	"fmt"
	"os"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/vpaza/bot/pkg/cache"
	"sigs.k8s.io/yaml"
)

type Facility struct {
	Facility                   string `json:"facility"`
	BotName                    string `json:"bot_name"`
	DiscordID                  string `json:"discord_id"`
	Description                string `json:"description"`
	NameFormat                 string `json:"name_format"`
	StaffFormat                string `json:"staff_format"`
	StaffTitleSeparator        string `json:"staff_title_separator"`
	RosterAPI                  string `json:"roster_api"`
	API                        string `json:"api"`
	BaseURL                    string `json:"base_url"`
	ShowOnline                 string `json:"show_online"`
	OnlineChannel              string `json:"online_channel"`
	NoControllersOnlineMessage string `json:"no_controllers_online_message"`
	Positions                  []struct {
		Name      string `json:"name"`
		Callsigns struct {
			Prefix []string `json:"prefix"`
			Suffix []string `json:"suffix"`
		} `json:"callsigns"`
	} `json:"positions"`
	Roles []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		If   []struct {
			Condition string `json:"condition"`
			Value     string `json:"value"`
		} `json:"if"`
	} `json:"roles"`
	UnknownControllers *struct {
		Enabled      bool `json:"enabled"`
		TempDisabled bool
		Channel      string `json:"channel"`
		RoleID       string `json:"role_id"`
	}
	cache *cache.Cache
}

var FacCfg map[string]*Facility

var (
	ErrorFacilityExists   = errors.New("facility already exists")
	ErrorFacilityNotFound = errors.New("facility not found")
	log                   = logger.Logger.WithField("component", "facility")
)

func init() {
	FacCfg = make(map[string]*Facility)
}

func FindFacility(f *Facility) (*Facility, error) {
	for _, cfg := range FacCfg {
		if cfg.DiscordID == f.DiscordID {
			return cfg, nil
		}
	}

	return nil, ErrorFacilityNotFound
}

func ParseFacilityConfig(file string, c *cache.Cache) (*Facility, error) {
	config, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &Facility{}
	err = yaml.Unmarshal(config, cfg)
	if err != nil {
		return nil, err
	}

	if cfg.StaffFormat == "all" && cfg.StaffTitleSeparator == "" {
		cfg.StaffTitleSeparator = "/"
		logger.Logger.WithField("component", "config").Warnf("Staff format is set to 'all' but no staff title separator is set. Defaulting to '/'")
	}

	if cfg.UnknownControllers == nil {
		cfg.UnknownControllers = &struct {
			Enabled      bool `json:"enabled"`
			TempDisabled bool
			Channel      string `json:"channel"`
			RoleID       string `json:"role_id"`
		}{
			Enabled:      false,
			TempDisabled: false,
			Channel:      "",
			RoleID:       "",
		}
	}

	if _, ok := FacCfg[cfg.Facility]; ok {
		return nil, ErrorFacilityExists
	}

	cfg.cache = c
	FacCfg[cfg.Facility] = cfg

	return cfg, nil
}

func (f *Facility) GetNoControllersOnlineMessage() string {
	if f.NoControllersOnlineMessage == "" {
		return fmt.Sprintf("There are currently no %s controllers online.", f.Facility)
	}

	return f.NoControllersOnlineMessage
}
