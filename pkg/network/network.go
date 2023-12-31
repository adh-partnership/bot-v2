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

package network

import (
	"net/url"

	adhnetwork "github.com/adh-partnership/api/pkg/network"
)

func init() {
	adhnetwork.UserAgent = "ADH/bot"
}

func Call(method, requrl string, contenttype string, formdata map[string]string, headers map[string]string) (int, []byte, error) {
	u, err := url.Parse(requrl)
	if err != nil {
		return 0, nil, err
	}

	data := url.Values{}

	for k, v := range formdata {
		data.Set(k, v)
	}

	return adhnetwork.HandleWithHeaders(
		method,
		u.String(),
		contenttype,
		data.Encode(),
		headers,
	)
}
