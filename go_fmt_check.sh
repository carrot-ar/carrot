#!/bin/bash

if [[ $(go fmt) ]]; then
  exit 1
else 
  exit 0
fi
