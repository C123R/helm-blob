#!/bin/sh -e

if [ -n "${HELM_BLOB_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
echo "Downloading and installing helm-blob v${version} ..."

# Downloding binary from github
url=""
if [ "$(uname)" = "Darwin" ]; then
    url="https://github.com/PTC-Global/helm-blob/releases/download/v${version}/helm-blob_darwin_amd64_v${version}.tar.gz"
elif [ "$(uname)" = "Linux" ]; then
    if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
        url="https://github.com/PTC-Global/helm-blob/releases/download/v${version}/helm-blob_linux_arm64_v${version}.tar.gz"
    else
        url="https://github.com/PTC-Global/helm-blob/releases/download/v${version}/helm-blob_linux_amd64_v${version}.tar.gz"
    fi
else
    url="https://github.com/PTC-Global/helm-blob/releases/download/v${version}/helm-blob_windows_amd64_v${version}.tar.gz"
fi

# Installing binary

mkdir -p "bin"
mkdir -p "releases/v${version}"

# Download with curl if possible.
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/v${version}.tar.gz"
else
    wget -q "${url}" -O "releases/v${version}.tar.gz"
fi
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"
mv "releases/v${version}/helm-blob" "bin/helm-blob" ||
    mv "releases/v${version}/helm-blob.exe" "bin/helm-blob"

mv "releases/v${version}/scripts/proxy.sh" "bin/proxy.sh"
mv "releases/v${version}/plugin.yaml" "bin/plugin.yaml"
mv "releases/v${version}/README.md" "bin/README.md"
mv "releases/v${version}/LICENSE" "bin/LICENSE"

rm -rf "releases/v${version}.tar.gz"
