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
