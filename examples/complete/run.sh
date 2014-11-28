#!/bin/bash

shopt -s nullglob

for file in testset/*.test
do
	echo "Test $file"

	./executor $file

	if [ $? -ne 0 ]; then
		echo "Error detected, will exit loop"

		break
	fi
done
