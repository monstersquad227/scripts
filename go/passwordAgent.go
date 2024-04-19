package main

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const DecryptKey = "2acb1bc6-7a90-47"
const MysqlUsername = "root"
const MysqlPassword = "1qaz@WSX"
const MysqlAddress = "192.168.1.87"
const MysqlPort = "3306"
const MysqlDatabases = "devops"
const MysqlCharset = "utf8"

func AesDecryptByGCM(Password, key string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(Password)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(dataByte) < nonceSize {
		customErr := errors.New("dataByte to short")
		return "", customErr
	}
	nonce, ciphertext := dataByte[:nonceSize], dataByte[nonceSize:]
	open, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(open), nil
}

func main() {
	app := gin.Default()
	app.POST("/machine/modify/password", func(c *gin.Context) {
		type machinePassword struct {
			ID       int    `json:"id"`
			IP       string `json:"ip"`
			Password string `json:"password"`
		}
		// 解析json
		mp := machinePassword{}
		err := c.ShouldBindJSON(&mp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		// 获取解密密码
		password, err := AesDecryptByGCM(mp.Password, DecryptKey)
		if err != nil {
			log.Println("password Failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		// 执行修改密码
		cmd := exec.Command("chpasswd")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Println("Failed to open pipe to command:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		go func() {
			defer stdin.Close()
			fmt.Fprintf(stdin, "root:%s\n", password)
		}()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Println("Exec Command Failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		// 修改数据库密码
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
			MysqlUsername,
			MysqlPassword,
			MysqlAddress,
			MysqlPort,
			MysqlDatabases,
			MysqlCharset)
		MysqlEngine, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Connect Databases Failed: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		query := "UPDATE machine SET instance_password = ? WHERE id = ?"
		result, err := MysqlEngine.Exec(query, mp.Password, mp.ID)
		if err != nil {
			log.Println("Databases Exec Failed: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    20000,
				"message": "failed",
				"data":    err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "successful",
			"data":    result,
		})
	})
	err := app.Run("0.0.0.0:38080")
	if err != nil {
		return
	}
}
