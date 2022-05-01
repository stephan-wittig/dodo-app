# dodo prototype

This is the prototype for the dodo platform, a project in MG4C3 at LSE.

## Components

´/api´ holds the code for the backend prototype, an application that can load, serve and fill templates for policies.

The frontend was prototyped in figma, and is not contained in thsi repository.

## How to...

### run

To run the backend prototype, (Go must be installed)[https://go.dev/doc/install].

```
cd api
go run .
```

You will be prompted to

1. Select a template to work with
2. Select whether to generate instructions or generate a document
3. (if generating a document) Select a file with variable values.

Examples for all files can be found in `/demo`, and this is where the applications will look for files.
