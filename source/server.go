package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prologic/bitcask"
)

// EnrollFlag enrolls the flag
func EnrollFlag(id int, flag string) {
	bid := []byte(strconv.Itoa(id))
	bflag := []byte(flag)
	db, err := bitcask.Open("databases/flag")
	defer db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		return
	}
	db.Put(bid, bflag)
	fmt.Printf("[+] Enroll { %v : %v }\n", id, flag)
}

// GetFlag gets the flag
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

// DeleteFlag deletes the flag
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

// CheckFlag checks whether the submitted flag is correct or not
func CheckFlag(id string, submission string) (bool, error) {
	bid := []byte(id)
	db, _ := bitcask.Open("databases/flag")
	defer db.Close()
	bflag, err := db.Get(bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] %v\n", err)
		return false, err
	}

	return string(bflag) == submission, nil

}

func main() {
	fmt.Print("[+] Starting server...\n")
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	EnrollFlag(0, "FLAG")
	// flag := GetFlag(0)
	// fmt.Println(flag)
	// flag2 := GetFlag(5)
	// fmt.Println(flag2)
	// DeleteFlag(0)
	r.POST("/submit", func(c *gin.Context) {
		c.Request.ParseForm()
		problemid := c.Request.Form["id"]
		submittedflag := c.Request.Form["flag"]
		correct, err := CheckFlag(problemid[0], submittedflag[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[-] %v\n", err)

		}
		c.HTML(200, "index.html", gin.H{"flag": submittedflag, "correct": correct})

	})
	r.Run()
}
