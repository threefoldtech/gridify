#!/bin/sh

git clone $REPO_URL repo && cd repo
ginit start -f demo/Procfile -e demo/.env

