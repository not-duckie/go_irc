package main

import ("fmt"
		"net"
		"log"
		"strings"
		"errors"
	)


type server struct{
	rooms map[string]*room
	commands chan command
}


func newServer()(*server){
	return &server{
		rooms:  make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server)run(){
	for cmd:= range s.commands{
		switch cmd.id{
		case CMD_NICK:
			s.nick(cmd.client,cmd.args)
		case CMD_JOIN:
			s.join(cmd.client,cmd.args)
		case CMD_ROOMS:
			s.listrooms(cmd.client,cmd.args)
		case CMD_MSG:
			s.msg(cmd.client,cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client,cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn){
	log.Printf("new client has connected %s\n",conn.RemoteAddr().String())
	c := &client{
		conn: conn,
		nick: "anonymous",
		commands: s.commands,
	}
	c.readInput()

}

func (s *server)nick(c *client,args []string){
	if len(args) < 2{
		c.err(errors.New("error in command usage"))
		return
	}
	c.nick = args[1]
	c.msg(fmt.Sprintf("you nick is now %s",c.nick))
}

func (s *server)join(c *client,args []string){
	if len(args) < 2{
		c.err(errors.New("error in command usage"))
		return
	}
	roomName := args[1]
	r,ok := s.rooms[roomName]
	if !ok{
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	if c.room != nil{
		s.quitCurrentRoom(c)
	}
	c.room = r	
	r.broadcast(c,fmt.Sprintf("%s has join the chat",c.nick))
	c.msg(fmt.Sprintf("you have entered the chat %s",r.name))
}
func (s *server)listrooms(c *client,args []string){
	var rooms []string
	for name:= range s.rooms{
		rooms = append(rooms,name)
	}
	c.msg(fmt.Sprintf("available rooms are: %s",strings.Join(rooms,", ")))
}

func (s *server)msg(c *client,args []string){
	if len(args) < 2{
		c.err(errors.New("error in command usage"))
		return
	}
	if c.room == nil{
		c.err(errors.New("you must join a room first"))
		return
	}
	c.room.broadcast(c,c.nick+": "+strings.Join(args[1:len(args)]," "))
}
func (s *server)quit(c *client,args []string){
	log.Printf("client has disconnected: %s",c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("bie bye")
	c.conn.Close()
}

func (s *server)quitCurrentRoom(c *client){
	if c.room != nil{
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c,fmt.Sprintf("%s has left the chat",c.nick))
	}
}