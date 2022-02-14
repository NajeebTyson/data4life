#!/bin/bash
psql postgresql://postgres:postgres123@postgres:5432/data4life -c "CREATE TABLE tokens (token VARCHAR(7) NOT NULL UNIQUE);"
