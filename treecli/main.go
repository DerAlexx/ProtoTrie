package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

func sendDelete(id int, token string, key int) (bool, error) {

	return false, nil
}

func sendChange(id int, token string, pair tree.Pair) (bool, error) {

	return false, nil
}

func sendInsert(id int, token string, pair tree.Pair) (bool, error) {

	return false, nil
}

func sendCreateTrie(size int) (bool, error) {

	return false, nil
}

func sendDeleteTrie(id int, token string) (bool, error) {

	return false, nil
}

func main() {
	fmt.Println("Hello Tree-CLI!")
	createTree := flag.Int("create-trie", -1, "Create a trie and enter the size of the leafs.")
	ID := flag.Int("id", 0, "Enter the ID of your trie")
	token := flag.String("token", "", "Enter a token fitting the ID of the trie in order to gain access to your trie.")

	insertBool := flag.Bool("insert", false, "Set this flag with the Flag key and value in order to insert a pair into your trie")
	changeBool := flag.Bool("change", false, "Set this flag with the Flag key and value in order to change a pair into your trie")
	delete := flag.Int("delete", -1, "Enter a key to delete, the function will remove an entire entry from the trie decided by the key.")

	key := flag.Int("key", -1, "Set this flag in order to pass a key")
	value := flag.String("value", "", "Set this flag in order to pass a value")

	deleteTrie := flag.Bool("delete-trie", false, "If this flag is set your trie will be deleted ~ DANGEROUS")

	flag.Parse()
	if *insertBool != false && *key != -1 && *value != "" {
		sendInsert(*ID, *token, tree.Pair{
			Key:   int(*key),
			Value: *value,
		})
	} else if *changeBool != false && *key != -1 && *value != "" {
		sendChange(*ID, *token, tree.Pair{
			Key:   int(*key),
			Value: *value,
		})
	} else if *delete != -1 {
		sendDelete(*ID, *token, *key)
	} else if *deleteTrie != false {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are sure you wanna delete the trie (y/n)")
		text, _ := reader.ReadString('\n')
		if text == "y" {
			fmt.Println("Yes")
		}
		fmt.Println("Trie will not be deleted")
	} else if *createTree != -1 {
		sendCreateTrie(*createTree)
	} else {
		fmt.Println("Please make sure your given arguments fit the required pattern for more info enter .. -h")
	}
}
