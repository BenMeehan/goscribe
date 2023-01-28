# GoScribe
GoScribe is a minimal pub-sub broker written in pure golang.

## Usage 
1. "sub _topic_" - subscribe to a topic
2. "unsub _topic_" - unsubscribe from a topic
3. "pub _topic_" - publish to a topic
4. "ls" - list all topics
5. "quit" - disconnect from the broker

## How to run
1. clone this repo 
    `git clone `
2. cd into the cloned directory
3. run `go build`
4. run the resultant binary `./goscribe [-h HOST] [-p PORT]` 

Note : HOST and PORT are optional. Defaults to 0.0.0.0 and 8090 

## Demo using telnet
![demo.gif](https://s9.gifyu.com/images/Screen-Recording-2023-01-28-at-3.15.12-PM.gif)

## Gotcha's
GoScribe is a very simple broker. Think of it as a chat room server but for pub-sub. It does not provide any message persistance or queueing for now. What is not recieved by the subscribers is lost for ever!

## TODO
- [x] client library
- [x] queueing messages inside broker
- [x] message persistence
- [x] health check if subscribers

***Thanks to [Pilutau](https://www.youtube.com/watch?v=Sphme0BqJiY) for inspiring this project.***