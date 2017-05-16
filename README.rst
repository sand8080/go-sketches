==============
Disjoints test
==============

Project purpose
---------------

* Test of disjoint set finding implementation on the large data
  (~20M of objects and ~200M relations)

Components
----------

Library with disjoint set finding implementation.
Manage script for drop/create/fill data into DB and finding the disjoint sets.

How to configure development environment
----------------------------------------

To start testing on a localhost we need to have:

* PostgreSQL_ of version 9.3 or greater,
* Golang_ of version 1.8.1,

Install PostgreSQL_ and development libraries: ::

  sudo apt-get install --yes postgresql postgresql-server-dev-all


Install golang as described in https://golang.org/doc/install

Install required libraries: ::

    go get github.com/lib/pq golang.org/x/crypto/ssh/terminal gopkg.in/urfave/cli.v1

Create a DB user for test: ::

  sudo -u postgres psql
  CREATE ROLE disjoint WITH NOSUPERUSER NOCREATEDB NOCREATEROLE LOGIN ENCRYPTED PASSWORD 'disjoint';

Create a DB and grant privileges to it: ::

  sudo -u postgres psql
  CREATE DATABASE disjoint;
  GRANT ALL ON DATABASE disjoint TO disjoint;

Check that all tests are passed: ::

  cd go-sketches && go test ./...

Build the executable: ::

  go build

How to deal with DB
-------------------

DB connection credentials should be passed as params for go-sketches executable.
Default params are: ::

    user: disjoint
    password: disjoint
    dbname: disjoint
    host: localhost
    port: 5432

Create new tables in DB: ::

  go-sketches create

Drop tables from DB: ::

  go-sketches drop

Fill data into tables: ::

  go-sketches fill


Check the help for the 'fill' command with option '--help' for creating required DB content.

How to run tests
----------------

For running the tests execute: ::

  go-sketches recalc

Result will be stored into 'components' table.

For instance creation of DB with 20M objects and 200M relations can be done with: ::

  go-sketches drop && go-sketches create && go-sketches fill -n 20000000 -r 200000000 -m 0 -M 10

DB data generation takes about ~1.5 hours on the Intel core-i5 with 4 cores with 16Gb of RAM.
Recalculation of the disjoint sets takes about 25 minutes.

.. _PostgreSQL: https://www.postgresql.org/
.. _Golang: https://golang.org
