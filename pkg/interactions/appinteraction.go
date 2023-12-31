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

package interactions

import "github.com/bwmarrin/discordgo"

type AppInteraction struct {
	Name       string
	ID         string
	Handler    func(*discordgo.Session, *discordgo.InteractionCreate)
	AppCommand *discordgo.ApplicationCommand
}

var (
	commands   map[string]*AppInteraction
	components map[string]*AppInteraction
)

func init() {
	commands = make(map[string]*AppInteraction)
	components = make(map[string]*AppInteraction)
}

func AddCommand(a *AppInteraction) {
	commands[a.Name] = a
}

func AddCommands(a []*AppInteraction) {
	for _, c := range a {
		AddCommand(c)
	}
}

func AddComponent(a *AppInteraction) {
	components[a.Name] = a
}

func AddComponents(a []*AppInteraction) {
	for _, c := range a {
		AddComponent(c)
	}
}

func RemoveCommand(name string) {
	delete(commands, name)
}

func RemoveComponent(name string) {
	delete(components, name)
}

func FindCommand(name string) (*AppInteraction, bool) {
	if command, ok := commands[name]; ok {
		return command, true
	}
	return nil, false
}

func FindComponent(name string) (*AppInteraction, bool) {
	if component, ok := components[name]; ok {
		return component, true
	}
	return nil, false
}

func GetCommands() []*AppInteraction {
	var cmds []*AppInteraction
	for _, c := range commands {
		cmds = append(cmds, c)
	}
	return cmds
}

func GetComponents() []*AppInteraction {
	var comps []*AppInteraction
	for _, c := range components {
		comps = append(comps, c)
	}
	return comps
}
