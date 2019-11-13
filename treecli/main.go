package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-pwn/messages"
	"github.com/ob-vss-ws19/blatt-3-pwn/tree"
)

/*
inserts information in a message which is necessary to delete an entry in the tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendDelete(id int, token string, key int) (bool, error) {
	message := &messages.DeleteRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(key),
	}
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Delete %d", key)
}

/*
inserts information in a message which is necessary to change an entry in the tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendChange(id int, token string, pair tree.Pair) (bool, error) {
	message := &messages.ChangeRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(pair.Key),
		Value: pair.Value,
	}
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Change %d %s", pair.Key, pair.Value)
}

/*
inserts information in a message which is necessary to find an entry in the tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendFind(id int, token string, key int) (bool, error) {
	message := &messages.FindRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(key),
	}
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Find %d", key)
}

/*
inserts information in a message which is necessary to insert an entry in the tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendInsert(id int, token string, pair tree.Pair) (bool, error) {
	message := &messages.InsertRequest{
		Token: token,
		Id:    int32(id),
		Key:   int32(pair.Key),
		Value: pair.Value,
	}
	fmt.Println("Sending CLI")
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Insert %d %s", pair.Key, pair.Value)
}

/*
inserts information in a message which is necessary to creat a tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendCreateTrie(size int) (bool, error) {
	message := &messages.CreateTreeRequest{
		Size_: int32(size),
	}
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Create Tree")
}

/*
inserts information in a message which is necessary to delete a tree
calls the function that sends the message to the service
returns true if that was successful
*/
func sendDeleteTrie(id int, token string) (bool, error) {
	message := &messages.DeleteTreeRequest{
		Token: token,
		Id:    int32(id),
	}
	se, er := remotesend(message)
	if er == nil && se {
		return true, nil
	}
	return false, fmt.Errorf("Cannot Change %d", id)
}

/*
starts the cli. The cli uses command line parameters
 #create a tree
 #delete a tree
 #search, insert and delete entries
*/
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

	deleteTrie := flag.Int("delete-trie", -1, "If this flag is set your trie will be deleted ~ DANGEROUS")

	flag.Parse()
	if *insertBool && *key != -1 && *value != "" {
		fmt.Println("Start Insert CLI")
		time.Sleep(5 * time.Second)
		_, _ = sendInsert(*ID, *token, tree.Pair{
			Key:   *key,
			Value: *value,
		})
		time.Sleep(5 * time.Second)
		fmt.Println("Stop Insert CLI")
	} else if *changeBool && *key != -1 && *value != "" {
		_, _ = sendChange(*ID, *token, tree.Pair{
			Key:   *key,
			Value: *value,
		})
	} else if *delete != -1 {
		_, _ = sendDelete(*ID, *token, *key)
	} else if *deleteTrie != -1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are sure you wanna delete the trie (yes/no)")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		g := strings.HasPrefix(text, "yes")
		fmt.Println(g)
		if g {
			fmt.Println("Trie will now be deleted")
			_, _ = sendDeleteTrie(*ID, *token)
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("Trie will not be deleted")
		}
	} else if *createTree != -1 {
		//fmt.Println("Before")
		time.Sleep(5 * time.Second)
		_, _ = sendCreateTrie(*createTree)
		time.Sleep(5 * time.Second)
		//fmt.Println("hahahahah")
	} else {
		fmt.Println("Please make sure your given arguments fit the required pattern for more info enter .. -h")
	}
}

/*
ClientRemoteActor represents a RemoteActore in the client
*/
type ClientRemoteActor struct {
	count int
	wg    *sync.WaitGroup
}

/*
Receive will receive the response from the Service. Either if the action connected to the message could be executed or not
*/
func (state *ClientRemoteActor) Receive(context actor.Context) {
	//fmt.Println(context.Message())
	switch msg := context.Message().(type) {
	case *messages.Response:
		state.count++
		fmt.Println(state.count, msg)
	case *actor.Stopped:
		state.wg.Done()
	default:
		fmt.Println("Test")
	}
}

/*
starts an actorsystem that can be reached remotely and can send messages to the service
*/
func remotesend(mess interface{}) (bool, error) {

	remote.Start("localhost:8090")

	var wg sync.WaitGroup

	context := actor.EmptyRootContext

	props := actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &ClientRemoteActor{0, &wg}
	})

	pid := context.Spawn(props)
	fmt.Println(pid)

	fmt.Println("Sleeping 5 seconds...")
	time.Sleep(5 * time.Second)
	fmt.Println("Awake...")
	fmt.Printf("Trying to connect")

	spawnResponse, err := remote.SpawnNamed("localhost:8091", "remote", "hello", 5*time.Second)
	if err == nil {
		context.RequestWithCustomSender(spawnResponse.Pid, mess, pid)
		return true, nil
	}
	wg.Wait()
	return false, errors.New("Cannot send to remote Controller")

}
