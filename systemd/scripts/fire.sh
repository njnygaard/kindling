#!/bin/sh

cd /home/drone/development/kindling
go build -o /home/drone/development/kindling -buildvcs=false
/home/drone/development/kindling/kindling
cp /home/drone/development/kindling/weather.png /var/www/sploo.sh/weather.png
