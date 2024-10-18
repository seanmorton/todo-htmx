#!/bin/bash
if [ "$#" -ne 1 ]; then
  echo "verison tag is required"
  exit 1
fi

docker build --tag seanmorton/todo-htmx:latest --platform linux/amd64 . &&
docker tag seanmorton/todo-htmx:latest seanmorton/todo-htmx:"$1" &&
docker push seanmorton/todo-htmx:latest &&
docker push seanmorton/todo-htmx:"$1" &&
git tag "$1" &&
git push --tags
