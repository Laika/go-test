package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prologic/bitcask"
)

// EnrollFlag enrolls the flag
func EnrollFlag(id string, flag string) error {
	bid := []byte(id)
	bflag := []byte(flag)
	db, err := bitcask.Open("databases/flag")
	defer db.Close()
	if err != nil {
		return err
	}
	if !db.Has(bid) {
		db.Put(bid, bflag)
	} else {
		fmt.Fprintln(os.Stderr, "Another flag has been already enrolled")
		return fmt.Errorf("Another flag has been already enrolled")
	}
	fmt.Printf("[+] Enroll { %v : %v }\n", id, flag)
	return nil
}

// GetFlag gets the flag
func GetFlag(id string) (string, error) {
	bid := []byte(id)
	db, _ := bitcask.Open("databases/flag")
	defer db.Close()
	bflag, err := db.Get(bid)
	if err != nil {
		return "", err
	}
	return string(bflag), nil
}

// DeleteFlag deletes the flag
func DeleteFlag(id string) error {
	bid := []byte(id)
	db, _ := bitcask.Open("databases/flag")
	defer db.Close()
	flag, err := db.Get(bid)
	if err != nil {
		return err
	}
	db.Delete(bid)
	fmt.Printf("[+] Delete { %v : %v }\n", id, string(flag))
	return nil
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

	EnrollFlag("0", "FLAG")

	r.PUT("/", func(c *gin.Context) {
		c.Request.ParseForm()
		problemid := c.Request.Form["id"]
		flag := c.Request.Form["flag"]
		err := EnrollFlag(problemid[0], flag[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[-] %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err})
		}

	})
	r.POST("/", func(c *gin.Context) {
		c.Request.ParseForm()
		problemid := c.Request.Form["id"]
		submittedflag := c.Request.Form["flag"]
		correct, err := CheckFlag(problemid[0], submittedflag[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[-] %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "correct": false})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err, "correct": correct})
		}

	})
	r.DELETE("/", func(c *gin.Context) {
		c.Request.ParseForm()
		problemid := c.Request.Form["id"]
		for key, value := range c.Request.Form {
			fmt.Fprintf(os.Stderr, "%v: %v\n", key, value)
		}
		err := DeleteFlag(problemid[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "[-] %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": err})
		}
	})
	r.Run()
}
