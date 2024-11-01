#!/usr/bin/env bash

# Check that new user can be successfully registered
curl -d \
    "{\"Name\":\"John Doe 1\",\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/register
echo ""

