// 1. Package declaration
// Every Go file starts with a package declaration.
// The 'main' package is special - it defines a standalone executable program, not a library.
package main

// 2. Import statements
// These import external packages that this program will use.
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// 3. Constants
// This defines a constant for the OpenAI API URL.
// Constants in Go are declared using the 'const' keyword.
const openAIURL = "https://api.openai.com/v1/chat/completions"

// 4. Global variables
// This declares a global variable to store the OpenAI API key.
// In Go, variables declared outside of functions are package-level variables.
var openAIKey string

// 5. Struct definitions
// Structs in Go are used to create custom data types.
// The `json` tags are used for JSON marshaling and unmarshaling.

// Message represents a single message in the chat.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents the structure of a request to the OpenAI API.
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// OpenAIResponse represents the structure of a response from the OpenAI API.
type OpenAIResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

// WebSocketMessage represents a message sent over WebSocket.
type WebSocketMessage struct {
	Text string `json:"text"`
}

// 6. More global variables
// This creates a map to store active WebSocket connections.
// The 'var' block allows declaring multiple variables together.
var (
	clients = make(map[*websocket.Conn]bool)
)

// 7. Main function
// The main function is the entry point of the Go program.
func main() {
	// 8. Environment variable retrieval
	// os.Getenv retrieves the value of an environment variable.
	openAIKey = os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable")
		return
	}

	// 9. Fiber app initialization
	// This creates a new instance of the Fiber web framework.
	app := fiber.New()

	// 10. Static file serving
	// This tells Fiber to serve static files from the "./static" directory.
	app.Static("/", "./static")

	// 11. Route handlers
	// These set up the routes for the web application.
	app.Get("/", handleHome)
	app.Get("/ws", websocket.New(handleWebSocket))

	// 12. Port configuration
	// This gets the port from an environment variable, or uses a default.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 13. Start the server
	// This starts the Fiber server on the specified port.
	fmt.Printf("Server starting on :%s\n", port)
	app.Listen(":" + port)
}

// 14. Home route handler
// This function handles requests to the root ("/") path.
func handleHome(c *fiber.Ctx) error {
	// It sends the index.html file as the response.
	return c.SendFile("./static/index.html")
}

// 15. WebSocket handler
// This function handles WebSocket connections.
func handleWebSocket(c *websocket.Conn) {
	// 16. Add client to the clients map
	// The clients map keeps track of all active WebSocket connections.
	clients[c] = true
	// This defers the removal of the client from the map until the function returns.
	defer delete(clients, c)

	// 17. Infinite loop to handle incoming messages
	for {
		var msg WebSocketMessage
		// ReadJSON reads a JSON message from the WebSocket connection.
		err := c.ReadJSON(&msg)
		if err != nil {
			break
		}
		// Start a new goroutine to handle the response streaming.
		// This allows multiple clients to be served concurrently.
		go streamResponse(msg.Text, c)
	}
}

// 18. Response streaming function
// This function streams responses from the OpenAI API to the client.
func streamResponse(message string, conn *websocket.Conn) {
	// 19. Prepare OpenAI API request
	openAIReq := OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "user", Content: message},
		},
		Stream: true,
	}
	// Marshal the request into JSON.
	reqBody, _ := json.Marshal(openAIReq)

	// 20. Create and send HTTP request to OpenAI API
	req, _ := http.NewRequest("POST", openAIURL, strings.NewReader(string(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling OpenAI API:", err)
		return
	}
	// Ensure the response body is closed when the function returns.
	defer resp.Body.Close()

	// 21. Read the streaming response
	reader := bufio.NewReader(resp.Body)
	isFirstToken := true
	for {
		// Read each line of the stream.
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading stream:", err)
			break
		}

		// 22. Process each line
		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}
		line = strings.TrimPrefix(line, "data: ")
		var aiResp OpenAIResponse
		err = json.Unmarshal([]byte(line), &aiResp)
		if err != nil {
			continue
		}

		// 23. Send processed content to WebSocket client
		if len(aiResp.Choices) > 0 {
			content := aiResp.Choices[0].Delta.Content
			if content != "" {
				if isFirstToken {
					// Send first token with "AI: " prefix.
					conn.WriteJSON(WebSocketMessage{Text: "AI: " + content})
					isFirstToken = false
				} else {
					// Send subsequent tokens without prefix.
					conn.WriteJSON(WebSocketMessage{Text: content})
				}
			}
		}
	}
}
