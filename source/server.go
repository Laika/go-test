package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prologic/bitcask"
)

// RegisterFlag enrolls a flag
func RegisterFlag(id int, flag string) error {
	bid := []byte(strconv.Itoa(id))
	bflag := []byte(flag)
	db, err := bitcask.Open("databases/flag")
	defer db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		return fmt.Errorf("Another flag has been already registered")
	}
	db.Put(bid, bflag)
	fmt.Printf("[+] Enroll { %v : %v }\n", id, flag)
	return nil
}

// GetFlag gets a flag
func GetFlag(id int) string {
	bid := []byte(strconv.Itoa(id))
	db, _ := bitcask.Open("databases/flag")
	defer db.Close()
	flag := ""
	bflag, err := db.Get(bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		return ""
	}
	flag = string(bflag)
	return flag
}

// DeleteFlag deletes a flag
func DeleteFlag(id int) {
	bid := []byte(strconv.Itoa(id))
	db, _ := bitcask.Open("databases/flag")
	defer db.Close()
	flag, err := db.Get(bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		return
	}
	db.Delete(bid)
	fmt.Printf("[+] Delete { %v : %v }\n", id, string(flag))
}

func main() {
	fmt.Print("[+] Starting server...\n")
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	data := "Hello Go/Gin"
	// RegisterFlag(0, "FLAG")
	// flag := GetFlag(0)
	// fmt.Println(flag)
	// flag2 := GetFlag(5)
	// fmt.Println(flag2)
	// DeleteFlag(0)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"data": data})
	})
	r.Run()
}
