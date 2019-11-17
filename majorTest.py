#!/bin/python3

'''
This script is a set of commands for a prasentation on this tree.
It is quiet short and will show some of the features of this trie.
'''

import os, time, sys
from keyboard import press_and_release

def TestInsert():
	'''
	Function to Test the Insert-Function of our Trie.
	'''
	os.system("go run ./treecli -insert -key 230 -value \" Das hier wird mein erster Knoten \" -id 545 -token 37")
	print("[+] Inserted first value")
	time.sleep(2.5)
	os.system("go run ./treecli -insert -key 342 -value \" Das gefällt mir sehr \" -id 545 -token 37")
	print("[+]Inserted second value --> Right Side")
	time.sleep(2.5)
	os.system("go run ./treecli -insert -key 12 -value \" Ist das Produkt aus 3x4 \" -id 545 -token 37")
	print("[+]Inserted third value --> Left Side")
	time.sleep(2.5)
	os.system("go run ./treecli -insert -key 25000 -value \" Wäre ein schönes Sümmchen Geld \" -id 545 -token 37")
	print("[+] Inserted 4th value --> Right Side")
	time.sleep(2.5)
	os.system("go run ./treecli -insert -key 226 -value \" Eine Interessante Zahl \" -id 545 -token 37")
	print("[+] Inserted 4th value --> Left Side")
	time.sleep(2.5)
	os.system("go run ./treecli -insert -key 228 -value \" Ist zwei größer als die Zahl zuvor \" -id 545 -token 37")
	print("[+] Inserted 5th value --> Left Side")
	time.sleep(2.5)

def TestChange():
	'''
	Function to change some values 
	'''
	os.system("go run ./treecli -insert -key 25000 -value \" Eigentlich will ich doch noch mehr Geld haben \" -id 545 -token 37")
	print("[+] Change value 25000")
	time.sleep(2.5)

def TestDelete():
	'''
	Function to Test the Delete Command of our Trie.
	'''
	os.system("go run ./treecli -delete 342 -id 545 -token 37")
	print("[-->]  Delete Value")
	time.sleep(2.5)
	os.system("go run ./treecli -delete 12 -id 545 -token 37")
	print("[-->]  Delete value")
	time.sleep(2.5)

def TestFind():
	'''
	Function to test the find command.
	'''
	os.system("go run ./treecli -find 226 -id 545 -token 37")
	print("[-->]  Found value")
	time.sleep(2.5)
	os.system("go run ./treecli -find 12 -id 545 -token 37")
	print("[-->]  Found value")
	time.sleep(2.5)

def TestTraverse():
	'''
	Function to test the traverse command of our Trie.
	'''
	os.system("go run ./treecli -traverse -id 545 -token 37")
	print("[-->]  Traversed Trie")
	time.sleep(2.5)

def TestDeleteTrie():
	'''
	Function to test the DeleteTrie Command.
	'''
	os.system("go run ./treecli -delete-trie {}")
	sys.stdout.write('yes')
	press_and_release('enter')
	print("[-->]  Deleted Trie")

def TestcreateTrie(size):
	'''
	Function to create a trie.
	@param size will be the size of the trie.
	'''
	os.system("go run ./treecli -create-trie {}".format(size))
	print("[+] Created Trie")
	time.sleep(2.5)

def RunGoTests():
	'''
	Function to run the go test command in order to see the results of the Test.
	'''
	os.system("go test tree/... -cover -timeout 10")
	print("[+] Created Trie")

def main():
	'''
	Programm entry point.
	Will execute all possible commands to test the trie.
	'''
	print("[+] Script started!")
	TestcreateTrie(2)
	TestInsert()
	TestChange()
	TestDelete()
	TestFind()
	TestDeleteTrie()
	TestTraverse()
	RunGoTests()
	print("[+] Script finished!")


if __name__ == "__main__":
	main()