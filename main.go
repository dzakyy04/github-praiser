package main

import (
	"github-praiser/services"

	"github.com/gofiber/fiber/v2"
)

type GitHubData struct {
	Data   interface{}
	Readme string
}

func main() {
	app := fiber.New()

	app.Get("/:username", func(ctx *fiber.Ctx) error {
		username := ctx.Params("username")

		githubData, err := services.GetGithubData(username)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(githubData)
	})

	app.Listen(":3000")
}
