# Go HTMX LLM Chat Application

This project is a real-time chat application built with Go, HTMX, and WebSockets, integrating with an AI language model for interactive conversations.

## Features

- Real-time chat interface
- AI-powered responses
- Markdown rendering for rich text formatting
- Code syntax highlighting
- Responsive design with Tailwind CSS

## Prerequisites

- Go 1.20 or later
- Docker (optional)

## Installation

1. Clone the repository:

```
git clone https://github.com/developersdigest/go-htmx-llm.git
cd go-htmx-llm
```

2. Set up your environment variables:

Create a `.env` file in the project root and add your OpenAI API key:

```
OPENAI_API_KEY=your_api_key_here
```

## Running the Application

### Without Docker

1. Install dependencies:

```
go mod download
```

2. Build and run the application:

```
go build -o main .
./main
```

### With Docker

1. Build the Docker image:

```
docker build -t go-htmx-llm .
```

2. Run the container:

```
docker run -p 8080:8080 --env-file .env go-htmx-llm
```

The application will be available at `http://localhost:8080`.

## Project Structure

- `main.go`: Main application file containing the server setup and WebSocket handling
- `static/index.html`: Frontend HTML file with HTMX integration
- `Dockerfile`: Instructions for building the Docker image
- `go.mod` and `go.sum`: Go module files for dependency management

## Technologies Used

- [Go](https://golang.org/): Backend server
- [Fiber](https://gofiber.io/): Web framework for Go
- [HTMX](https://htmx.org/): Frontend interactivity
- [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API): Real-time communication
- [Tailwind CSS](https://tailwindcss.com/): Styling
- [Marked](https://marked.js.org/): Markdown parsing
- [Highlight.js](https://highlightjs.org/): Code syntax highlighting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.
