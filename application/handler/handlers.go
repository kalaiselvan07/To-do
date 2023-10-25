package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

type JwtClaims struct {
	CompanyId string `json:"companyId,omitempty"`
	Username  string `json:"username,omitempty"`
	Role      string `json:"role,omitempty"`
	jwt.StandardClaims
}

type todo struct {
	ID        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

const ip = "localhost"

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(ip, true) {
		return nil
	}
	return fmt.Errorf("token is invalid")
}

var db *bbolt.DB

func InitDB(dbPath string) error {
	var err error
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Todos"))
		return err
	})
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func GetTodos(c *gin.Context) {
	todos := make([]todo, 0)

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var t todo
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			todos = append(todos, t)
		}

		return nil
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve todos"})
		return
	}

	c.IndentedJSON(http.StatusOK, todos)
}

func GetTodoById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	var t todo
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))
		v := b.Get([]byte(strconv.Itoa(id)))

		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &t)
	})

	if err != nil || t.ID == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func AddTodo(c *gin.Context) {
	var t todo
	if err := c.BindJSON(&t); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))

		id, _ := b.NextSequence()
		t.ID = int(id)

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(t.ID)), buf)
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to add todo"})
		return
	}

	c.IndentedJSON(http.StatusCreated, t)
}

func UpdateTodoById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	var t todo
	err = c.BindJSON(&t)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	t.ID = id
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))
		if b.Get([]byte(strconv.Itoa(id))) == nil {
			return nil
		}

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(id)), buf)
	})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to update todo"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func Authenticate(c *gin.Context) {
	reqBody := c.Request.Body
	url := "http://authn:8088/authenticate"
	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		log.Printf("Error: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseMap map[string]string
	if err := json.Unmarshal(body, &responseMap); err != nil {
		log.Printf("Error unmarshalling response: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(resp.StatusCode, responseMap)

}

func ValidateToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		validationEndpoint := "http://authn:8088/validate"
		apiKey := context.GetHeader("apikey")

		req, err := http.NewRequest("POST", validationEndpoint, nil)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			context.Abort()
			return
		}
		req.Header.Add("apikey", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to validate token"})
			context.Abort()
			return
		}
		defer resp.Body.Close()

		var claims JwtClaims
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&claims); err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse validation response"})
			context.Abort()
			return
		}

		context.Set("claims", claims)
		context.Next()
	}
}

func AuthorizeRequest(role string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var jwtClaims JwtClaims
		if claims, exists := context.Get("claims"); exists {
			jwtClaims = claims.(JwtClaims)
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No claims provided"})
			return
		}

		claimsBytes, err := json.Marshal(jwtClaims)
		if err != nil {
			log.Printf("Error encoding claims: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding claims"})
			return
		}

		baseURL := "http://authz:8089/authorize/"
		req, err := http.NewRequest("POST", baseURL+role, bytes.NewBuffer(claimsBytes))
		if err != nil {
			log.Printf("Error creating request: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error executing request: %s", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
	}
}
