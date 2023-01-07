#!/bin/sh

cd /home/drone/development/kindling
go build -o /home/drone/development/kindling
/home/drone/development/kindling/kindling
cp /home/drone/development/kindling/out.png /var/www/kindle/image.png
