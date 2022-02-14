#!/bin/bash
psql postgresql://postgres:postgres123@postgres:5432/data4life -c "DROP TABLE tokens;"
