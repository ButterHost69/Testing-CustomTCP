// package main

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"
// 	"time"
// )

// type TCP struct {
// 	SourcePort uint16
// 	DestPort   uint16
// 	SeqNumber  uint32
// 	AckNumber  uint32
// 	DataOffset uint16
// 	Flags      uint16
// 	WindowSize uint16
// 	Checksum   uint16
// 	UrgentPtr  uint16
// }

// func main() {
// 	var ip string
// 	var srcPort uint16
// 	var dstPort uint16

// 	fmt.Println("Hello World...")
// 	fmt.Print("Enter IP: ")
// 	fmt.Scan(&ip)

// 	fmt.Println("Entered IP: ", ip)
// 	fmt.Print("Enter Source Port: ")
// 	fmt.Scan(&srcPort)

// 	fmt.Print("Enter Destination Port: ")
// 	fmt.Scan(&dstPort)

// 	ip = strings.TrimSpace(ip)
// 	pip := net.ParseIP(ip)
// 	if pip == nil {
// 		log.Fatal("Could Not Parse IP: ")
// 	}

// 	remoteAddr := &net.IPAddr{
// 		IP: pip,
// 	}

// 	// Craft TCP Packet
// 	tcpPacket := TCP{
// 		SourcePort: srcPort,
// 		DestPort: dstPort,
// 		SeqNumber: uint32(10001),
// 		AckNumber: uint32(0),
// 		DataOffset: 5 << 4,
// 		Flags: uint16(0x02),
// 		WindowSize: uint16(5840),
// 		UrgentPtr: 0,
// 	}

// 	tcpByte := []byte {
// 		byte(tcpPacket.SourcePort >> 8) , byte(tcpPacket.SourcePort & 0x00ff),
// 		byte(tcpPacket.DestPort >> 8), byte(tcpPacket.DestPort & 0xff),
// 		byte(tcpPacket.SeqNumber >> 24), byte(tcpPacket.SeqNumber >> 16), byte(tcpPacket.SeqNumber >> 8), byte(tcpPacket.SeqNumber),
// 		byte(tcpPacket.AckNumber >> 24), byte(tcpPacket.AckNumber >> 16), byte(tcpPacket.AckNumber >> 8), byte(tcpPacket.AckNumber),
// 		byte(tcpPacket.DataOffset >> 8), byte(tcpPacket.DataOffset & 0xff),
// 		byte(tcpPacket.Flags >> 8), byte(tcpPacket.Flags & 0xff),
// 		byte(tcpPacket.WindowSize >> 8), byte(tcpPacket.WindowSize & 0xff),
// 		byte(tcpPacket.Checksum >> 8), byte(tcpPacket.Checksum & 0xff),
// 		byte(tcpPacket.UrgentPtr >> 8), byte(tcpPacket.UrgentPtr & 0xff),
// 	}

// 	fmt.Println("TCP Packet: ", tcpPacket)
// 	fmt.Println("TCP byte: ", tcpByte)

// 	conn, err := net.DialIP("ip4:tcp", nil, remoteAddr)
// 	if err != nil {
// 		log.Fatal("Error in IP: ", err)
// 	}

// 	timeout := 3 * time.Second
// 	err = conn.SetDeadline(time.Now().Add(timeout))
// 	if err != nil {
// 		log.Fatal("Timeout Error: ", err)
// 	}

// 	defer conn.Close()

// 	// Send ICMP Echo Request
// 	// request := []byte{8, 0, 0, 0, 0, 0, 0, 0} // ICMP Echo Request
// 	// _, err = conn.Write(request)
// 	// if err != nil {
// 	// 	log.Fatal("Send Error: ", err)
// 	// }

// 	// // Read Response
// 	// response := make([]byte, 1024)
// 	// _, err = conn.Read(response)
// 	// if err != nil {
// 	// 	log.Fatal("Receive Error: ", err)
// 	// }

// 	// fmt.Println("IP is reachable!")

// 	//Send TCP Packet
// 	_, err = conn.Write(tcpByte)
// 	if err != nil {
// 		log.Fatal("Ero in Sending: ", err)
// 	}

// }

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
    stun "github.com/ccding/go-stun/stun"
)

type TCP struct {
	SourcePort uint16
	DestPort   uint16
	SeqNumber  uint32
	AckNumber  uint32
	DataOffset uint8 // The DataOffset is still in 16-bits but should be calculated as 5 for a standard 20-byte header
	Flags      uint8
	WindowSize uint16
	Checksum   uint16
	UrgentPtr  uint16
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		// Read incoming data
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed:", conn.RemoteAddr())
			return
		}

		// Echo the data back
		fmt.Printf("Received: %s", string(buf[:n]))
		conn.Write(buf[:n])
	}
}

