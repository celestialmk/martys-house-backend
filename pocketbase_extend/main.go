package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v5"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	// Load .env file
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// environmentPath := filepath.Join(dir, ".env")
	// err = godotenv.Load(environmentPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))

		e.Router.GET("/hello/:name", func(c echo.Context) error {
			name := c.PathParam("name")

			return c.JSON(http.StatusOK, map[string]string{"message": "Hello there this is a new start " + name})
		} /* optional middlewares */)

		// e.Router.POST("/api/v1/chat", func(c echo.Context) error {
		// 	data := apis.RequestInfo(c).Data

		// 	question := data["message"]

		// 	var str_question string

		// 	str_question, ok := question.(string)

		// 	if !ok {
		// 		log.Fatal("Question not string")
		// 	}

		// 	ctx := context.Background()
		// 	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	model := client.GenerativeModel("gemini-1.5-flash")

		// 	resp, err := model.GenerateContent(ctx, genai.Text(str_question))

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	return c.JSON(http.StatusOK, map[string]any{"response": resp.Candidates[0].Content.Parts})

		// })

		return nil
	})

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

}
