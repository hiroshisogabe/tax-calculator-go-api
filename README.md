# Tax Calculator Go API

## Why this repo?

It started as an idea to write in `Go` the service that provides the tax calculation in the `Next.js` application [`tax-calculator-nextjs`](https://github.com/hiroshisogabe/tax-calculator-nextjs).

Before writing the code in this repository, I started "A Tour of Go" from the Go official website and decided to push to the [`a-tour-of-go`](https://github.com/hiroshisogabe/a-tour-of-go) repository the resolution of the exercises I could complete.

## Was AI used in this repo?

Yes, Google Gemini 3 was asked to provide some "server" implementation in Go based on the tax calculation that I implemented in the `Next.js` application [`tax-calculator-nextjs`](https://github.com/hiroshisogabe/tax-calculator-nextjs).

The first answer from Google Gemini was the creation of a `main.go` with everything, which is not ideal from my perspective. Nevertheless I decided to keep the suggested `main.go` file and committed to the repository, thus we can evaluate the refactoring which splits some pieces of code into new file(s) or specific functions as helpers.

The refactoring phase also had contribution from Google Gemini after tweaking it a bit, including the test files and scenarios.

In addition Github Copilot using Claude Haiku 4.5 was also helpful to clarify and execute specific tasks because it can see the whole codebase and make adjustments directly.

## Tech notes

### Install

If you don't have `Go` installed in your environment, you can check the [Download and install](https://go.dev/doc/install) page from the official website. Without `go`, nothing will work as you might expect.

### Start

After cloning the repo you can start the server locally with the following command in the root path of the project:
```bash
go run main.go
```

If everything goes well, as we can see in the code snippet below from [`main.go`](./main.go), it'll start the server on port `8080` and show the message:
```go
port := ":8080"
fmt.Printf("Tax Engine running on http://localhost%s\n", port)
```

As the server is available at `http://localhost:8080`, the full path will be `http://localhost:8080/calculate` because as we can see also from [`main.go`](./main.go), the endpoint is provided in the path `/calculate`:
```go
http.HandleFunc("/calculate", calculateHandler)
```

### Test

As the main idea of the repository is to expose an endpoint that calculates the tax based on some input, it's possible to either test the endpoint by sending a request, e.g. [`curl`](#using-curl), when it's running locally or through the written tests with [`go test`](#running-_testgo-files) command.

#### Using `curl`

If the server is running, the following examples use [`curl`](https://curl.se/) to send a `POST` request to the endpoint `/calculate`:

```bash
# Success
# Response: {"success":true,"data":{"baseAmount":100,"taxAmount":8.799999999999999,"total":108.8,"rate":0.088,"state":"NY","year":2024}}
curl -X POST http://localhost:8080/calculate \
    -H "Content-Type: application/json" \
    -d '{
            "amount": 100,
            "state": "ny",
            "year": 2024,
            "productCategory": "electronics"
        }'
```
```bash
# Fail
# Response: {"success":false,"error":"Amount must be greater than zero"}
curl -X POST http://localhost:8080/calculate \
    -H "Content-Type: application/json" \
    -d '{
            "amount": -50,
            "state": "NY",
            "year": 2024,
            "productCategory": "electronics"
        }'
```

> The payload was sent in the `-d` option as JSON.

#### Running `*_test.go` files

The following samples show how to run a specific test file, e.g. the [`calculator_test.go`](./pkg/calculator/calculator_test.go), or all the tests: 

```bash
# Run specific test file
go test ./pkg/calculator -v

# Run all test files
go test ./... -v
```

##### Running specific test case

It's also possible to run a specific test case from a test file due to the fact the tests follow the [TableDrivenTests](https://go.dev/wiki/TableDrivenTests) approach and we provided a name to each case.


```bash
# FindRate test function

# case "Valid NY Rule"
go test -v -run TestFindRate/Valid_NY_Rule ./pkg/calculator 
```

```bash
# Calculate test function

# case: "Standard 10% tax"
go test -v -run TestCalculate/Standard_10%_tax ./pkg/calculator
```

> Details from the commands above:
> - `-v` option stands for `verbose` mode
> - `...` runs tests in the current directory and all sub-directories

## Is there something planned to do next?

### Have a possible client for consuming the endpoint

As the implementation of this repository was based on the tax calculation from the [`tax-calculator-nextjs`](https://github.com/hiroshisogabe/tax-calculator-nextjs) application, it would be awesome to consume the endpoint provided by `go` from there, instead of having the calculation logic in the Server Action.

### Dockerize and deploy to cloud?

It was a suggestion from Google Gemini, create a dockerfile and think about deploying it, given the purpose of exposing an endpoint shouldn't be only to run it locally.

For sure it suggested the [Google Could Run](https://cloud.google.com/) as a "Industry Standard".
