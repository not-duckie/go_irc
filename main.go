package main


func main(){
	s := newServer()
	go s.run()
	listener,err:= net.Listen("tcp",":8080")
	if err!=nil{
		log.Fatalln("failed to bind to port",err.Error())
	}
	defer listener.Close()
	log.Printf("started the server on port :8080")

	for {
		conn,err := listener.Accept()
		if err!=nil{
			log.Println("Failed to accept connection")
			conitnue
		}
		go s.newClient(conn)

	}

}
