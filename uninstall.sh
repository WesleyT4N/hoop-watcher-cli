#!/usr/bin/env bash

# get home directory
INSTALL_PATH=$HOME/bin/hoop-watcher-cli

# check if hoop-watcher is installed
if [ ! -f $INSTALL_PATH ]; then
  echo "hoop-watcher is not installed"
  exit 1
fi
rm $INSTALL_PATH
rm $HOME/bin/nba_teams.json
echo "uninstalled hoop-watcher-cli from $INSTALL_PATH"
