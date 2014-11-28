#!/bin/bash

shopt -s nullglob

for file in testset/*.test
do
	echo "Test $file"

	./executor $file

	if [ $? -ne 0 ]; then
		echo "Error detected."

		echo "Reduce original file to its minimum."

		tavor --format-file vending.tavor reduce --input-file $file --exec "./executor TAVOR_DD_FILE" --exec-argument-type argument --exec-exact-exit-code > $file.reduced

		echo "Saved to $file.reduced"

		break
	fi
done
