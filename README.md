# Running the Application

## Running with npm and go (Recommended)

In order to run the application locally for development do the following:

1. Install all dependencies.

2. Copy the .env.example files in /application and /application/gui to /application/.env and /application/gui/.env

   ```sh
   cp application/.env.example application/.env && \\
   cp application/gui/.env.example application/gui/.env
   ```

3. Run the entire stack from the project root with:

   ```sh
   npm run start:all
   ```

OR

3. Run the frontend by navigating to /application/gui and running:

   ```sh
   npm ci && \\
   npm run start
   ```

4. In another terminal run the backend by navigating to /application and running:

   ```sh
   go run ./cmd
   ```

## Running with docker

In order to run the application with docker, follow these steps:

1. Disable any locally running mongodb service.

2. navigate to /application and run:

   ```sh
   docker compose up
   ```

# Testing the Application

To test the application and get the coverage percentage, follow these steps:

1. Open a terminal and navigate to the `application` directory:
   ```sh
   .../TKE-passwordless-authentication/application
   ```
2. Then run this command to run the tests located in the `tests` directory and generate code coverage percentage for code in `internal`:
   ```sh
   go test -cover -coverpkg=chalmers/tkey-group22/application/internal ./tests
   ```

# Viewing documentation

To view the documentation in localhost:

1. Install go doc

2. From the `application` directory, run:

   ```sh
   godoc -http=:6060
   ```

3. Navigate to `http://localhost:6060/pkg/chalmers/tkey-group22/application/internal/` to view documentation for the `internal` package.
