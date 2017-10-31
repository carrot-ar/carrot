// The MIT License (MIT)

// Copyright (c) 2016 Jupp Müller

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// https://github.com/jupp0r/go-priority-queue/blob/master/priorty_queue_test.go

//The following code has been modified from the original to fit our implementation.

package carrot

import (
	"sort"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue()
	elements := []float64{5, 3, 7, 8, 6, 1, 9}
	for _, e := range elements {
		pq.Insert(e, e)
	}

	sort.Float64s(elements)
	for _, e := range elements {
		item, err := pq.Pop()
		if err != nil {
			t.Fatalf(err.Error())
		}

		i := item.(float64)
		if e != i {
			t.Fatalf("expected %v, got %v", e, i)
		}
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	pq := NewPriorityQueue()
	pq.Insert("foo", 3)
	pq.Insert("bar", 4)
	pq.UpdatePriority("bar", 2)

	item, err := pq.Pop()
	if err != nil {
		t.Fatal(err.Error())
	}

	if item.(string) != "bar" {
		t.Fatal("priority update failed")
	}
}

func TestPriorityQueueLen(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Len() != 0 {
		t.Fatal("empty queue should have length of 0")
	}

	pq.Insert("foo", 1)
	pq.Insert("bar", 1)
	if pq.Len() != 2 {
		t.Fatal("queue should have length of 2 after 2 inserts")
	}
}

func TestDoubleAddition(t *testing.T) {
	pq := NewPriorityQueue()
	pq.Insert("foo", 2)
	pq.Insert("bar", 3)
	pq.Insert("bar", 1)

	if pq.Len() != 2 {
		t.Fatal("queue should ignore inserting the same element twice")
	}

	item, _ := pq.Pop()
	if item.(string) != "foo" {
		t.Fatal("queue should ignore duplicate insert, not update existing item")
	}
}

func TestPopEmptyQueue(t *testing.T) {
	pq := NewPriorityQueue()
	_, err := pq.Pop()
	if err == nil {
		t.Fatal("should produce error when performing pop on empty queue")
	}
}

func TestUpdateNonExistingItem(t *testing.T) {
	pq := NewPriorityQueue()

	pq.Insert("foo", 4)
	pq.UpdatePriority("bar", 5)

	if pq.Len() != 1 {
		t.Fatal("update should not add items")
	}

	item, _ := pq.Pop()
	if item.(string) != "foo" {
		t.Fatalf("update should not overwrite item, expected \"foo\", got \"%v\"", item.(string))
	}
}
