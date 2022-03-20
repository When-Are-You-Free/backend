# When Are You Free - Backend

## Concept

Allows you to create and edit a shared calendar for meet-ups, allowing you to
easily find out which of the atteendes have time. There's no account required.

Each user is identified by an HTTP header. It's up to the client to generate
the identifier and make sure it can't easily be bruteforced. In addition to
that, a nickname can be added. Both the name and the identifier will be saved
in the database.

The client should use persistent storage for the identifier, as recovery
requires manually retrieving the identifier from the database.

## Running

To run, simply installed go 1.18 or later and run:

```
go run .
```

## API Docs

There's docs in form of a swagger.yml in the root of the repository.

Might host this on GH pages at some point.
