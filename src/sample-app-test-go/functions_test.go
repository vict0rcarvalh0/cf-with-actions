package main

import "testing"

func TestAdd(t *testing.T) {
    result := Add(3, 5)
    expected := 8
    if result != expected {
        t.Errorf("Add(3, 5) = %d; want %d", result, expected)
    }
}

func TestIsEven(t *testing.T) {
    result := IsEven(4)
    if !result {
        t.Errorf("IsEven(4) = %v; want %v", result, true)
    }

    result = IsEven(5)
    if result {
        t.Errorf("IsEven(5) = %v; want %v", result, false)
    }
}

func TestReverse(t *testing.T) {
    result := Reverse("golang")
    expected := "gnalog"
    if result != expected {
        t.Errorf("Reverse('golang') = %s; want %s", result, expected)
    }
}
