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
