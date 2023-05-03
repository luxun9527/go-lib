package standard

import (
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	os.Truncate("S:\\go-lib\\standard\\a.txt", 0)
	os.WriteFile("S:\\go-lib\\standard\\a.txt", []byte("sss"), os.ModePerm)
}
func TestTruncate(t *testing.T) {
	file, err := os.OpenFile("S:\\go-lib\\standard\\a.txt", os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString("abcd")

}
func TestOpen(t *testing.T) {
	file, err := os.OpenFile("S:\\go-lib\\standard\\a.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	file.WriteString("------")

}
