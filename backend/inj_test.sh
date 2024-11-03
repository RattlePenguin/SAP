#!/usr/bin/env bash

# Correct login details
curl -d \
    "{\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

# Incorrect login details
curl -d \
    "{\"Email\":\"totallyfakeemail@email.com OR 1=1\",\"Password\":\"totallywrongpass OR 1=1\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""
