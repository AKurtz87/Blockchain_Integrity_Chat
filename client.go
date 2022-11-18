package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// ################################# START DATABASE CONNECTION ########################################

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "password"
    dbName := "chat"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    //log.Println("CONNECTION OK")
    if err != nil {
        panic(err.Error())
    }
    return db
}

// ################################# END DATABASE CONNECTION ########################################


// ################################# START MAKE HASH ################################################

func hashString (stringToHash string) string {
	h := sha256.New()
	h.Write([]byte(stringToHash))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}

// ################################# END MAKE HASH #################################################

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	log.Println("READ_PUMP")

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	db := dbConn()
	
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	
// ################################# START LOADING CHAT BACKUP ########################################

	type chat struct {
		id int
		time string
		user  string
		message string
		hash string
		proof string
		blockchain string
	}

	var messageBackup []string

	var hashHistory []string

	hashHistory = append(hashHistory, "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

	for j := 2; ; j++{
		selDBA, err := db.Query("SELECT * FROM chat WHERE id=?", j - 1)

		if err != nil {
			panic(err.Error())
		}

		empA := new (chat)
		
		for selDBA.Next() {
			var id int
			var time, user, message, hash, proof, blockchain string
			err = selDBA.Scan(&id, &time, &user, &message, &hash, &proof, &blockchain)
			if err != nil {
				panic(err.Error())
			}
			empA.id = id
			empA.time = time
			empA.user = user
			empA.message = message
			empA.hash = hash
			empA.proof = proof
			empA.blockchain = blockchain
			}

			if empA.user == "" {
				break
			} else {
				messageBackup = append(messageBackup, empA.message)
				hashHistory = append(hashHistory, empA.hash)
				
			}

			data := empA.user + " => " + empA.message

			c.hub.broadcast <- []byte(data)			
	  }

//################################### BACKUP RIPRISTINATO READY TO RECEIVE OR SENT MESSAGES #################################

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		
		c.hub.broadcast <- message
		
		messageSlice := [6]string{}

		t := time.Now()
		messageTime := t.Format("2006-01-02 15:04:05")	
		messageSlice[0] = messageTime

		messageString := string(message)

		reUser := regexp.MustCompile(`\[ (.*?)\ ]`)
		submatchallUser := reUser.FindAllString(messageString, -1)
		for _, element := range submatchallUser {
			element = strings.Trim(element, "[ ")
			element = strings.Trim(element, " ]")
			messageSlice[1] = element
		}
		

		reText := regexp.MustCompile(`\>    (.*)`)
		submatchallText := reText.FindAllString(messageString, -1)
		for _, element := range submatchallText {
			element = strings.Trim(element, "> ")
			messageSlice[2] = element
		}

		//############## GENERATE HASH ##############
			
			hash := hashString(messageSlice[2])
			
			for i := 0; ; i++ {
				
				delta := strconv.Itoa(i)
				tryHash := hashString(hash + delta)
					secret := "AAA"
					if !strings.HasPrefix(tryHash, secret) {
						//fmt.Println(tryHash, " do more work!")
						//time.Sleep(10 * time.Millisecond)
						//time.Sleep(1 * time.Second)
						continue
					} else {
						messageSlice[3] = hash
						fmt.Println(tryHash, " work done!")
						messageSlice[3] = tryHash
						hashHistory = append(hashHistory, tryHash)
						fmt.Println(hashHistory)
						messageSlice[4] = delta

							messageSlice[5] = hashString(messageSlice[2] + hashHistory[len(hashHistory) - 2])
							fmt.Println("VALORE DA ASCHARE PER BLOCKCHAIN")
							fmt.Println(messageSlice[2] + hashHistory[len(hashHistory) - 2])
							

							insForm, err := db.Prepare("INSERT INTO `chat` (`time`, `user`, `message`, `hash`, `proof`, `blockchain` ) VALUES (?, ?, ?, ?, ?, ?)")
				
								if err != nil {
									panic(err.Error())
								}
	
							insForm.Exec(messageSlice[0], messageSlice[1], messageSlice[2], messageSlice[3], messageSlice[4], messageSlice[5])

								break
					}	
				}			
	}

	defer db.Close()
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	log.Println("WRITE_PUMP")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			//fmt.Println(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}



