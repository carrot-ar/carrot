package carrot

import "testing"

func TestNewClientList(t *testing.T) {
	_, err := NewClientList()

	if err != nil {
		t.Fatalf("client list failed to create, %v", err)
	}
}

func TestClientsInsert(t *testing.T) {
	clientList, _ := NewClientList()

	client := &Client{}
	err := clientList.Insert(client)

	if err != nil {
		t.Fatalf("failed to insert into client list, %v", err)
	}
}

func TestClientsRelease(t *testing.T) {
	clientList, _ := NewClientList()

	client := &Client{}
	clientList.Insert(client)

	clientList.Release(0)

	for _, client := range clientList.clients {
		if client.Valid() {
			t.Fatalf("no clients should be valid!")
		}
	}
}
