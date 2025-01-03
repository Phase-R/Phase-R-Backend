package controllers

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type DietParams struct {
	Height           int     `json:"height"`
	Weight           int     `json:"weight"`
	Age              int     `json:"age"`
	BMI              float64 `json:"bmi"`
	Gender           string  `json:"gender"`
	Goal             string  `json:"goal"`
	ActivityLevel    string  `json:"activity_level"`
	Duration         int     `json:"duration"`
	TargetCal        int     `json:"target_cal"`
	TargetProtein    int     `json:"target_protein"`
	TargetFat        int     `json:"target_fat"`
	TargetCarbs      int     `json:"target_carbs"`
	Cuisine          string  `json:"cuisine"`
	MealChoice       string  `json:"meal_choice"`
	Allergies        string  `json:"allergies"`
	OtherPreferences string  `json:"other_preferences"`
	Variety          string  `json:"variety"`
}

type SubstituteParams struct {
	Food             string `json:"food"`
	Allergies        string `json:"allergies"`
	OtherPreferences string `json:"other_preferences"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func validateSubstituteParams(params SubstituteParams) error {
	// Check for invalid or missing values
	for _, v := range []string{
		params.Food, params.Allergies, params.OtherPreferences,
	} {
		if v == "" || v == "unknown" {
			return errors.New("some parameters are invalid or missing")
		}
	}
	return nil
}

func validateDietParams(params DietParams) error {
	// Check for invalid or missing values
	for _, v := range []string{
		params.Gender, params.Goal, params.ActivityLevel, params.Cuisine, params.MealChoice,
		params.Allergies, params.OtherPreferences, params.Variety,
	} {
		if v == "" || v == "unknown" {
			return errors.New("some parameters are invalid or missing")
		}
	}
	return nil
}

func Diet_Sub(ctx *gin.Context) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	api_key := os.Getenv("MODEL_TOKEN")

	client := openai.NewClient(
		option.WithAPIKey(api_key),
		option.WithBaseURL("https://models.inference.ai.azure.com"),
	)

	const promptTemplate = `Generate an alternate food for {{.Food}} with a similar nutritional profile. The food should be suitable for someone with {{.Allergies}} allergies and {{.OtherPreferences}} preferences. Mention only the dish and nothing else.`

	var params SubstituteParams

	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input!"})
		return
	}

	if err := validateSubstituteParams(params); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse template!"})
		return
	}

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(finalPrompt.String()),
		}),
		Model: openai.F("gpt-4o-mini"),
	})

	if err != nil {
		panic(err.Error())
	}

	response := chatCompletion.Choices[0].Message.Content
	ctx.JSON(http.StatusOK, gin.H{"message": response})
}

func Monthly_Diet_Gen(ctx *gin.Context) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file!")
	}
	// ctx.JSON(http.StatusAccepted, gin.H{"message": "Hello World!"})
	api_key := os.Getenv("MODEL_TOKEN")

	client := openai.NewClient(
		option.WithAPIKey(api_key),
		option.WithBaseURL("https://models.inference.ai.azure.com"),
	)

	const promptTemplate = `
    Generate a meal plan for a person with the following details: 
    Height: {{.Height}} cm, Weight: {{.Weight}} kg, Age: {{.Age}} years, BMI: {{.BMI}}, Gender: {{.Gender}}, 
    Goal: {{.Goal}}, Activity Level: {{.ActivityLevel}}, and Duration: {{.Duration}} weeks. 
    The plan should target {{.TargetCal}} kcal daily with a macro distribution of 
    {{.TargetProtein}} g of protein, {{.TargetCarbs}} g of carbs, and {{.TargetFat}} g of fat. 

    Use foods from {{.Cuisine}} cuisine and follow the {{.MealChoice}} preference. Address any allergies specified: {{.Allergies}}. 
    Ensure {{.Variety}} variety across meals and adhere to these additional preferences: {{.OtherPreferences}}. 

    Structure the meal plan over 7 days, and include the following five mandatory meals each day: 
    Breakfast, Afternoon Snack, Lunch, Evening Snack, and Dinner. 
    Output the meal plan as a JSON in the following format: 

    [
      {
        "day1": [
          {"Type": "Breakfast", "Time": "08:00 AM", "Cals": 400, "Foods": []}, 
          {"Type": "Afternoon Snack", "Time": "11:00 AM", "Cals": 200, "Foods": []}, 
          {"Type": "Lunch", "Time": "01:00 PM", "Cals": 600, "Foods": []}, 
          {"Type": "Evening Snack", "Time": "04:30 PM", "Cals": 200, "Foods": []}, 
          {"Type": "Dinner", "Time": "07:30 PM", "Cals": 500, "Foods": []}
        ]
      },
      ... 
    ]

    Do not include any introduction, explanation, or content outside of the JSON.
`

	var params DietParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input!"})
		return
	}

	if err := validateDietParams(params); err != nil {
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
