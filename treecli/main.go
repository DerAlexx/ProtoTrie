package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

func sendDelete(id int, token string, key int) (bool, error) {
	message := &messages.DeleteRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(key),
	}
	se, er := remotesend(message)
	if er != nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Delete %d", key)
}

func sendChange(id int, token string, pair tree.Pair) (bool, error) {
	message := &messages.ChangeRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(pair.Key),
		Value: pair.Value,
	}
	se, er := remotesend(message)
	if er != nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Change %d %s", pair.Key, pair.Value)
}

func sendInsert(id int, token string, pair tree.Pair) (bool, error) {
	message := &messages.InsertRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(pair.Key),
		Value: pair.Value,
	}
	se, er := remotesend(message)
	if er != nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Insert %d %s", pair.Key, pair.Value)
}

func sendCreateTrie(size int) (bool, error) {
	message := &messages.CreateTreeRequest{
		Size_: int32(size),
	}
	se, er := remotesend(message)
	if er != nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Create Tree")
}

func sendDeleteTrie(id int, token string) (bool, error) {
	message := &messages.DeleteTreeRequest{
		Token: token,
		Id:    int32(id),
	}
	se, er := remotesend(message)
	if er != nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Change %d", id)
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
	if *insertBool && *key != -1 && *value != "" {
		_, _ = sendInsert(*ID, *token, tree.Pair{
			Key:   *key,
			Value: *value,
		})
	} else if *changeBool && *key != -1 && *value != "" {
		_, _ = sendChange(*ID, *token, tree.Pair{
			Key:   *key,
			Value: *value,
		})
	} else if *delete != -1 {
		_, _ = sendDelete(*ID, *token, *key)
	} else if *deleteTrie {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are sure you wanna delete the trie (y/n)")
		text, _ := reader.ReadString('\n')
		if text == "y" {
			fmt.Println("Yes")
		}
		fmt.Println("Trie will not be deleted")
	} else if *createTree != -1 {
		_, _ = sendCreateTrie(*createTree)
	} else {
		fmt.Println("Please make sure your given arguments fit the required pattern for more info enter .. -h")
	}
}

/*
ClientRemoteActor will
*/
type ClientRemoteActor struct {
	count int
}

/*
Receive will
*/
func (state *ClientRemoteActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		fmt.Println(state.count, msg)
	}
default:
	fmt.Println("Test")
}

func remotesend(mess interface{}) (bool, error) {

	remote.Start("localhost:8090")

	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor { return &ClientRemoteActor{} })
	pid := context.Spawn(props)
	fmt.Println(pid)

	spawnResponse, err := remote.SpawnNamed("localhost:8091", "remote", "hello", time.Second)
	if err != nil {
		context.RequestWithCustomSender(spawnResponse.Pid, mess, pid)
		return true, nil
	}

	return false, errors.New("Cannot send to remote Controller")

}
