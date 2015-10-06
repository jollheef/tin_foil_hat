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
