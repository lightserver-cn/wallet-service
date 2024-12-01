# Define the path to the Docker Compose configuration file
yaml_file="docker-compose/docker-compose.yaml"

# Build the server-side binary package
# Execute the build process located in the dockerfile/build.sh script to compile the server-side application
sh ./dockerfile/build.sh server

# Remove any existing log directory and its contents to ensure a clean state
rm -rf ./runtime/log

# Remove any existing Docker volume directories for PostgreSQL and Redis to ensure a clean state
rm -rf ./docker-compose/volumes/{postgres,redis}

# Create a new log directory for storing log files generated during application runtime
mkdir -p ./runtime/log

# Create new Docker volume directories for persisting PostgreSQL and Redis data
mkdir -p ./docker-compose/volumes/{postgres,redis}

# Stop and remove containers, networks, and volumes defined in the Docker Compose configuration file
docker-compose -f "$yaml_file" down

# Set up the runtime environment using the specified Docker Compose configuration file and run the project
# The -f parameter specifies the configuration file path
# The --build parameter indicates to rebuild the image before starting the container
# The -d parameter runs the container in detached mode (background)
docker-compose -f "$yaml_file" up --build -d
