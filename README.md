# Redis Migrator

Redis Migrator will take the keys from one server/db and
transfer them to another server/db.

## How does it work?

Simply it iterates over the keys in SOURCE and recreates them on DESTINATION.

This will not wipe out any existing data in DESTINATION, with the exception
of Lists which have to be replaced entirely on the DESTINATION. Strings will
obviously be replaced. Hashes will only be updated, any pre-existing hash fields
set on the DESTINATION will remain if they are not replaced by the SOURCE.

Any pre-existing keys on DESTINATION that are not on SOURCE will remain untouched.

## Usage

```bash
Usage: redis-migrator --source-host SOURCE-HOST --destination-host DESTINATION-HOST
  [--source-port SOURCE-PORT] [--source-db SOURCE-DB] [--source-user SOURCE-USER]
  [--source-pass SOURCE-PASS] [--destination-port DESTINATION-PORT]
  [--destination-db DESTINATION-DB] [--destination-user DESTINATION-USER]
  [--destination-pass DESTINATION-PASS] [--source-filter SOURCE-FILTER]
  [--destination-prefix DESTINATION-PREFIX] [--verbose]

Options:
  --source-host SOURCE-HOST
                         source redis server hostname
  --source-port SOURCE-PORT
                         source redis server port [default: 6379]
  --source-db SOURCE-DB
                         source redis server db index [default: 0]
  --source-user SOURCE-USER
                         source redis server auth username
  --source-pass SOURCE-PASS
                         source redis server auth password
  --destination-host DESTINATION-HOST
                         destination redis server hostname
  --destination-port DESTINATION-PORT
                         destination redis server port [default: 6379]
  --destination-db DESTINATION-DB
                         destination redis server db index [default: 0]
  --destination-user DESTINATION-USER
                         destination redis server auth username
  --destination-pass DESTINATION-PASS
                         destination redis server auth password
  --source-filter SOURCE-FILTER
                         source keys filter string [default: *]
  --destination-prefix DESTINATION-PREFIX
                         destination key prefix to prepend
  --verbose, -v          Print a lot more info
  --help, -h             display this help and exit
```

## Install

#### Centos

RPM hosted on [yum.jc21.com](https://yum.jc21.com)

#### Go Get

```bash
go get github.com/jc21/redis-migrator
```


#### Building

```bash
git clone https://github.com/jc21/redis-migrator && cd redis-migrator
go build -o bin/redis-migrator cmd/redis-migrator/main.go
./bin/redis-migrator -h
```
