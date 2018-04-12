#!/usr/bin/env bash

echo "start clean log file"

cd /work/btcpool/build
cd run_gbtmaker
rm log_gbtmaker/*
echo "clean gbtmaker log complete"

cd /work/btcpool/build
cd run_jobmaker
rm log_jobmaker/*
echo "clean jobmaker log complete"

cd /work/btcpool/build
cd run_sserver
rm log_sserver/*
echo "clean sserver log complete"

cd /work/btcpool/build
cd run_sharelogger
rm log_sharelogger/*
echo "clean sharelogger log complete"

cd /work/btcpool/build
cd run_slparser
rm log_slparser/*
echo "clean slparser log complete"

cd /work/btcpool/build
cd run_statshttpd
rm log_statshttpd/*
echo "clean statshttpd log complete"

