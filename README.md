# Master Control Program v2

A rewrite of the old [Master Control Program](https://github.com/adh-partnership/bot)
to be more configurable (aka, not hardcoded), modular, and be able to be used in multiple
guilds instead of running multiple bots, one for each guild.

While this may not happen, it may prove useful for some implementations.

## Installation

TODO

## Configuration

### config.yaml

This configuration defines the basic global configuration settings for the bot.

Example config.yaml:

```yaml
database:
  host: 127.0.0.1
  port: "3306"
  user: root
  password: bot
  database: bot
  auto_migrate: true
discord:
  appid: "1234567"
  token: bot.token.discord
```

### Facility Configuration

These are YAML files generally in `facilities/` and are used to define facility guilds.

An example:

```yaml
# Used for internal identification purposes
facility: ZAN
# The nickname the bot should set for itself in this guild
bot_name: Go Bot
description: Anchorage ARTCC
# The Guild ID for the facility's discord
discord_id: 123456789012345678
# What name format do we set for users in the guild?
# - first_cid (John 888888, John 888888 | ATM)
# - first_last (John Doe, John Doe | ATM)
#   *NOTE* If this is set and staff_format is all, or the name is long, that user's last name will be truncated to last initial
#   to fit within Discord's 32 character limit
# - first_last_initial (John D., John D. | ATM)
# - first_last_initial_oi (John D. - JD, John D. - JD | ATM)
# default: first_last
name_format: first_last
# At the end of the name we'll append the staff role(s)
# - highest: The highest staff role from ATM, DATM, TA, EC, FE, WM, INS, MTR, AEC, AFE, AWM
# - all: All staff roles from ATM, DATM, TA, EC, FE, WM, INS, MTR, AEC, AFE, AWM separated by the defined seperator
# - none: No staff roles
# default: highest
staff_format: highest
# Separator to use between the staff role(s) when staff_format is all
# Defaults to "/" when not defined
staff_title_separator: "/"
# Position table, used for the who's online embedded message
# that will be posted to the positions_channel
positions:
  - name: Oceanic
    callsigns:
      - prefix: ["ZAN"]
        suffix: ["FSS"]
  - name: Enroute
    callsigns:
      prefix: ["ANC"]
      suffix: ["CTR"]
  - name: A11
    callsigns:
      prefix: ["ANC"]
      suffix: ["APP"]
  - name: FAI TRACON
    callsigns:
      prefix: ["FAI"]
      suffix: ["APP"]
  - name: CAB
    callsigns:
      prefix:
        - ADQ
        - AKN
        - ANC
        - BIG
        - BET
        - EIL
        - EDF
        - ENA
        - FAI
        - FBK
        - FRN
        - JNU
        - LHD
        - MRI
      suffix:
        - TWR
        - GND
        - DEL
# The ID number of the channel to post the who's online message to
# This channel should ideally be setup so only the bot posts to it
online_channel: 1011814580848177195
# Message to use when no controllers are online
# Defaults to "There are currently no (Facility ID) controllers online."
no_controllers_online_message: |
  There are currently no ZAN controllers online.
# Role assignments
# Users who have a configured role but do not meet the conditions will have the role
# removed.
# Conditions are defined in the `if` block and are evaluated as a logical OR.
# Conditions:
# - has_role: Checks if the user has the specified role in the roster
# - controller_type: Check if the user is a `home` controller, `visitor` controller, or `none`
# - rating: Checks user rating (S1, S2, etc)
# - unknown: Checks if the account is linked on the roster
roles:
  - id: 123456789012345678 # Role ID
    name: ZAN Senior Staff # Role Name (primarily for logging)
    if: # Conditions evaluated as a logical OR
      - condition: has_role
        value: ATM
      - condition: has_role
        value: DATM
      - condition: has_role
        value: TA
  - id: 123456789012345678
    name: ZAN Staff
    if:
      - condition: has_role
        value: EC
      - condition: has_role
        value: FE
      - condition: has_role
        value: WM
  - id: 123456789012345678
    name: ZAN Assistance Staff
    if:
      - condition: has_role
        value: facilities
      - condition: has_role
        value: events
      - condition: has_role
        value: web
  - id: 123456789012345678
    name: Training Team
    if:
      - condition: has_role
        value: INS
      - condition: has_role
        value: MTR
  - id: 123456789012345678
    name: FE Team
    if:
      - condition: has_role
        value: facilities
  - id: 123456789012345678
    name: Events Team
    if:
      - condition: has_role
        value: events
  - id: 123456789012345678
    name: Instructors
    if:
      - condition: has_role
        value: INS
  - id: 123456789012345678
    name: Mentors
    if:
      - condition: has_role
        value: MTR
  - id: 123456789012345678
    name: ZAN Members
    if:
      - condition: controller_type
        value: home
  - id: 123456789012345678
    name: Visitors
    if:
      - condition: controller_type
        value: visitor
  - id: 123456789012345678
    name: Pilot
    if:
      - condition: controller_type
        value: none
  - id: 123456789012345678
    name: Unverified
    if:
      - condition: unknown
        value: true
# Base URL of the API
# Note that the Bot will use this API to fetch event information for the /event
# slash command as well as the roster. The API, if not ADH Partnership's API
# must support the `/v1/events` `/v1/event/:id` and `/v1/user/all` endpoints
# with the same response format as ADH Partnership's API
#
# ADH Partnership's Swagger Docs can be accessed at https://api.zanartcc.org
api: https://api.zanartcc.org
# Base URL of Front End
base_url: https://www.zanartcc.org
# If you want to handle unknown controller notifications from the bot, define this:
unknown_controllers:
  # Set enabled to true to enable the feature, when this is true will also enable the
  # /unknown slash command to toggle this feature during runtime
  enabled: true
  # The channel ID to send the notifications to, ideally should be a senior staff only channel
  channel: "123456789012345678"
  # The role ID that can use the slash command
  role_id: "123456789012345678"
```

## Running

To start the bot, it assumes that in your PWD you have a `config.yaml` file and there will be a `facilities` directory with yaml files. If this is not the case
you can override with arguments.

```shell
./bot start --config /path/to/config.yaml --facility-configs-path /path/to/facilities-directory
```

By default log level is set to `info`, if you need lower, you can add `--log-level <level>` prior to the subcommand.

Level is one of:

- error
- warn
- info (default)
- debug
- trace

Logs will always print for at and above the defined level.

```shell
./bot --log-level debug start --config /path/to/config.yaml --facility-configs-path /path/to/facilities-directory
```

## Kubernetes Manifests

Not provided. A way to build it is the define a `ConfigMap` to mount for the config yaml.

For more information, consult [the Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/)

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details
