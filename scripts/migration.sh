#!/usr/bin/env bash

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATION_RUN="docker compose --project-directory $REPO_ROOT -f $REPO_ROOT/docker-compose.dev.yml run --rm --build platform-api-migration"

# Function to display usage
usage() {
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  up                    Run migrations (all modules or specific)"
    echo "  down                  Rollback migrations (all modules or specific)"
    echo "  status                Check migration status for all modules"
    echo "  create                Create a new migration file"
    echo ""
    echo "Options:"
    echo "  --module=<name>       Specify a module (for up/down/status commands)"
    echo "  --dry-run             Run in dry-run mode (for up/down commands)"
    echo "  --direction=<dir>     Migration direction: up or down (for status command)"
    echo "  --name=<name>         Migration name (for create command)"
    echo ""
    echo "Examples:"
    echo "  $0 up                                     # Run all migrations"
    echo "  $0 up --module=account                    # Run migrations for account module"
    echo "  $0 up --dry-run                          # Check all migrations (dry-run)"
    echo "  $0 up --module=account --dry-run         # Check account migrations (dry-run)"
    echo "  $0 down                                   # Rollback all migrations"
    echo "  $0 down --module=workflows                # Rollback workflows migrations"
    echo "  $0 status                                 # Check status for all modules"
    echo "  $0 status --direction=down                # Check rollback status for all modules"
    echo "  $0 status --module=account                # Check status for account module"
    echo "  $0 status --module=account --direction=down # Check rollback status for account module"
    echo "  $0 create --name=add_user_table --module=account  # Create new migration"
    exit 1
}

# Parse command
COMMAND=$1
shift || usage

# Parse options
MODULE=""
DRY_RUN=""
NAME=""
DIRECTION=""

for arg in "$@"; do
    case $arg in
        --module=*)
            MODULE="${arg#*=}"
            ;;
        --dry-run)
            DRY_RUN="--dry-run"
            ;;
        --name=*)
            NAME="${arg#*=}"
            ;;
        --direction=*)
            DIRECTION="${arg#*=}"
            ;;
        *)
            echo -e "${RED}Unknown option: $arg${NC}"
            usage
            ;;
    esac
done

# Execute command
case $COMMAND in
    up)
        if [ -n "$MODULE" ]; then
            echo -e "${GREEN}Running migrations for module: $MODULE${NC}"
            if [ -n "$DRY_RUN" ]; then
                echo -e "${YELLOW}Mode: Dry-run${NC}"
            fi
            DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=up ${MODULE:+--module=$MODULE} $DRY_RUN
        else
            if [ -n "$DRY_RUN" ]; then
                echo -e "${YELLOW}Checking migrations for all modules (dry-run)...${NC}"
                DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=up $DRY_RUN
            else
                echo -e "${GREEN}Running migrations for all modules...${NC}"
                DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=up
            fi
        fi
        ;;

    down)
        if [ -n "$MODULE" ]; then
            echo -e "${YELLOW}Rolling back migrations for module: $MODULE${NC}"
            if [ -n "$DRY_RUN" ]; then
                echo -e "${YELLOW}Mode: Dry-run${NC}"
            fi
            DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=down ${MODULE:+--module=$MODULE} $DRY_RUN
        else
            if [ -n "$DRY_RUN" ]; then
                echo -e "${YELLOW}Checking rollback for all modules (dry-run)...${NC}"
                DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=down $DRY_RUN
            else
                echo -e "${YELLOW}Rolling back migrations for all modules...${NC}"
                DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=down
            fi
        fi
        ;;

    status)
        DIRECTION="${DIRECTION:-up}"
        if [ -n "$MODULE" ]; then
            echo -e "${GREEN}Checking migration status for module: $MODULE (direction: $DIRECTION)...${NC}"
            DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=$DIRECTION --module=$MODULE --dry-run
        else
            echo -e "${GREEN}Checking migration status for all modules (direction: $DIRECTION)...${NC}"
            DOCKER_BUILDKIT=1 $MIGRATION_RUN --direction=$DIRECTION --dry-run
        fi
        ;;

    create)
        if [ -z "$NAME" ]; then
            echo -e "${RED}Error: --name is required for create command${NC}"
            usage
        fi

        if [ -z "$MODULE" ]; then
            echo -e "${RED}Error: --module is required for create command${NC}"
            usage
        fi

        MIGRATION_PATH="$REPO_ROOT/apps/platform-api/internal/$MODULE/infrastructure/database/migrations"
        TIMESTAMP=$(date +%Y%m%d%H%M%S)

        mkdir -p "$MIGRATION_PATH"

        UP_FILE="$MIGRATION_PATH/${TIMESTAMP}_${MODULE}_${NAME}.up.sql"
        DOWN_FILE="$MIGRATION_PATH/${TIMESTAMP}_${MODULE}_${NAME}.down.sql"

        echo "-- Migration: $NAME" > "$UP_FILE"
        echo "-- Created: $(date)" >> "$UP_FILE"
        echo "" >> "$UP_FILE"
        echo "-- Write your UP migration SQL here" >> "$UP_FILE"

        echo "-- Rollback: $NAME" > "$DOWN_FILE"
        echo "-- Created: $(date)" >> "$DOWN_FILE"
        echo "" >> "$DOWN_FILE"
        echo "-- Write your DOWN migration SQL here" >> "$DOWN_FILE"

        echo -e "${GREEN}✅ Created migration files:${NC}"
        echo -e "${BLUE}  ↑ UP:   $UP_FILE${NC}"
        echo -e "${BLUE}  ↓ DOWN: $DOWN_FILE${NC}"
        echo ""
        echo -e "${YELLOW}Next steps:${NC}"
        echo "  1. Edit the UP migration file to add your schema changes"
        echo "  2. Edit the DOWN migration file to add rollback logic"
        echo "  3. Run: ${BLUE}make migration-up module=$MODULE${NC} to test"
        ;;

    *)
        echo -e "${RED}Unknown command: $COMMAND${NC}"
        usage
        ;;
esac