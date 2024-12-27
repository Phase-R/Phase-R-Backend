package controllers

import (
	"os"
	"log"
	"bytes"
	"errors"
	"context"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type DietParams struct {
	Plan             string `json:"plan"`
	Activity         string `json:"activity"`
	TargetCal        string `json:"target_cal"`
	TargetProtein    string `json:"target_protein"`
	TargetFat        string `json:"target_fat"`
	TargetCarbs      string `json:"target_carbs"`
	Cuisine          string `json:"cuisine"`
	MealChoice       string `json:"meal_choice"`
	Occupation       string `json:"occupation"`
	Allergies        string `json:"allergies"`
	OtherPreferences string `json:"other_preferences"`
	Variety          string `json:"variety"`
	Budget           string `json:"budget"`
}

type Message struct {
	Role	string	`json:"role"`
	Content	string	`json:"content"`
}

func validateParams(params DietParams) error {
	// Check for invalid or missing values
	for _, v := range []string{
		params.Plan, params.Activity, params.TargetCal, params.TargetProtein,
		params.TargetFat, params.TargetCarbs, params.Cuisine, params.MealChoice,
		params.Occupation, params.Allergies, params.OtherPreferences, params.Variety, params.Budget,
	} {
		if v == "" || v == "unknown" {
			return errors.New("some parameters are invalid or missing")
		}
	}
	return nil
}

func Monthly_Diet_Gen(ctx *gin.Context) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file!")
	}
	// ctx.JSON(http.StatusAccepted, gin.H{"message": "Hello World!"})
	api_key := os.Getenv("GITHUB_TOKEN")

	client := openai.NewClient(
		option.WithAPIKey(api_key),
		option.WithBaseURL("https://models.inference.ai.azure.com"),
	)

	const promptTemplate = `
	Generate a meal plan for {{.Plan}} with a daily activity level being {{.Activity}}.
	Target calories: {{.TargetCal}} kcal. Macro Distribution: {{.TargetProtein}} g protein, 
	{{.TargetCarbs}} g of carbs and {{.TargetFat}} g of fat. Give this in the form of a table with days: 
	[Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday] with meals being mandatory 
	[Breakfast, Afternoon Snack, Lunch, Evening Snack, Dinner]. Don't tell anything other than the table 
	in the form of HTML table as mentioned. The foods should mainly belong to {{.Cuisine}} cuisine and should be {{.MealChoice}}.
	`

	var params DietParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input!"})
		return
	}

	if err := validateParams(params); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"headers": gin.H{"X-Error": "Some invalid parameters were found!!!"},
		})
		return
	}

	tmp, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template!"})
		return
	}

	var finalPrompt bytes.Buffer
	err = tmp.Execute(&finalPrompt, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute template!"})
		return
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(finalPrompt.String()),
		}),
		Model: openai.F("gpt-4o"),
	})
	if err != nil {
		panic(err.Error())
	}
	
	response := chatCompletion.Choices[0].Message.Content
	ctx.JSON(http.StatusOK, gin.H{"message": response})
}