#!/bin/bash

set -e

cat test.eml | /app --test | grep "5065427c-23d3-47ca-b6e0-946ea0e8c4be"

