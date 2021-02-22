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
Usage: redis-migrator --source-host SOURCE-HOST [--source-port SOURCE-PORT]
  [--source-db SOURCE-DB] [--source-user SOURCE-USER] [--source-pass SOURCE-PASS]
  --destination-host DESTINATION-HOST [--destination-port DESTINATION-PORT]
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
  --version              display version and exit
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
go build -ldflags="-X main.version=1.0.1" -o bin/redis-migrator cmd/redis-migrator/main.go
./bin/redis-migrator -h
```


## Real world example

Copying data from one db to another on the same instance. This is a production instance
where the keys were expiring mid-way through the migration, hence why the total
number of keys migrated doesn't match the initial count.

```bash
time redis-migrator --source-host 127.0.0.1 --source-db 0 --destination-host 127.0.0.1 --destination-db 1
SOURCE:
  Server:   127.0.0.1:6379
  DB Index: 0
  Auth:     None
DESTINATION:
  Server:   127.0.0.1:6379
  DB Index: 1
  Auth:     None
2021-02-22 07:08:37 INFO => Source has 622028 keys
2021-02-22 07:08:37 INFO => Destination has 0 keys
2021-02-22 07:08:37 INFO => Found 622028 keys on SOURCE with Key filter: *
2021-02-22 07:08:37 INFO => Migration running, each dot is ~1,000 keys
.......... .......... .......... .......... .......... 8% (50012 / 622028)
.......... .......... .......... .......... .......... 16% (100028 / 622028)
.......... .......... .......... .......... .......... 24% (150037 / 622028)
.......... .......... .......... .......... .......... 32% (200046 / 622028)
.......... .......... .......... .......... .......... 40% (250058 / 622028)
.......... .......... .......... .......... .......... 48% (300073 / 622028)
.......... .......... .......... .......... .......... 56% (350084 / 622028)
.......... .......... .......... .......... .......... 64% (400090 / 622028)
.......... .......... .......... .......... .......... 72% (450102 / 622028)
.......... .......... .......... .......... .......... 80% (500113 / 622028)
.......... .......... .......... .......... .......... 88% (550127 / 622028)
.......... .......... .......... .......... .......... 96% (600147 / 622028)
.......... .......... ..
2021-02-22 07:16:33 INFO => Migration completed with 622004 keys, 0 skipped :)

real    7m55.358s
user    0m52.892s
sys     4m32.788s
```
