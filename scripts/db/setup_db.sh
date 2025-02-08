#!/bin/bash

OPTIONS="-config=./dbconfig.yml"
ENVS=$(cat ../../.env | xargs)

# Migration up
echo "Migration up"
env $ENVS sql-migrate up $OPTIONS
env $ENVS sql-migrate status $OPTIONS

# Seeder up
echo "Seeder up"
env $ENVS sql-migrate up $OPTIONS -env seeder
env $ENVS sql-migrate status $OPTIONS -env seeder
