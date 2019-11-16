#!/bin/python3

import os,time


def testExpand():
	os.system("go run ./treecli -create-trie 1")
	print("Created trie")
	time.sleep(5)
	os.system("go run ./treecli -insert -key 230 -value \"hahah\" -id 545 -token 37")
	print("Inserted first value")
	time.sleep(5)
	os.system("go run ./treecli -insert -key 234 -value \"PaulskleinerPenis\" -id 545 -token 37")
	print("Inserted second value")
	time.sleep(5)
	os.system("go run ./treecli -insert -key 235 -value \"PaulskleinerPenis\" -id 545 -token 37")
	print("Inserted third value")
	time.sleep(5)
	os.system("go run ./treecli -traverse -id 545 -token 37")
	time.sleep(2)


def main():
	print("Welcome to trietestingScript")
	print("Current Test will be Expand and Travserse after expanding")
	testExpand()


if __name__ == "__main__":
	main()

