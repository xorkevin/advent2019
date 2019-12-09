#!/usr/bin/env bash

basedir=$PWD
days=$(find . -maxdepth 1 -type d -name 'day*' | sort)

for i in $days
do
  echo "BENCHMARKING $i..."
  cd "$basedir/$i"
  make bench
done
