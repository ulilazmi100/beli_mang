import os
import subprocess
import argparse
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger()

def load_env_file(env_file_path):
    logger.info("Loading environment variables from file...")
    if os.path.exists(env_file_path):
        with open(env_file_path) as f:
            for line in f:
                if line.startswith('export'):
                    key_value = line.split(' ', 1)[1].strip()
                    key, value = key_value.split('=', 1)
                    os.environ[key] = value
                    logger.info(f"Set {key}={value}")
    else:
        raise FileNotFoundError(f"Environment file not found at path: {env_file_path}")

def run_migration(direction, db_type, db_user, db_password, db_host, db_port, db_name, migration_path, verbose):
    db_uri = f"{db_type}://{db_user}:{db_password}@{db_host}:{db_port}/{db_name}?sslmode=disable"

    verbosity = ""
    if verbose:
        verbosity = "-verbose"

    command = ["migrate", "-database", db_uri, "-path", migration_path, verbosity, direction]
    

    if direction == "down":
        logger.info("Starting down migration...")
        # Automate the confirmation by sending "y\n" to the subprocess
        process = subprocess.Popen(command, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        output, error = process.communicate(input=b'y\n')
        logger.info(output.decode())
        if error:
            logger.error(error.decode())
    elif direction == "up":
        logger.info("Starting up migration...")
        process = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        logger.info(process.stdout.decode())
        if process.stderr:
            logger.error(process.stderr.decode())
    else:
        raise ValueError("Invalid migration direction. Use 'up' or 'down'.")

    logger.info(f"{direction.capitalize()} migration completed.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Database Migration Script")
    parser.add_argument('-dbType', type=str, default="postgres", help="Database type (default: postgres)")
    parser.add_argument('-dbUser', type=str, default="postgres", help="Database user (default: postgres)")
    parser.add_argument('-dbPassword', type=str, default="postgres", help="Database password (default: postgres)")
    parser.add_argument('-dbHost', type=str, default="localhost", help="Database host (default: localhost)")
    parser.add_argument('-dbPort', type=int, default=5432, help="Database port (default: 5432)")
    parser.add_argument('-dbName', type=str, default="postgres", help="Database name (default: postgres)")
    parser.add_argument('-migrationPath', type=str, default="db/migrations", help="Path to migration files (default: db/migrations)")
    parser.add_argument('-setEnvVars', action='store_true', help="Set environment variables from .env file")
    parser.add_argument('-envFilePath', type=str, default="./.env", help="Path to .env file (default: ./.env)")
    parser.add_argument('-verbose', action='store_true', help="Run migrations in verbose mode")
    parser.add_argument('-help', action='store_true', help="Show help")

    args = parser.parse_args()

    if args.help:
        parser.print_help()
        exit()

    if args.setEnvVars:
        load_env_file(args.envFilePath)

    logger.info("Starting database migration script...")

    run_migration("down", args.dbType, args.dbUser, args.dbPassword, args.dbHost, args.dbPort, args.dbName, args.migrationPath, args.verbose)
    logger.info("Waiting before starting up migration...")
    run_migration("up", args.dbType, args.dbUser, args.dbPassword, args.dbHost, args.dbPort, args.dbName, args.migrationPath, args.verbose)

    logger.info("Database migration script completed.")
