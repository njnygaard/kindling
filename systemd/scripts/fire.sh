#!/bin/sh

cd /home/drone/development/kindling
go build -o /home/drone/development/kindling -buildvcs=false
/home/drone/development/kindling/kindling
convert /home/drone/development/kindling/weather.png -monochrome -colors 2 -depth 1 -strip png:/var/www/sploo.sh/weather.png
