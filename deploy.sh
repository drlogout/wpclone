#!/usr/bin/env bash

if [ -z "$WPCLONE_INSTALL_DIR" ]; then
    echo "WPCLONE_INSTALL_DIR is not set. Please set it to the directory where you want to install WP-Clone."
    exit 1
fi

make build-all
cp -r bin/* "$WPCLONE_INSTALL_DIR"
