#!/bin/sh

set -x

cd /home/drone/development/kindling
rm kindling
rm weather.png
go build -o /home/drone/development/kindling
/home/drone/development/kindling/kindling
convert /home/drone/development/kindling/weather.png -monochrome -colors 2 -depth 1 -strip png:/var/www/sploo.sh/weather.png

