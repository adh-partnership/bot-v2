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

// This is currently deadcode.... decided I don't really want to hit the VATUSA API
// for every unique user across all the ADH partner discords once daily or even every couple
// of days for this as that is a group of 100+ users all at once generally. This can be managed
// manually for now, at least until there is a VATUSA API endpoint I can hit once to get all of the staff
// in one go.

package division

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/vpaza/bot/pkg/cache"
	"github.com/vpaza/bot/pkg/network"
)

var Cache *cache.Cache

var log = logger.Logger.WithField("component", "start")
var VATUSAStaffRegexp = regexp.MustCompile(`^US\d+$`)

type VATUSAResponse struct {
	Data struct {
		roles []struct {
			Facility string `json:"facility"`
			Role     string `json:"role"`
		}
	} `json:"data"`
}

func IsDivisionStaff(cid string) bool {
	val, err := Cache.Get(cid)
	isStaff := false
	if err != nil {
		isStaff = lookupUser(cid)
		Cache.Set(cid, isStaff, time.Hour*72)
	} else {
		isStaff = val.(bool)
	}

	return isStaff
}

func lookupUser(cid string) bool {
	status, content, err := network.Call(
		"GET",
		fmt.Sprintf("https://api.vatusa.net/v2/user/%s", cid),
		"application/json",
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("Failed to get user info for %s: %s", cid, err)
		return false
	}

	if status != 200 {
		log.Errorf("Failed to get events for %s: %s", cid, content)
		return false
	}

	var vatusa VATUSAResponse
	err = json.Unmarshal(content, &vatusa)
	if err != nil {
		log.Errorf("Failed to unmarshal user info for %s: %s", cid, err)
		return false
	}

	for _, r := range vatusa.Data.roles {
		// Check if r.Facility is ZHQ and r.Role matches US\d+ regexp
		if r.Facility == "ZHQ" && VATUSAStaffRegexp.MatchString(r.Role) {
			return true
		}
	}

	return false
}
