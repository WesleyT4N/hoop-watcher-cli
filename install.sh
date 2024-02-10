#!/usr/bin/env bash

INSTALL_PATH=$HOME/bin/hoop-watcher-cli
cp ./cmd/hoop-watcher-cli/hoop-watcher-cli $INSTALL_PATH
cp ./cmd/hoop-watcher-cli/nba_teams.json $HOME/bin/nba_teams.json
echo "installed hoop-watcher-cli to $INSTALL_PATH"
