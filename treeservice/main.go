package main

import (
	"fmt"

	"github.com/ob-vss-ws19/blatt-3-pwn/tree"

	"github.com/AsynkronIT/protoactor-go/actor"
)

func main() {
	fmt.Println("Hello Tree-Service!")
	/*
		l := tree.NewLeaf()
		fmt.Printf("Contains 5 %t \n", l.Contains(5))
		for i := 0; i < 10; i++ {
			l.Insert(i, "Hallo")
		}
		fmt.Printf("Contains 5 %t \n", l.Contains(5))
		fmt.Printf("Can insert %t \n", l.Insert(7, "kakak"))
		fmt.Printf("Vorher %d \n", l.Size())
		l.Erase(5)
		fmt.Printf("Contains 5 %t \n", l.Contains(5))
		fmt.Printf("Nachher %d \n", l.Size()) */
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.Nodeactor{}
	})
	pid := context.Spawn(props)
	context.Send(pid, &tree.Insertmessage{
		Element: tree.Pair{
			Key:   12,
			Value: "Hallo vom Planeten Paulanius",
		},
	})
	var myname string
	fmt.Scanf("%s", &myname)
}
