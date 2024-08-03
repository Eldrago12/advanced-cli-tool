#!/bin/bash

go build -o act main.go

sudo mv act /usr/local/bin/

if command -v act &> /dev/null
then
    echo "act successfully installed!"
else
    echo "Installation failed."
fi
