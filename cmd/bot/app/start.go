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

package app

import (
	"context"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/adh-partnership/api/pkg/logger"
	"github.com/urfave/cli/v2"

	"github.com/vpaza/bot/internal/bot"
	"github.com/vpaza/bot/internal/commands"
	"github.com/vpaza/bot/internal/facility"
	"github.com/vpaza/bot/pkg/cache"
	"github.com/vpaza/bot/pkg/config"
	"github.com/vpaza/bot/pkg/jobs"
)

var log = logger.Logger.WithField("component", "start")

func newStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start discord bot",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   "config.yaml",
				Usage:   "Load configuration from `FILE`",
				Aliases: []string{"c"},
				EnvVars: []string{"CONFIG"},
			},
			&cli.StringFlag{
				Name:    "facility-configs-path",
				Usage:   "Path of facility configs",
				EnvVars: []string{"FACILITY_CONFIGS_PATH"},
				Value:   "facilities",
			},
		},
		Action: func(c *cli.Context) error {
			log.Infof("Starting bot...")
			log.Infof("config=%s", c.String("config"))

			log.Infof("Creating in-memory cache...")
			cacheService := cache.NewCache(context.Background(), time.Minute*3, time.Minute*3)

			log.Infof("Loading configuration...")
			_, err := config.ParseConfig(c.String("config"))
			if err != nil {
				return err
			}
			// Walk the facility configs path and load them
			err = filepath.WalkDir(c.String("facility-configs-path"), func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() && path != c.String("facility-configs-path") {
					log.Infof("Skipping %s as it's a directory", path)
					return filepath.SkipDir
				}

				if strings.HasSuffix(path, "yaml") || strings.HasSuffix(path, "yml") {
					log.Infof("Loading facility config %s", path)
					_, err = facility.ParseFacilityConfig(path, cacheService)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}

			log.Debugf("FacCfg=%+v", facility.FacCfg)

			log.Infof("Building jobs...")
			jobs.BuildJobs()

			log.Infof("Starting bot...")
			sess, err := bot.Start()
			if err != nil {
				return err
			}

			log.Infof("Starting jobs async...")
			jobs.Start(sess)

			log.Infof("Bot should be running...")

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
			<-stop

			log.Infof("Gracefully shutting down...")
			sess.Close()
			for _, facility := range facility.FacCfg {
				err := commands.Unregister(sess, facility.DiscordID)
				if err != nil {
					log.Warnf("Failed to unregister commands: %s", err)
				}
			}
			log.Infof("Done")

			return nil
		},
	}
}
