#!/bin/bash

isMigrationsFolderExist() {
  if [[ ! -d ./../db/migration ]]; then
    return 0
  fi
  return 1
}

if [ ! $(( which migrate | wc -l )) -eq 1 ]; then
  echo "[ERROR]: Please install Golang Database Migrations"
  echo "[ERROR]: Follow instruction on link below"
  echo "[ERROR]: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
  exit 1
fi

echo ""
echo "1. Create Database Migrations"
echo "2. Run Migrations"
echo "3. Rollback Migrations"
echo "4. Drop Migrations"

echo -n "Input [1-4]: "
read menu
echo

if [[ menu -eq 1 ]]; then
  if [[ isMigrationsFolderExist -eq 0 ]]; then
    echo "[ERROR]: Folder migration not exists"
    exit 1
  fi

  echo "Creating database migrations file..."
  
  echo -n "Input the name of migrations file (without space): "
  read fileName

  migrate create -ext sql -dir ../db/$fileName
  exit 0
elif [[ menu -eq 2 ]]; then
  echo "not implemented"
elif [[ menu -eq 3 ]]; then 
  echo "not implemented"
elif [[ menu -eq 4 ]]; then
  echo "not implemented"
else
  echo "The options is no exists"
  exit 1
fi