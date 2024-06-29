#!/bin/bash

#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# GitHub repository owner and name
REPO_OWNER="ashish10alex"
REPO_NAME="formatdataform"
BINARY="formatdataform"

# Detect the operating system and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

echo "Detected OS: $OS"
echo "Detected architecture: $ARCH"


# if os is darwin then change it to Darwin
# if os is linux then change it to Linux
case $OS in
    darwin)
        OS="Darwin"
        ;;
    linux)
        OS="Linux"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Translate architecture names to match GitHub release naming
case $ARCH in
    x86_64)
        ARCH="x86_64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "OS name for relase url template: $OS"
echo "Architecture name for relase url template: $ARCH"

# Get the latest release download URL for the appropriate tar.gz file
RELEASE_URL=$(curl -s https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest | \
    grep "browser_download_url.*${OS}.*${ARCH}.*tar.gz" | cut -d '"' -f 4)


echo "Url to get latest release: $RELEASE_URL"

# Check if the URL is empty
if [ -z "$RELEASE_URL" ]; then
    echo "Could not find a release for OS: $OS and ARCH: $ARCH"
    exit 1
fi

echo "Creating temporary directory for download and extraction"
TMP_DIR=$(mktemp -d)

cd $TMP_DIR

echo "Downloading the latest release into the directory"
curl -L -o release.tar.gz $RELEASE_URL

echo "Extracting the .tar.gz file"
tar -xzvf release.tar.gz


# Check if the binary is found
if [ -z "$BINARY" ]; then
    echo "Could not find the binary file."
    exit 1
fi

echo "Moving the $BINARY binary to /usr/local/bin (requires sudo)"
sudo mv $BINARY /usr/local/bin/

echo "Making the $BINARY executable"
sudo chmod +x /usr/local/bin/$BINARY

echo "Going back to: "
cd -
rm -rf $TMP_DIR

echo "Installation completed"

echo " `$BINARY --version`"
echo ""
echo "Try  $BINARY --help  to see the available options"

