#!/usr/bin/env bash

# Check that server is running
curl "http://localhost:8000/api/healthchecker"
echo ""

# Check that user can be successfully registered
curl -d \
    "{\"Name\":\"John Doe 2\",\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/register
echo ""