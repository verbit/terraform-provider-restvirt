#!/usr/bin/env bash
GOOS=$(awk '{print tolower($0)}' <<<"$(uname -s)")
GOARCH=$(uname -m)
case $GOARCH in
x86_64)
  GOARCH=amd64
  ;;
esac

URL=$(curl -s https://api.github.com/repos/verbit/terraform-provider-restvirt/releases/latest | jq -r ".assets[] | select(.name|endswith(\"${GOOS}_${GOARCH}.zip\")) | .browser_download_url")
TARGET=~/.terraform.d/plugins/github.com/verbit/restvirt/
mkdir -p $TARGET
cd $TARGET || { echo "couldn't cd into $TARGET"; exit 1; }
echo "Downloading provider from $URL"
curl -L -O "$URL"
echo "done"
