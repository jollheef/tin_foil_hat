[![Build Status](https://travis-ci.org/jollheef/tin_foil_hat.svg?branch=master)](https://travis-ci.org/jollheef/tin_foil_hat)
[![GoDoc](https://godoc.org/github.com/jollheef/tin_foil_hat?status.svg)](http://godoc.org/github.com/jollheef/tin_foil_hat)
[![Coverage Status](https://coveralls.io/repos/jollheef/tin_foil_hat/badge.svg?branch=master&service=github)](https://coveralls.io/github/jollheef/tin_foil_hat?branch=master)
[![Go Report Card](http://goreportcard.com/badge/jollheef/tin_foil_hat)](http://goreportcard.com/report/jollheef/tin_foil_hat)

# Tin foil hat
Unix-way contest checking system.

### Components
#### Counter
Count scoreboard.
#### Checker
Manage services checkers.
#### Receiver
Read flags from teams.
#### Steward
Generic database interface.
#### Vexillary
Generate and check flags.
#### Pulse
Manage rounds.
#### Scoreboard
Web scoreboard.

# Deploy

### Depends

    $ emerge dev-db/postgresql

### Build

    $ ./build.sh && ./test.sh

### Run

    $ sudo psql -U postgres
    postgres=# CREATE DATABASE tinfoilhat;
    postgres=# CREATE USER tfh WITH password 'STRENGTH_PASSWORD'
    postgres=# GRANT ALL privileges ON DATABASE tinfoilhat TO tfh;

After that you need to fix 'connection' parameter in configuration file.
(And other parameters, of course)

Now, run it!

    $ ./bin/tinfoilhat ./src/tinfoilhat/config/tinfoilhat.toml
