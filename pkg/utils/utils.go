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

package utils

import (
	"encoding/json"
	"os"
	"time"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Now() *time.Time {
	t := time.Now()
	return &t
}

func EnvOrDefault(env string, def string) string {
	if v, ok := os.LookupEnv(env); ok {
		return v
	}
	return def
}

func PointerOf[T any](v T) *T {
	return &v
}

func Trim(s string, length int) string {
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}

func MapJSON(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}