func main() {
	var ip string
	var srcPort uint16
	var dstPort uint16

	client := stun.NewClient()
    client.SetServerAddr("stun.l.google.com:19302")
    client.SetLocalPort(8090)
    client.SetLocalIP("0.0.0.0")

    nat, host, err := client.Discover()
    fmt.Println(nat.String())
    fmt.Println(host)
    fmt.Println(err)

	fmt.Println("Hello World...")
	fmt.Print("Enter IP: ")
	fmt.Scan(&ip)

	fmt.Println("Entered IP: ", ip)
	fmt.Print("Enter Source Port: ")
	fmt.Scan(&srcPort)

	fmt.Print("Enter Destination Port: ")
	fmt.Scan(&dstPort)

	ip = strings.TrimSpace(ip)
	pip := net.ParseIP(ip)
	if pip == nil {
		log.Fatal("Could Not Parse IP: ")
	}

	remoteAddr := &net.IPAddr{
		IP: pip,
	}

	// Craft TCP Packet - SYN Packet
	// s := fmt.Sprintf("%b", 80)
	tcpPacket := TCP{
		SourcePort: srcPort,
		DestPort:   dstPort,
		SeqNumber:  uint32(1001),
		AckNumber:  uint32(0),
		// Set DataOffset to 5, which is the correct value for a TCP header with no options
		DataOffset: 0x50, // This will be 80 in decimal or 0x50
		Flags:      0x02,
		WindowSize: uint16(5840),
		UrgentPtr:  0,
	}

	tcpByte := []byte{
		// Split fields to match the TCP header structure
		byte(tcpPacket.SourcePort >> 8), byte(tcpPacket.SourcePort & 0x00ff),
		byte(tcpPacket.DestPort >> 8), byte(tcpPacket.DestPort & 0xff),
		byte(tcpPacket.SeqNumber >> 24), byte(tcpPacket.SeqNumber >> 16), byte(tcpPacket.SeqNumber >> 8), byte(tcpPacket.SeqNumber),
		byte(tcpPacket.AckNumber >> 24), byte(tcpPacket.AckNumber >> 16), byte(tcpPacket.AckNumber >> 8), byte(tcpPacket.AckNumber),
		byte(tcpPacket.DataOffset),
		//byte(tcpPacket.DataOffset & 0xff),  // DataOffset is 5 (20 bytes of header)
		byte(tcpPacket.Flags),
		//byte(tcpPacket.Flags & 0xff),
		byte(tcpPacket.WindowSize >> 8), byte(tcpPacket.WindowSize & 0xff),
		byte(tcpPacket.Checksum >> 8), byte(tcpPacket.Checksum & 0xff),
		byte(tcpPacket.UrgentPtr >> 8), byte(tcpPacket.UrgentPtr & 0xff),
	}

	fmt.Println("TCP Packet: ", tcpPacket)
	fmt.Println("Offset: ", tcpPacket.DataOffset)
	fmt.Println("TCP byte: ", tcpByte)

	// Send packet using net.DialIP
	conn, err := net.DialIP("ip4:tcp", nil, remoteAddr)
	if err != nil {
		log.Fatal("Error in IP: ", err)
	}

	timeout := 3 * time.Second
	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		log.Fatal("Timeout Error: ", err)
	}


	// Send the crafted TCP packet
	_, err = conn.Write(tcpByte)
	if err != nil {
		log.Fatal("Error in Sending: ", err)
	}

	// Close TCP Connection - Hopefuly hole punching is done
	conn.Close()
	localIp := "0.0.0.0:" + strconv.Itoa(int(srcPort))
	dstIp := ip + ":" + strconv.Itoa(int(dstPort))

	var mode string
	fmt.Println("Syn Packet Sent")
	fmt.Print("Enter TCP Mode [(d)Dial / (r)Recieve]: ")
	if mode == "d" {
		// Define source address

		localAddr, err := net.ResolveTCPAddr("tcp", localIp) // Change IP & port as needed
		if err != nil {
			fmt.Println("Error resolving local address:", err)
			return
		}

		// Define destination address

		remoteAddr, err := net.ResolveTCPAddr("tcp", dstIp)
		if err != nil {
			fmt.Println("Error resolving remote address:", err)
			return
		}

		// Dial with specific source address
		conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
		if err != nil {
			fmt.Println("Error dialing:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Connected to", conn.RemoteAddr(), "from", conn.LocalAddr())
	} else {
		laddr, err := net.ResolveTCPAddr("tcp", localIp) // Listen on all interfaces, port 8080
		if err != nil {
			fmt.Println("Error resolving address:", err)
			os.Exit(1)
		}

		// Start listening on the specified TCP address
		listener, err := net.ListenTCP("tcp", laddr)
		if err != nil {
			fmt.Println("Error starting TCP listener:", err)
			os.Exit(1)
		}
		defer listener.Close()

		fmt.Println("Listening on", laddr)

		for {
			// Accept incoming connections
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			fmt.Println("New connection from", conn.RemoteAddr())

			// Handle connection in a separate goroutine
			go handleConnection(conn)
		}
	}
}
