package main

type room struct{
	name string
	memebers map[net.Addr]*client

}
