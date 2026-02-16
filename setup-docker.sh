#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="/docker"
MARIADB_DIR="/docker/mariadb"
COMPOSE_FILE="$TARGET_DIR/docker-compose.yml"
SOURCE_COMPOSE="$SCRIPT_DIR/docker-compose.yml"
ENV_FILE="$TARGET_DIR/.env"

copied_compose=0

if docker info > /dev/null 2>&1; then
    DOCKER_CMD="docker"
else
    echo "User cannot access docker directly, will use sudo"
    DOCKER_CMD="sudo docker"
fi

if [ ! -d "$TARGET_DIR" ]; then
    echo "Creating $TARGET_DIR"
    sudo mkdir -p "$TARGET_DIR"
    sudo chown "$USER":"$USER" "$TARGET_DIR"
fi

if [ ! -d "$MARIADB_DIR" ]; then
    echo "Creating $MARIADB_DIR"
    sudo mkdir -p "$MARIADB_DIR"
    sudo chown 1000:1000 "$TARGET_DIR"
    sudo chown 1000:1000 "$MARIADB_DIR"
fi

if [ ! -f "$COMPOSE_FILE" ]; then
    echo "docker-compose.yml not found in $TARGET_DIR, copying from script directory"
    sudo cp "$SOURCE_COMPOSE" "$COMPOSE_FILE"
    sudo chown "$USER":"$USER" "$COMPOSE_FILE"
    copied_compose=1
fi

if [ ! -f "$ENV_FILE" ]; then
    echo "Generating $ENV_FILE"
    MARIADB_ROOT_PASSWORD=$(openssl rand -base64 24)
    MARIADB_PASSWORD=$(openssl rand -base64 24)

    cat > "$ENV_FILE" <<EOF
MARIADB_ROOT_PASSWORD=$MARIADB_ROOT_PASSWORD
MARIADB_DATABASE=poda
MARIADB_USER=poda
MARIADB_PASSWORD=$MARIADB_PASSWORD
TZ=Europe/Athens
EOF
fi

if [ "$copied_compose" -eq 1 ]; then
    echo "Checking docker availability.."
    if ! $DOCKER_CMD info > /dev/null 2>&1; then
        echo "Error: Docker is not running or not accessible even with sudo"
        exit 1
    fi

    echo "Checking docker-compose availability.."
    if ! command -v docker-compose > /dev/null 2>&1; then
        echo "Error: docker-compose not found in PATH"
        exit 1
    fi

    echo "Starting up docker stack.."
    cd "$TARGET_DIR"
    if [ "$DOCKER_CMD" = "sudo docker" ]; then
        sudo docker-compose up -d
    else
        docker-compose up -d
    fi
fi

echo "Done!"
