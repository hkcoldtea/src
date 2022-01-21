#!/bin/bash
if [ ! -f ~/go/bin/rectcrop ]; then
	exit 0
fi
folder=""
OLDIFS=$IFS
IFS="
"
for filename in "$@"
do
	folder=${folder:-"${filename%.*}"}
	mkdir -p "${folder}"
	~/go/bin/rectcrop -input="$filename" -format=jpeg -output="${folder}/${filename%.*}.jpg"
done
IFS=$OLDIFS
rmdir --ignore-fail-on-non-empty "${folder}"
