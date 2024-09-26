package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		e.Router.GET("/hello/:name", func(c echo.Context) error {
			name := c.PathParam("name")

			return c.JSON(http.StatusOK, map[string]string{"message": "Hello there this is a new start " + name})
		} /* optional middlewares */)

		e.Router.POST("/api/v1/chat", func(c echo.Context) error {
			data := apis.RequestInfo(c).Data

			question := data["message"]

			var str_question string

			str_question, ok := question.(string)

			if !ok {
				log.Fatal("Question not string")
			}

			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

			if err != nil {
				log.Fatal(err)
			}
			model := client.GenerativeModel("gemini-1.5-flash")

			resp, err := model.GenerateContent(ctx, genai.Text(str_question))

			if err != nil {
				log.Fatal(err)
			}

			return c.JSON(http.StatusOK, map[string]any{"response": resp.Candidates[0].Content.Parts})

		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
