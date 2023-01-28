package main

// Topic
// Name :  Name of the topic
// Subscribers : Map of clients subscribed to the topic. Using map as it is faster to delete entries.
type Topic struct {
	name        string
	subscribers map[*Emptor]bool
}
