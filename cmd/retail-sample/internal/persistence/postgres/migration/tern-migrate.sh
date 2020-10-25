#!/bin/sh

# while ! nc -z $HOST $PORT; do
while ! curl http://$HOST:$PORT/ 2>&1 | grep '52'
do
  echo "waiting for $HOST..."
  sleep 1
done

echo "running tern migrate"

tern migrate

echo "done"