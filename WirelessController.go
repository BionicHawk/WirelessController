package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

const (
	EOM string = "<|EOM|>"
	MIN_W int = 0
	MAX_W int = 1366
	MIN_H int = 0
	MAX_H int = 768
)

var current_x int = 0
var current_y int = 0

func main() {
	const Type = "tcp"
	var Host string= "localhost"
	var Port uint16 = 8080

	if len(os.Args) > 2 {
		Host = os.Args[1]
		pre_port, err := strconv.ParseUint(os.Args[2], 10, 16)
		if err != nil {
			fmt.Println("Not valid port")
			os.Exit(-1)
		}
		Port = uint16(pre_port)
	}

	port_string := strconv.Itoa(int(Port))

	listen, err := net.Listen(Type, Host + ":" + port_string)

	if err != nil {
		println("Couldn't initialize the tcp server!")
		os.Exit(-1)
	}

	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Printf("Couldn't Accept the connection: %s\n", err.Error())
			continue
		}

		HandleConnection(conn)
		
	}

}

func HandleConnection(Connection net.Conn) {
	defer Connection.Close()

	for {
		buffer := make([]byte, 1024)
		_, err := Connection.Read(buffer)
		
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err.Error())
		}
		
		serv_inp := string(buffer)

		HandleServerInput(&serv_inp)
	}

}

func HandleServerInput(input *string) {
	const N_COMPS int = 3

	if strings.Contains(*input, EOM) {
		coords := strings.Split(*input, EOM)
		coords_size := len(coords)

		for i:=0; i < coords_size; i+=1 {
			coord := coords[i]
			components := strings.Split(coord, ",")
			components_size := len(components)
			
			if components_size == N_COMPS {
				fmt.Printf("\r%d, %d        ", current_x, current_y)
				HandleCoords(components)
			}
		}
	}
}

func HandleCoords(coords []string) {
	amp := 90
	pre_x, err_x := strconv.ParseFloat(coords[2], 64)	
	pre_y, err_y := strconv.ParseFloat(coords[0], 64)	
	_, err_z := strconv.ParseFloat(coords[1], 64)
	
	if err_x != nil || err_y != nil || err_z != nil {
		return
	}

	current_x -= int(pre_x) * amp
	current_y -= int(pre_y) * amp

	robotgo.Move(current_x, current_y)

	if current_x > MAX_W {
		current_x = MAX_W
	}
	
	if current_x < MIN_W {
		current_x = MIN_W
	}

	if current_y > MAX_H {
		current_y = MAX_H
	}

	if current_y < MIN_H {
		current_y = MIN_H
	}

}
