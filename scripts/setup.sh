#!/bin/bash

# Data Lake - Setup Script
# This script helps with initial setup of the project

set -e

echo "╔════════════════════════════════════════════╗"
echo "║   Data Lake - Initial Setup Script        ║"
echo "╚════════════════════════════════════════════╝"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if required tools are installed
check_requirements() {
    echo "Checking requirements..."

    local missing_tools=()

    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    else
        print_success "Go is installed ($(go version))"
    fi

    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    else
        print_success "Docker is installed"
    fi

    if ! command -v docker-compose &> /dev/null; then
        missing_tools+=("docker-compose")
    else
        print_success "Docker Compose is installed"
    fi

    if ! command -v psql &> /dev/null; then
        print_warning "psql is not installed (optional, but recommended for DB management)"
    else
        print_success "psql is installed"
    fi

    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        echo ""
        echo "Please install missing tools:"
        echo "  - Go: https://golang.org/doc/install"
        echo "  - Docker: https://docs.docker.com/get-docker/"
        echo "  - Docker Compose: https://docs.docker.com/compose/install/"
        exit 1
    fi

    echo ""
}

# Setup .env file
setup_env() {
    if [ -f .env ]; then
        print_warning ".env file already exists"
        read -p "Do you want to overwrite it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Skipping .env setup"
            return
        fi
    fi

    echo "Creating .env file from template..."
    cp .env.example .env
    print_success ".env file created"

    echo ""
    print_info "Please edit .env file and fill in your credentials:"
    echo "  1. WakaTime API credentials (CLIENT_ID, CLIENT_SECRET)"
    echo "  2. Google OAuth credentials (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET)"
    echo "  3. Generate API_KEY: openssl rand -hex 32"
    echo "  4. Generate API_USER_ID: uuidgen"
    echo ""

    # Generate API_KEY and UUID
    read -p "Do you want to generate API_KEY and UUID now? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        if command -v openssl &> /dev/null; then
            API_KEY=$(openssl rand -hex 32)
            sed -i.bak "s/your_secret_api_key_here_generate_with_openssl/$API_KEY/" .env
            print_success "Generated and set API_KEY"
        else
            print_warning "openssl not found, please generate API_KEY manually"
        fi

        if command -v uuidgen &> /dev/null; then
            USER_ID=$(uuidgen | tr '[:upper:]' '[:lower:]')
            sed -i.bak "s/00000000-0000-0000-0000-000000000001/$USER_ID/" .env
            print_success "Generated and set API_USER_ID: $USER_ID"
        else
            print_warning "uuidgen not found, please generate UUID manually"
        fi

        rm -f .env.bak
    fi

    echo ""
}

# Install Go dependencies
install_dependencies() {
    echo "Installing Go dependencies..."
    go mod download
    go mod verify
    print_success "Go dependencies installed"
    echo ""
}

# Start PostgreSQL
start_database() {
    read -p "Do you want to start PostgreSQL with Docker? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        echo "Starting PostgreSQL..."
        docker-compose -f docker.yaml up -d

        # Wait for PostgreSQL to be ready
        echo "Waiting for PostgreSQL to be ready..."
        sleep 5

        if docker-compose -f docker.yaml ps | grep -q "Up"; then
            print_success "PostgreSQL is running"
        else
            print_error "Failed to start PostgreSQL"
            exit 1
        fi
    else
        print_info "Skipping PostgreSQL setup"
        print_warning "Make sure PostgreSQL is running and DATABASE_URL in .env is correct"
    fi
    echo ""
}

# Create database user
create_user() {
    read -p "Do you want to create a user in the database? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        # Get UUID from .env
        USER_ID=$(grep API_USER_ID .env | cut -d '=' -f2 | tr -d '"' | tr -d "'")

        if [ -z "$USER_ID" ]; then
            print_error "API_USER_ID not found in .env"
            return
        fi

        read -p "Enter username: " username
        read -p "Enter email: " email

        echo "Creating user in database..."
        docker exec -i postgres_datalake psql -U postgres -d datalake <<-EOSQL
INSERT INTO users (id, username, email)
VALUES ('$USER_ID', '$username', '$email')
ON CONFLICT (id) DO NOTHING;
EOSQL

        print_success "User created with ID: $USER_ID"
    else
        print_info "Skipping user creation"
    fi
    echo ""
}

# Build application
build_app() {
    read -p "Do you want to build the application? (Y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        echo "Building application..."
        mkdir -p bin
        go build -o bin/datalake cmd/main.go
        print_success "Application built: bin/datalake"
    else
        print_info "Skipping build"
    fi
    echo ""
}

# Final instructions
print_instructions() {
    echo "╔════════════════════════════════════════════╗"
    echo "║          Setup Complete! 🎉                ║"
    echo "╚════════════════════════════════════════════╝"
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. Edit .env file and add your OAuth credentials:"
    echo "   ${YELLOW}nano .env${NC}"
    echo ""
    echo "   Get WakaTime credentials: https://wakatime.com/apps"
    echo "   Get Google credentials: https://console.cloud.google.com/"
    echo ""
    echo "2. Run the application:"
    echo "   ${YELLOW}make run${NC}"
    echo "   or"
    echo "   ${YELLOW}./bin/datalake${NC}"
    echo ""
    echo "3. Authorize with OAuth providers (check logs for URLs)"
    echo ""
    echo "4. Test API:"
    echo "   ${YELLOW}export API_KEY=\$(grep API_KEY .env | cut -d '=' -f2 | tr -d '\"')${NC}"
    echo "   ${YELLOW}curl -H \"X-API-Key: \$API_KEY\" http://localhost:8080/api/v1/wakatime/stats${NC}"
    echo ""
    echo "5. (Optional) Start monitoring:"
    echo "   ${YELLOW}make monitoring-up${NC}"
    echo ""
    echo "For more information, see README.md"
    echo ""
}

# Main setup flow
main() {
    check_requirements
    setup_env
    install_dependencies
    start_database
    create_user
    build_app
    print_instructions
}

# Run main function
main

