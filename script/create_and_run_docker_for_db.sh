#!/bin/bash

if [[ $(which docker | wc -l) -eq 0  ]]; then
	echo "Docker engine doesn't exist"
	echo "Please install docker first"
	exit 1
fi

echo "Check if theres sumber_sari_volume volume exist"
echo
if [[ $(docker volume ls | grep sumber_sari_volume | wc -l) -eq 0 ]]; then
  echo "sumber_sari_volume volume doesn't exist yet"
  echo "Creating the volume"
  echo
  docker volume create sumber_sari_volume
  echo
  echo "Volume created"
else
  echo "sumber_sari_volume was exist"
fi
echo

echo "Check if theres sumber_sari_db container exist"
echo
if [[ $(docker ps | grep sumber_sari_db | wc -l) -eq 0 ]]; then
  echo "sumber_sari_db container doesn't exist yet"
  echo "Creating the container"
  echo
  docker run -d --name sumber_sari_db \
  -v sumber_sari_volume:/var/lib/mysql \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  -e MYSQL_ROOT_PASSWORD= \
  -e MYSQL_DATABASE=sumber_sari_db \
  -p 3308:3306 \
  mysql:8
  echo
  echo "Container created"
else
  echo "sumber_sari_db container was exist"
fi
echo

echo "Preparing db container success"
echo "Happy hacking!"

exit 0