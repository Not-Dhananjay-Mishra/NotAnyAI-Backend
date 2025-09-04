package router

import (
	"encoding/base64"
	"log"
	"os"
	mongodb "server/database/MongoDB"
	"server/models"
	codingmodel "server/models/CodingModel"
	"server/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type recived struct {
	Agent string `json:"agent"`
	Query string `json:"query"`
	Img   string `json:"img"`
}

func HandleConn(conn *websocket.Conn, username string) {
	defer func() {
		log.Println(utils.Magenta("Cleaning up user: ", username))
		//conn.ReadJSON(utils.Response{Text: "Error Encountered Report to Developer"})
		delete(utils.LiveConn, conn)
		delete(utils.UserConn, username)
		delete(utils.ConnUser, conn)
		delete(utils.MemoryStore, username)

		conn.Close()
	}()
	client := models.GeminiModel()
	for {
		var data recived
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println("error reading data")
			return
		}
		if data.Agent == "normal" {
			log.Println(utils.Green("Json"))

		} else if data.Img != "" {
			imgBytes, _ := base64.StdEncoding.DecodeString(data.Img)
			sus := uuid.New().String()
			fileName := "uploads/" + sus + "-" + username + ".png"
			err = os.WriteFile(fileName, imgBytes, 0644)
			if err != nil {
				log.Println("File write error:", err)
			}
			userPrompt := data.Query
			log.Println(utils.Green("Image received and saved as", fileName))
			models.AddImgToMemoryUSER(username, userPrompt, fileName, imgBytes)
			//fmt.Println(utils.MemoryStore[username])
			conn.WriteJSON(utils.Response{Text: "Image received successfully"})
			//models.PrintMemobyUsername(username)
			aires := models.ImageModel(client, fileName, userPrompt, "png", conn, username, imgBytes)
			//aires := models.ModelWithTools(client, utils.MemoryStore[username], username, conn)
			conn.WriteJSON(utils.Response{Text: aires})
			continue

		} else if data.Agent == "code" {
			replyCh := make(chan map[string]string, 1)
			codingmodel.UnderProcessCode <- codingmodel.Process{Data: codingmodel.Sus{Query: data.Query}, ReplyCh: replyCh, Conn: conn}
			result := <-replyCh
			conn.WriteJSON(result)
			continue
		}

		receivedData := data
		log.Println(utils.Blue(receivedData.Query))
		if receivedData.Query != "" {
			linit, _ := mongodb.GetUserLimit(username)
			if linit <= 0 {
				conn.WriteJSON(map[string]string{"processing": "You have exhausted your limit. Please talk to developer to continue using the service."})
				continue
			}
			mongodb.UpdateLimit(username, linit-1)
			prompt := receivedData.Query
			models.AddToMemoryUSER(username, prompt)
			aires := models.ModelWithTools(client, utils.MemoryStore[username], username, conn)
			conn.WriteJSON(utils.Response{Text: aires})
		}
	}

}
