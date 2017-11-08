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
