#!/usr/bin/env bash
set -euo pipefail

TARGET_DIR="/docker"
COMPOSE_FILE="$TARGET_DIR/docker-compose.yml"

if docker info > /dev/null 2>&1; then
    DOCKER_CMD="docker"
else
    echo "User cannot access Docker directly, will use sudo"
    DOCKER_CMD="sudo docker"
fi

if [ -d "$TARGET_DIR" ]; then
    cd "$TARGET_DIR"

    if [ -f "$COMPOSE_FILE" ]; then
        echo "Shutting down docker stack.."
        if [ "$DOCKER_CMD" = "sudo docker" ]; then
            sudo docker-compose down
        else
            docker-compose down
        fi
    else
        echo "docker-compose.yml not found in $TARGET_DIR, skipping down"
    fi

    echo "Removing $TARGET_DIR.."
    cd /
    sudo rm -rf "$TARGET_DIR"

    if $DOCKER_CMD image inspect poda:latest >/dev/null 2>&1; then
        echo "Removing local image.."
        $DOCKER_CMD image rm poda:latest
    fi

    echo "Cleanup complete."
else
    echo "Error: $TARGET_DIR does not exist, nothing to clean."
fi
