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
	"encoding/binary"
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

// Attempt - 3
func CalculateTCPHeaderChecksum(src_ip, dst_ip net.IP, data []byte) uint16{
	pseudo_header := make([]byte, 12)

	fmt.Println("Checksum Source IP: ", src_ip.To4())
	fmt.Println("Checksum Destination IP: ", dst_ip.To4())

	copy(pseudo_header[0:4], src_ip.To4()) // IP.To4() returns a slice of bytes
	copy(pseudo_header[4:8], dst_ip.To4())
	pseudo_header[8] = 0 // Extra `0` Padding
	pseudo_header[9] = 6 // Set Protocol to 6

	// Calculate Header + Data length
	// len() returns 32bit int -> Convert to 16
	//tcp_length16 := uint16(len(data))
	binary.BigEndian.PutUint16(pseudo_header[10:], uint16(len(data)))
	// pseudo_header[10] = tcp_length16 >> 8
	// pseudo_header[11] = uint8(byte(tcp_length16) & 0x00ff)

	fmt.Println("Pseaudo Header Bits: ", pseudo_header)

	data = append(data, pseudo_header...)

	checksum := uint32(0)

	fmt.Println("Calculating Checksum: ")
	for i := 0 ; i < len(data) - 3 ; i += 4 {
		byte32 := uint32(data[i]) << 24 + uint32(data[i + 1]) << 16 + uint32(data[i + 2]) << 8 + uint32(data[i + 3])
		fmt.Println("Convert Data to 32 bit: ", byte32)
		
		checksum += byte32
	}

	// If the len of data is not divisible by 4 - Add Padded Bits to data and add to checksum - FIXME: Could Fail as data should ideally be in 16bit. Soln: Create a separaet overflow manager
	if len(data) % 4 != 0 {
		fmt.Println("Checksum Need Padded Data Bits: ")
		extra_bytes := (len(data) % 4)

		byte32 := uint32(0)
		for i := len(data) - extra_bytes ; i < len(data) ; i ++ {
			byte32 = byte32 << 8 // FIXME: Could be incorrect - maybe the shift is wrong
			byte32 += uint32(data[i])
		}

		fmt.Println("Extra Padded Bits Value: ", byte32)
		checksum += byte32
	}

	fmt.Println("Checksum Sum: ", checksum)

	// Fold 32 bit checksum - Add any bit in pos > 16 to the end until num is 16 bits long
	for checksum > 0xffff {
		extra := checksum >> 16
		fmt.Println("Extra Bit: ", extra)
		fmt.Println("Masked: ", checksum & 0xffff)
		checksum = (checksum & 0xffff) + extra
	}

	fmt.Println("Folded Checksum: ", checksum)

	var finalChecksum uint16
	// finalChecksum = ^uint16(checksum)
	finalChecksum = ^uint16(checksum + 1) // Temp dont + 1
	fmt.Println("Final Complement Checksum: ", finalChecksum)

	return finalChecksum
}

// func CalculateTCPHeaderChecksum(data []byte) uint16 {
// 	fmt.Println("Calculating Checksum")
// 	finalChecksum := uint16(0)

// 	// Iterate through the data in 16-bit chunks
// 	for i := 0; i < len(data); i += 2 {
// 		// Take the current byte and the next byte, if available
// 		byte16 := uint16(data[i]) << 8
// 		if i+1 < len(data) {
// 			byte16 |= uint16(data[i+1])
// 		}
// 		fmt.Println("Current Byte: ", byte16)
// 		finalChecksum += byte16
// 	}

// 	// Fold 32-bit result into 16-bit checksum
// 	for finalChecksum > 0xffff {
// 		finalChecksum = (finalChecksum >> 16) + (finalChecksum & 0xffff)
// 	}

// 	// One's complement of the checksum
// 	finalChecksum = ^finalChecksum
// 	fmt.Println("Checksum: ", finalChecksum)
// 	return finalChecksum
// }

// func CalculateTCPHeaderChecksum(data []byte) uint16 {
// 	fmt.Println("Calculating Checksum")
// 	finalChecksum := uint16(0)
// 	i := 1
// 	for i = 1 ; i <= len(data) ; i *= 2{
// 		fmt.Println("Itr Number: ", i)
// 		byte16 := uint16(data[i-1]) << 8 | uint16(data[i])
// 		fmt.Println("Current Byte: ", byte16)
// 		finalChecksum += byte16
// 	}

// 	fmt.Println("Checksum: ",finalChecksum)
// 	finalChecksum = ^finalChecksum
// 	fmt.Println("Complement Checksum: ",finalChecksum)
// 	return finalChecksum
// }

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
		SourcePort: srcPort, // uint16
		DestPort:   dstPort, // uint16
		SeqNumber:  uint32(1001),
		AckNumber:  uint32(0),
		// Set DataOffset to 5, which is the correct value for a TCP header with no options
		DataOffset: 0x50, // This will be 80 in decimal or 0x50 - uint8
		Flags:      0x02, // uint8
		WindowSize: uint16(5840),
		Checksum: uint16(0),
		UrgentPtr:  0, // uint16
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

	checksum := CalculateTCPHeaderChecksum(net.ParseIP("192.168.29.182"), pip, tcpByte)
	tcpByte[16] = byte(checksum >> 8)
	tcpByte[17] = byte(checksum & 0x00ff)
	fmt.Println("Checksum: ", checksum)
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


	//Send the crafted TCP packet
	for i := range(5) {
		fmt.Println("Attempt Number: ", i )
		_, err = conn.Write(tcpByte)
		if err != nil {
			log.Fatal("Error in Sending: ", err)
		}
	}

	// Close TCP Connection - Hopefuly hole punching is done
	conn.Close()
	localIp := "0.0.0.0:" + strconv.Itoa(int(srcPort))
	dstIp := ip + ":" + strconv.Itoa(int(dstPort))

	var mode string
	fmt.Println("Syn Packet Sent")
	fmt.Print("Enter TCP Mode [(d)Dial / (r)Recieve]: ")
	fmt.Scan(&mode)
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
