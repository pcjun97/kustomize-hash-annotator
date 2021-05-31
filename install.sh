#!/bin/bash
set -e

# Require $XDG_CONFIG_HOME to be set
if [[ -z "$XDG_CONFIG_HOME" ]]; then
  echo "You must define XDG_CONFIG_HOME to use a kustomize plugin"
  echo "Add 'export XDG_CONFIG_HOME=\$HOME/.config' to your .bashrc or .zshrc"
  exit 1
fi

PLUGIN_PATH="$XDG_CONFIG_HOME/kustomize/plugin/pcjun97/v1/hashannotator"
PLUGIN_NAME="HashAnnotator"

mkdir -p $PLUGIN_PATH

echo "Copying exec plugin to the kustomize plugin path..."
echo "cp $PLUGIN_NAME $PLUGIN_PATH/$PLUGIN_KIND"
cp $PLUGIN_NAME "$PLUGIN_PATH/$PLUGIN_KIND"
