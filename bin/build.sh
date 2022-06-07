#! /bin/sh
# GOOS=linux GOARCH=amd64 go build -o owl.linux . 
docker build --platform=linux/amd64 -t translucentlink/scoring:$1 .
docker push translucentlink/scoring:$1