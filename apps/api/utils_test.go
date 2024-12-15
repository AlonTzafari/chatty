package main

import "testing"

func TestRemove(t *testing.T) {
	slice := []string{"one", "two", "three"}
	slice = remove(slice, 1)
	if slice[1] == "two" {
		t.Fatal("remove failed ", slice)
	}

}

func TestFindIndex(t *testing.T) {
	slice := []string{"one", "two", "three"}
	i := findIndex(slice, "two")
	if i != 1 {
		t.Fatal("i", i)
	}
}
