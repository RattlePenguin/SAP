#!/usr/bin/env bash

##### Server Running Tests #####

# PRE-TEST: Wipe database
##### rm golang.db

# Check that server is running
curl "http://localhost:8000/api/healthchecker"
echo ""


##### 2FA Route Testing #####

# Check that new user can be successfully registered
curl -d \
    "{\"Name\":\"John Doe 1\",\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/register
echo ""

# Check that same email address cannot be used
curl -d \
    "{\"Name\":\"John Doe 2\",\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/register
echo ""

# Check that user can be logged in
json=$(curl -d \
    "{\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login)
echo ""
userid=$(echo "$json" | cut -d '"' -f14)


# Generate OTP for User
curl -d \
    "{\"user_id\":\"$userid\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/otp/generate

# Verify Incorrect OTP produces error
read -n6 token
bad_token=$(( 10#$token + 1 ))
curl -d \
    "{\"user_id\":\"$userid\",\"token\":\"$bad_token\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/otp/verify

# Verify OTP for User
curl -d \
    "{\"user_id\":\"$userid\",\"token\":\"$token\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/otp/verify

# Validate Incorrect OTP for User
curl -d \
    "{\"user_id\":\"$userid\",\"token\":\"$bad_token\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/otp/validate

# Validate OTP for User
curl -d \
    "{\"user_id\":\"$userid\",\"token\":\"$token\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/otp/validate


##### Vulnerability Tests #####

# Verify cookie given to already logged in user


# SQL Vulnerability 1: SELECT user with no 2FA enabled
# SQL Vulnerability 2: SELECT specific user with no 2FA enabled
# SQL Vulnerability 3: Disable user's Otp_enabled value

# XSS Vulnerability 1: Stealing Cookies
