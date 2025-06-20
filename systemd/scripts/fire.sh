#!/bin/sh

# Remove the old file.
rm /home/drone/development/kindling/trmnl/*

# Run Kindling
/home/drone/development/kindling/kindling

# Convert and Deploy the images to the site.
convert /home/drone/development/kindling/trmnl/weather.png -monochrome -colors 2 -depth 1 -strip png:/var/www/sploo.sh/trmnl/weather.png

