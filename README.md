# Simple HTTP session checker

## Table of Contents

- [About](#about)

## About <a name = "about"></a>

Simple HTTP session (cookie, redis) checker. By GET /incr, you will get the following response.

{"count":3,"hostname":"cf238ff21641"}

count will increase each time you GET /incr depending on the session.
