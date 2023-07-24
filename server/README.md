# Server for handling shelter donation management

## Requirements:

- go: https://go.dev/doc/install or `brew install go` if you have `brew` installed.

For reference, I am using `go` version `go1.20.6 linux/amd64`

## How to run:

From the `openapi-seattle-example/server` directory, run:

```
go run cmd/main.go
```

which should also download any missing dependencies. The server will then run on localhost:3333

Then for the we can `curl` things like:

`register`:

```
curl localhost:3333/donations/register \
    -H "content-type: application/json" \
    -d '{"name": "Andrew", "type": "money", "quantity": 3, "description": "hello"}'
```

which returns:

```
{"date":"0001-01-01T00:00:00Z","description":"hello","id":1,"name":"Andrew","quantity":3,"type":"money"}
```

`distribute`:

```
curl localhost:3333/donations/distribute \
    -H "content-type: application/json" \
    -d '{"donation_id": 1, "type": "money", "quantity": 3, "description": "hello"}'
```

which returns:

```
{"date":"2023-07-23T04:23:19.578944529-07:00","description":"hello","donation_id":1,"id":1,"quantity":3,"type":"money"}
```

`inventory report`:

```
curl localhost:3333/donations/report/inventory
```

which returns:

```
{"clothing":[],"food":[],"money":[{"distributions":[{"date":"2023-07-24T02:50:41.362104077-07:00","description":"hello","donation_id":1,"id":1,"quantity":3,"type":"money"}],"donation":{"date":"2023-07-24T02:50:38.811969141-07:00","description":"hello","id":1,"name":"Andrew","quantity":3,"type":"money"}}]}
```

`donor report`:

```
{"report":[{"donations":[{"quantity":3,"quantity_distributed":3,"type":"money"}],"name":"Andrew"}]}
```

## Design choices/things I considered:

- `/donations/register`: This should pretty much accept anything that gets put in. I thought it was missing a `description` because otherwise it's hard to distinguish between types of food/clothing. We use ids to distinguish between different donations. This could be a uuid, but we use an auto-incrementing id for ease of implementation.
- `/donations/distribute`: I felt like this should be itemized and distributions should be associated with the original donation so that it's easier to trace when things come from. I also added a `description` for the same reasons as above and it could also include information about where this is going to so that people who donate have a better idea of where their donations are going. We use ids to distinguish between different distributions. This could be a uuid, but we use an auto-incrementing id for ease of implementation..
- `/donations/report/inventory`: The instructions were a little unclear about this, so I split this by type and then by the status of each donation registration/distribution. Alternatively this could've just been a summary of everything for each type of thing, but that doesn't make a lot of sense in the case that you have different types of donations for like clothes. Like you don't want to say that 1 sock and 1 sweatshirt means that you have 2 clothes. Stuff like that.
- `/donations/report/donors`: This just said summary, so I just grouped by types for each donor.

## Technology choices:

- `api definition style`: I chose `OpenAPI` for this because it documents interactions with the backend in a way that is consolidated to a single place. Although I spent a while fiddling around to get this to generate server code I could use and update. The most popular `OpenAPI` generation tool is a little problematic to use because it generates the server code and expects you to edit those files and add it to a `.ignore` file. The problem with this is that this makes it hard to edit or update existing APIs without manually determining the changes you should make based on your existing code and the newly generated one. The nicest method I've seen for using `OpenAPI` for servers involves using interfaces to define the servers and implementing them in a separate file so that it doesn't get overwritten on updates to the generated code.
- `language`: I chose `go`. This could've been anything, but I've been writing `go` so I thought this would be faster for me to write at the moment. I would normally write this in `Python` with `typing`, but I didn't like the generated python server code for the same reasons as I listed above.
- `storage`: I used `psql` because it is lightweight and faster for iterating. Other things I've considered include just writing things to and reading from a file which would've also been fine for this.
- `interactions with the server`: I'm just using `curl` for now because it is faster for iterating locally, at the cost of having a lot of copied and pasted commands.
- `testing`: I'm just doing this in `python` because it is faster for me.
- `logging`: I used `zap` for structured logging, although I think that maybe normal logging would've been better for testing things locally because it'd be more readable. Structured logging is good if we have like a logs processor that needs the logs to be structured, but going that far is probably a bit much for this assignment.
