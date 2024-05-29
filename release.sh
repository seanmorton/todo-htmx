#!/bin/bash
docker tag seanmorton/todo-htmx:latest seanmorton/todo-htmx:"$1" &&
docker push seanmorton/todo-htmx:latest &&
docker push seanmorton/todo-htmx:"$1" &&
git tag "$1" &&
git push --tags
