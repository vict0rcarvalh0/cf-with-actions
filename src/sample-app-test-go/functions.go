package main

func Add(a int, b int) int {
    return a + b
}

func IsEven(num int) bool {
    return num%2 == 0
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}
