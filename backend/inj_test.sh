#!/usr/bin/env bash

# Correct login details
curl -d \
    "{\"Email\":\"john.doe@email.com\",\"Password\":\"johndoe123\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

# Correct username, incorrect password
curl -d \
    "{\"Email\":\"john.doe@email.com\",\"Password\":\"totallynotreal\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

# Incorrect login details
curl -d \
    "{\"Email\":\"totallyfakeemail@email.com;;;\",\"Password\":\"totallywrongpass\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

curl -d \
    "{\"Email\":\"totallyfakeemail@email.com\";\",\"Password\":\"totallywrongpass\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

curl -d \
    "{\"Email\":\"totallyfakeemail@email.com\" OR 1=1\",\"Password\":\"totallywrongpass\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

curl -d \
    "{\"Email\":\"x\" OR \"1\"=\"1\",\"Password\":\"x\" OR \"1\"=\"1\"}" \
    -H "Content-Type: application/json" http://localhost:8000/api/auth/login
echo ""

# "SELECT * FROM `users` WHERE `users`.`email` = "x" OR "1"="1" AND `users`.`password` = "" ORDER BY `users`.`id` LIMIT 1"