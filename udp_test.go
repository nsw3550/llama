package udprobe

import (
	"bytes"
	"net"
	"testing"

	pb "github.com/nsw3550/udprobe/proto"
	"google.golang.org/protobuf/proto"
)

func TestProbeProtobuf(t *testing.T) {
	// Test with good data first
	signature := []byte("abcdefghij")
	data := &pb.Probe{
		Signature: signature,
		Tos:       46,
		Sent:      123456789,
	}
	marshaled, err := proto.Marshal(data)
	if err != nil {
		t.Fatal("Failed to marshal probe data:", err)
	}

	unmarshaled := &pb.Probe{}
	err = proto.Unmarshal(marshaled, unmarshaled)
	if err != nil {
		t.Fatal("Failed to unmarshal probe data:", err)
	}

	// Compare the actual structs
	if !bytes.Equal(unmarshaled.Signature, data.Signature) ||
		unmarshaled.Tos != data.Tos ||
		unmarshaled.Sent != data.Sent {
		t.Error("Data unmarshaled, but lost in translation")
	}

	// Now verify that bad data doesn't work
	badData := []byte{1, 2, 3, 4, 5}
	err = proto.Unmarshal(badData, &pb.Probe{})
	if err == nil {
		t.Error("No error returned for bad data (though unmarshal might succeed with empty/corrupt proto)")
	}
}

func TestSetTos(t *testing.T) {
	// Resolve a local addr
	myAddr, _ := net.ResolveUDPAddr("udp", ":0")
	// Create a connection
	conn, _ := net.ListenUDP("udp", myAddr)
	// Set the ToS value
	tosVal := 240
	newTos := byte(tosVal)
	SetTos(conn, newTos)
	// Verify the ToS value
	val := GetTos(conn)
	if val != newTos {
		t.Error("New ToS value not set correctly. Set", tosVal, "and got",
			val, "instead.")
	}
}
