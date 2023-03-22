#!/bin/bash

# Start the Go server in the background
go run server/main.go &

# Start the SvelteKit frontend server
npm run dev
