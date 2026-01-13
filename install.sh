#!/bin/bash

set -e

if [ "$1" = "-t" ]; then
    # Copy templates only
    ENV_FILE=~/.personal-blog-generator/.env
    if [ ! -f "$ENV_FILE" ]; then
        echo "Configuration file not found at $ENV_FILE. Please run the setup script first."
        exit 1
    fi
    source "$ENV_FILE"
    if [ ! -d "$TEMPLATE_PATH" ]; then
        mkdir -p "$TEMPLATE_PATH"
    fi
    echo "Copying templates..."
    cp -r templates/* "$TEMPLATE_PATH"/
    echo "Templates copied successfully!"
    exit 0
elif [ "$1" = "-a" ]; then
    # Copy admin-files only
    ENV_FILE=~/.personal-blog-generator/.env
    if [ ! -f "$ENV_FILE" ]; then
        echo "Configuration file not found at $ENV_FILE. Please run the setup script first."
        exit 1
    fi
    source "$ENV_FILE"
    if [ ! -d "$ADMIN_FILES_PATH" ]; then
        mkdir -p "$ADMIN_FILES_PATH"
    fi
    echo "Copying admin files..."
    cp -r admin-files/* "$ADMIN_FILES_PATH"/
    echo "Admin files copied successfully!"
    exit 0
fi

echo "Personal Blog Generator Local Setup Script"
echo "This script will set up the blog generator on your local machine."
echo

# Prompt for config
echo "Configuring environment variables..."
echo "Press Enter to accept defaults."
echo

read -p "Enter APP_PORT (default 8080): " APP_PORT
APP_PORT=${APP_PORT:-8080}

read -p "Enter DB_PATH (default ~/.personal-blog-generator/blog.db): " DB_PATH
DB_PATH=${DB_PATH:-~/.personal-blog-generator/blog.db}

read -p "Enter ADMIN_FILES_PATH (default ~/.personal-blog-generator/admin-files): " ADMIN_FILES_PATH
ADMIN_FILES_PATH=${ADMIN_FILES_PATH:-~/.personal-blog-generator/admin-files}

read -p "Enter TEMPLATE_PATH (default ~/.personal-blog-generator/templates): " TEMPLATE_PATH
TEMPLATE_PATH=${TEMPLATE_PATH:-~/.personal-blog-generator/templates}

read -p "Enter OUTPUT_PATH (default ~/html-outputs): " OUTPUT_PATH
OUTPUT_PATH=${OUTPUT_PATH:-~/html-outputs}

# Expand ~ in paths
DB_PATH="${DB_PATH/#\~/$HOME}"
TEMPLATE_PATH="${TEMPLATE_PATH/#\~/$HOME}"
OUTPUT_PATH="${OUTPUT_PATH/#\~/$HOME}"
ADMIN_FILES_PATH="${ADMIN_FILES_PATH/#\~/$HOME}"

# Create directories
echo "Creating necessary directories..."
mkdir -p ~/bin
mkdir -p ~/.personal-blog-generator
mkdir -p "$TEMPLATE_PATH"
mkdir -p "$OUTPUT_PATH"
mkdir -p "$ADMIN_FILES_PATH"
mkdir -p "$(dirname "$DB_PATH")"
echo "Directories created."
echo

# Create .env
echo "Creating configuration file..."
cat > ~/.personal-blog-generator/.env << EOF2
# Database configuration
DB_PATH=$DB_PATH

# Server configuration
APP_PORT=$APP_PORT

# Paths configuration
TEMPLATE_PATH=$TEMPLATE_PATH
OUTPUT_PATH=$OUTPUT_PATH
ADMIN_FILES_PATH=$ADMIN_FILES_PATH
EOF2

echo "Configuration file created."
echo

# Copy binary
echo "Copying binary to ~/bin/..."
cp personal-blog-generator ~/bin/
chmod +x ~/bin/personal-blog-generator
echo "Binary copied and made executable."
echo

# Copy templates
echo "Copying templates..."
cp -r templates/* "$TEMPLATE_PATH"/
echo "Templates copied."
echo

# Copy admin-files
echo "Copying admin files..."
cp -r admin-files/* "$ADMIN_FILES_PATH"/
echo "Admin files copied."
echo

echo "Setup complete!"
echo
echo "To run the blog generator:"
echo "  ~/bin/personal-blog-generator"
echo
echo "Your blog will be generated in: $OUTPUT_PATH"
echo "Admin interface will be available at: http://localhost:$APP_PORT/admin/"
echo
echo "Next steps:"
echo "1. Run the binary: ~/bin/personal-blog-generator"
echo "2. Open your browser to http://localhost:$APP_PORT/admin/ to manage your blog"
echo "3. Publish your site from the admin interface"
