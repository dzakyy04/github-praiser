package handlers

import (
	"github-praiser/services"

	"github.com/gofiber/fiber/v2"
)

type PraisingRequest struct {
	Username string `json:"username"`
}

func HandlePraising(ctx *fiber.Ctx) error {
	var request PraisingRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    400,
				"message": "Bad Request",
			},
			"error": "Invalid request body",
		})
	}

	if request.Username == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    400,
				"message": "Bad Request",
			},
			"error": "Username is required",
		})
	}

	githubData, err := services.GetGithubData(request.Username)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    500,
				"message": "Internal Server Error",
			},
			"error": err.Error(),
		})
	}

	prompt := services.CreatePrompt(request.Username, githubData.Data, githubData.Readme)

	praising, err := services.GenerateAIResponse(prompt)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": fiber.Map{
				"code":    500,
				"message": "Internal Server Error",
			},
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.Map{
			"code":    200,
			"message": "OK",
		},
		"data": fiber.Map{
			"praising": praising,
		},
	})
}
