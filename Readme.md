# Task

Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the [Prof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

# Solution
A simple protocol over TCP was implemented, to allow client and server communicate errors, messages and allow client to understand if it should solve the POW challenge. 
On connection, server may send a challenge to a client. The client should respond with the solution. If solution is correct, server sends a quote.

I don't have good enough mathematics background to decide what POW algorithm would suite best, so more effort was spent on the extensibility, code readability and tests.

## POW challenge choice
I looked through CPU-bound, memory-bound and network-bound algorithms and stopped on CPU/GPU-bound algorithm that is based on the SHA-256 hash function that is basically the same as in the Bitcoin.
It might be not the best choice, however I decided to use it because:
 - It is easy to understand and implement;
 - It is widely used, so I have more trust that it can't be exploited;
 - Difficulty can be easily varied to make it harder for fraud clients.

CPU-bound POW challenge could be more harmful for low hardware clients, so server-client protocol was designed in such a way to allow not to send challenge to legitimate clients and also increase challenge difficulty for fraud clients.

One other algorithm that I strongly considered is network-bound [Guided tour puzzle protocol](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol). However, I decided that it doesn't make sense to use it for this task as we have only one server and use the same machine for guide servers would only increase server load.  

# How to run

## Docker
### Build

```bash
docker build -f docker/server.dockerfile -t wisdom-server:latest .
docker build -f docker/client.dockerfile -t wisdom-client:latest .
```

### Run

```bash
docker run --rm -it --name=wisdom-server wisdom-server:latest 0.0.0.0:8081
docker run --rm -it wisdom-client:latest $(docker inspect -f '{{- printf "%s:%d" .NetworkSettings.IPAddress 8081 -}}' wisdom-server)
```

## From source
```bash
go run cmd/server/main.go
go run cmd/client/main.go
```

# Tests
## Run unit tests
```bash
go test intertal/...
```
## Run integration tests
```bash
go test tests/...
```

## Run fuzzy tests
```bash
go test -v -fuzz=FuzzWriteOK -fuzztime=30s ./internal/protocol/
go test -v -fuzz=FuzzSolve -fuzztime=30s ./internal/challenge
```

# Project structure

`cmd` - folder with client and server entrypoints\
`internal` - folder where all implementation is located\
`tests` - folder for integration/functional tests\
`docker` - folder with Dockerfiles
