package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Phase-R/Phase-R-Backend/activities/db"
	"github.com/Phase-R/Phase-R-Backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
    "github.com/nrednav/cuid2"
	"gorm.io/gorm"
)


func ProgressController(ctx *gin.Context){
    type bodystruct struct {
        JWT        string `json:"jwt"` 
        ActivityID string `json:"ActivityID"`
        DrillID    string `json:"DrillID"`
    }
    var body bodystruct
    err:=ctx.ShouldBindJSON(&body);if err!=nil{
        ctx.JSON(404,gin.H{
            "messsage":"could not bind body",
        })
        return
    }
    isEmail,Email:=checkUserLoggedIn(body.JWT);
    if !isEmail {

        ctx.JSON(404,gin.H{
            "message":"some error occurred",
        })
        return
    }
    var user models.User
    if err := db.DB.Where("email = ?", Email).First(&user).Error; err != nil {
        log.Fatal(err)
        ctx.JSON(404,gin.H{
            "message":"cant fetch user id",
        })
        return
    }
    userID := user.ID
    err = MarkAllEntriesFalseFirstTime(userID, body.ActivityID, body.DrillID, db.DB)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Could not initialize progress",
        })
        return
    }



    err = CompleteDrillAndCheckProgress(userID, body.ActivityID, body.DrillID, db.DB)
	if err != nil {
        log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not update progress",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Drill and associated progress successfully updated",
	})
}

// func MarkAllEntriesFalseFirstTime(userID, activityID, drillID string, db *gorm.DB) error {
//     var userActivity models.UserActivityMapping
//     // if err := db.Where("email = ? AND activity_id = ?", userID, activityID).Preload("ActivityMappings.DrillMappings").First(&userActivity).Error; err != nil {
//     //     return err
//     // }

//     // Check if the activity has already been initialized (i.e., if any entry is marked as completed)
//     initialized := false
//     if userActivity.Completed {
//         initialized = true
//     }

//     for _, activityMapping := range userActivity.ActivityMappings {
//         if activityMapping.Completed {
//             initialized = true
//         }
//         for _, drillMapping := range activityMapping.DrillMappings {
//             if drillMapping.Completed {
//                 initialized = true
//             }
//         }
//     }

//     // If not initialized, mark all entries as completed: false
//     if !initialized {
//         userActivity.Completed = false
//         log.Println(userID)
//         userActivity.Email=userID
//         userActivity.ActivityID=activityID
//         userActivity.Completed=false
//         if err := db.Save(&userActivity).Error; err != nil {
//             return err
//         }

//         for _, activityMapping := range userActivity.ActivityMappings {
//             activityMapping.Completed = false
//             if err := db.Save(&activityMapping).Error; err != nil {
//                 return err
//             }

//             for _, drillMapping := range activityMapping.DrillMappings {
//                 drillMapping.Completed = false
//                 if err := db.Save(&drillMapping).Error; err != nil {
//                     return err
//                 }
//             }
//         }
//     }

//     return nil
// }


func MarkAllEntriesFalseFirstTime(userID, activityID, drillID string, db *gorm.DB) error {
    db.AutoMigrate(&models.ActivityAndActivityTypeMapping{},&models.UserActivityMapping{},&models.DrillCompletion{})
    var userActivity models.UserActivityMapping

    // Try to find the user activity record
    err := db.Where("user_id = ? AND activity_id = ?", userID, activityID).Preload("ActivityMappings.DrillMappings").First(&userActivity).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            // Record not found, create a new one
            userActivity = models.UserActivityMapping{
                ID:cuid2.Generate(),
                UserID:     userID,
                ActivityID: activityID,
                Completed:  false,
            }

            // Create the initial mappings for activity and drills
            newActivityMappings := []models.ActivityAndActivityTypeMapping{}
            activityTypes := []models.ActivityType{}

            // Get all activity types related to the activity
            if err := db.Where("activity_id = ?", activityID).Find(&activityTypes).Error; err != nil {
                return err
            }

            for _, activityType := range activityTypes {
                drillMappings := []models.DrillCompletion{}

                // Get all drills related to the activity type
                drills := []models.Drill{}
                if err := db.Where("activity_type_id = ?", activityType.ID).Find(&drills).Error; err != nil {
                    return err
                }

                for _, drill := range drills {
                    
                    drillMappings = append(drillMappings, models.DrillCompletion{
                        ID:cuid2.Generate(),
                        Drill:    drill.ID,
                        Completed:  false,
                    })
                }

                newActivityMappings = append(newActivityMappings, models.ActivityAndActivityTypeMapping{
                    ID:cuid2.Generate(),
                    ActivityType:activityType.Title,
                    UserActivityMapID: activityType.ID,
                    Completed:      false,
                    DrillMappings:  drillMappings,
                })
            }

            userActivity.ActivityMappings = newActivityMappings

            // Save the new user activity record
            if err := db.Create(&userActivity).Error; err != nil {
                return err
            }
        } else {
            // Some other database error
            return err
        }
    } else {
        // Record found, check if it's initialized
        initialized := userActivity.Completed

        for _, activityMapping := range userActivity.ActivityMappings {
            if activityMapping.Completed {
                initialized = true
            }
            for _, drillMapping := range activityMapping.DrillMappings {
                if drillMapping.Completed {
                    initialized = true
                }
            }
        }

        // If not initialized, mark all entries as completed: false
        if !initialized {
            userActivity.Completed = false
            if err := db.Save(&userActivity).Error; err != nil {
                return err
            }

            for _, activityMapping := range userActivity.ActivityMappings {
                activityMapping.Completed = false
                if err := db.Save(&activityMapping).Error; err != nil {
                    return err
                }

                for _, drillMapping := range activityMapping.DrillMappings {
                    drillMapping.Completed = false
                    if err := db.Save(&drillMapping).Error; err != nil {
                        return err
                    }
                }
            }
        }
    }

    return nil
}




// func CompleteDrillAndCheckProgress(userID, activityID, drillID string, db *gorm.DB) error {
// 	// Auto-migrate to ensure the schema is updated
// 	db.AutoMigrate(&models.DrillCompletion{})

// 	// Find or create the drill completion entry
// 	var drillCompletion models.DrillCompletion
// 	err := db.Where("activity_map_id = ? AND drill = ?", activityID, drillID).First(&drillCompletion).Error
// 	if err != nil && err == gorm.ErrRecordNotFound {
// 		// Record not found, create a new one
// 		drillCompletion = models.DrillCompletion{
//             ID : "123abc",
// 			ActivityMapID: activityID,
// 			Drill:         drillID,
// 			Completed:     false, // default value
// 		}
// 		if err := db.Create(&drillCompletion).Error; err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}

// 	// Check if the drill is already completed
// 	if drillCompletion.Completed {
// 		return nil // Nothing to do as the drill is already completed
// 	}

// 	// Mark the drill as completed
// 	drillCompletion.Completed = true
// 	if err := db.Save(&drillCompletion).Error; err != nil {
// 		return err
// 	}

// 	// Find or create the activity type mapping
// 	var activityMapping models.ActivityAndActivityTypeMapping
// 	err = db.Where("user_activity_map_id = ? AND activity_type = ?", userID, activityID).Preload("DrillMappings").First(&activityMapping).Error
// 	if err != nil && err == gorm.ErrRecordNotFound {
// 		// Record not found, create a new one
// 		activityMapping = models.ActivityAndActivityTypeMapping{
// 			UserActivityMapID: userID,
// 			ActivityType:      activityID,
// 			Completed:         false, // default value
// 		}
// 		if err := db.Create(&activityMapping).Error; err != nil {
// 			return err
// 		}
// 	} else if err != nil {
// 		return err
// 	}

// 	// Check if all drills are completed for this activity type
// 	allDrillsCompleted := true
// 	for _, drill := range activityMapping.DrillMappings {
// 		if !drill.Completed {
// 			allDrillsCompleted = false
// 			break
// 		}
// 	}

// 	// If all drills are completed, mark the activity type as completed
// 	if allDrillsCompleted {
// 		activityMapping.Completed = true
// 		if err := db.Save(&activityMapping).Error; err != nil {
// 			return err
// 		}

// 		// Find or create the user activity mapping
// 		var userActivityMapping models.UserActivityMapping
// 		err = db.Where("email = ? AND activity_id = ?", userID, activityID).Preload("ActivityMappings").First(&userActivityMapping).Error
// 		if err != nil && err == gorm.ErrRecordNotFound {
// 			// Record not found, create a new one
// 			userActivityMapping = models.UserActivityMapping{
// 				Email:        userID,
// 				ActivityID:   activityID,
// 				Completed:    false, // default value
// 			}
// 			if err := db.Create(&userActivityMapping).Error; err != nil {
// 				return err
// 			}
// 		} else if err != nil {
// 			return err
// 		}

// 		// Check if all activity types are completed for this activity
// 		allActivityTypesCompleted := true
// 		for _, activityType := range userActivityMapping.ActivityMappings {
// 			if !activityType.Completed {
// 				allActivityTypesCompleted = false
// 				break
// 			}
// 		}

// 		// If all activity types are completed, mark the activity as completed
// 		if allActivityTypesCompleted {
// 			userActivityMapping.Completed = true
// 			if err := db.Save(&userActivityMapping).Error; err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func CompleteDrillAndCheckProgress(userID, activityID, drillID string, db *gorm.DB) error {
// 	// Find or create the drill completion entry
// 	var drillCompletion models.DrillCompletion
// 	err := db.Where("activity_map_id = ? AND drill = ?", activityID, drillID).First(&drillCompletion).Error
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		return err
// 	}

// 	// Create a new record if it does not exist
// 	if err == gorm.ErrRecordNotFound {
// 		// Check if the related activity type mapping exists
// 		var activityMapping models.ActivityAndActivityTypeMapping
// 		if err := db.Where("id = ?", activityID).First(&activityMapping).Error; err != nil {
// 			return fmt.Errorf("related activity mapping does not exist: %w", err)
// 		}

// 		drillCompletion = models.DrillCompletion{
// 			ID:             "123abc", // Generate a unique ID as needed
// 			ActivityMapID:  activityID,
// 			Drill:          drillID,
// 			Completed:      false,
// 		}
// 		if err := db.Create(&drillCompletion).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// Check if the drill is already completed
// 	if drillCompletion.Completed {
// 		return nil // Nothing to do as the drill is already completed
// 	}

// 	// Mark the drill as completed
// 	drillCompletion.Completed = true
// 	if err := db.Save(&drillCompletion).Error; err != nil {
// 		return err
// 	}

// 	// Find the activity type mapping
// 	var activityMapping models.ActivityAndActivityTypeMapping
// 	if err := db.Where("user_activity_map_id = ? AND activity_type = ?", userID, activityID).Preload("DrillMappings").First(&activityMapping).Error; err != nil {
// 		return err
// 	}

// 	// Check if all drills are completed for this activity type
// 	allDrillsCompleted := true
// 	for _, drill := range activityMapping.DrillMappings {
// 		if !drill.Completed {
// 			allDrillsCompleted = false
// 			break
// 		}
// 	}

// 	// If all drills are completed, mark the activity type as completed
// 	if allDrillsCompleted {
// 		activityMapping.Completed = true
// 		if err := db.Save(&activityMapping).Error; err != nil {
// 			return err
// 		}

// 		// Find the user activity mapping
// 		var userActivityMapping models.UserActivityMapping
// 		if err := db.Where("email = ? AND activity_id = ?", userID, activityID).Preload("ActivityMappings").First(&userActivityMapping).Error; err != nil {
// 			return err
// 		}

// 		// Check if all activity types are completed for this activity
// 		allActivityTypesCompleted := true
// 		for _, activityType := range userActivityMapping.ActivityMappings {
// 			if !activityType.Completed {
// 				allActivityTypesCompleted = false
// 				break
// 			}
// 		}

// 		// If all activity types are completed, mark the activity as completed
// 		if allActivityTypesCompleted {
// 			userActivityMapping.Completed = true
// 			if err := db.Save(&userActivityMapping).Error; err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }


// func CompleteDrillAndCheckProgress(userID, activityID, drillID string, db *gorm.DB) error {
//     // Find the user activity mapping
//     var userActivity models.UserActivityMapping
//     err := db.Where("user_id = ? AND activity_id = ?", userID, activityID).
//         Preload("ActivityMappings.DrillMappings").
//         First(&userActivity).Error
//     if err != nil {
//         return err
//     }

//     // Mark the specific drill as completed
//     drillCompleted := false
//     for _, activityMapping := range userActivity.ActivityMappings {
//         for i, drillMapping := range activityMapping.DrillMappings {
//             if drillMapping.Drill == drillID {
//                 userActivity.ActivityMappings[i].DrillMappings[i].Completed = true
//                 drillCompleted = true
//                 break
//             }
//         }
//         if drillCompleted {
//             break
//         }
//     }

//     if !drillCompleted {
//         return fmt.Errorf("drill with ID %s not found in user activity mapping", drillID)
//     }

//     // Check if all drills in this activity type are completed
//     for _, activityMapping := range userActivity.ActivityMappings {
//         allDrillsCompleted := true
//         for _, drillMapping := range activityMapping.DrillMappings {
//             if !drillMapping.Completed {
//                 allDrillsCompleted = false
//                 break
//             }
//         }

//         // If all drills for this activity type are completed, mark the activity type as completed
//         if allDrillsCompleted {
//             activityMapping.Completed = true
//         } else {
//             activityMapping.Completed = false
//         }
//     }

//     // Check if all activity types in this activity are completed
//     allTypesCompleted := true
//     for _, activityMapping := range userActivity.ActivityMappings {
//         if !activityMapping.Completed {
//             allTypesCompleted = false
//             break
//         }
//     }

//     // Mark the entire activity as completed if all activity types are completed
//     userActivity.Completed = allTypesCompleted

//     // Save the updated user activity and its mappings
//     if err := db.Save(&userActivity).Error; err != nil {
//         return err
//     }

//     for _, activityMapping := range userActivity.ActivityMappings {
//         if err := db.Save(&activityMapping).Error; err != nil {
//             return err
//         }

//         for _, drillMapping := range activityMapping.DrillMappings {
//             if err := db.Save(&drillMapping).Error; err != nil {
//                 return err
//             }
//         }
//     }

//     return nil
// }

func CompleteDrillAndCheckProgress(userID, activityID, drillID string, db *gorm.DB) error {
    var userActivity models.UserActivityMapping

    // Fetch the user activity with related activity mappings and drill mappings
    if err := db.Where("user_id = ? AND activity_id = ?", userID, activityID).
        Preload("ActivityMappings.DrillMappings").
        First(&userActivity).Error; err != nil {
        return err
    }

    // Flag to check if the drill was found
    drillFound := false

    // Iterate over the activity mappings
    for i := range userActivity.ActivityMappings {
        // Iterate over the drill mappings
        for j := range userActivity.ActivityMappings[i].DrillMappings {
            if userActivity.ActivityMappings[i].DrillMappings[j].Drill == drillID {
                drillFound = true

                // Mark the drill as completed
                if !userActivity.ActivityMappings[i].DrillMappings[j].Completed {
                    userActivity.ActivityMappings[i].DrillMappings[j].Completed = true

                    // Save the updated drill mapping
                    if err := db.Save(&userActivity.ActivityMappings[i].DrillMappings[j]).Error; err != nil {
                        return err
                    }
                }
                break // Exit the loop after finding and processing the drill
            }
        }
    }

    if !drillFound {
        return fmt.Errorf("drill ID %s not found for activity ID %s", drillID, activityID)
    }

    return nil
}



func checkUserLoggedIn(jwt string)(bool,string){
    claims,err:=verifyToken(jwt);if err!=nil{
        log.Println(err)
		return false,""
    }
    ID:=claims["iss"].(string)
    return true,ID
}



func verifyToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("SECRET_KEY")
    log.Println(secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("unable to extract claims")
}





func GetUserProgress(ctx *gin.Context) {
    type ProgressResponse struct {
        ActivityID   string `json:"activity_id"`
        ActivityTitle string `json:"activity_title"`
        Completed    bool   `json:"completed"`
        ActivityTypes []struct {
            ActivityTypeID string `json:"activity_type_id"`
            Title          string `json:"title"`
            Completed      bool   `json:"completed"`
            Drills         []struct {
                DrillID   string `json:"drill_id"`
                Title     string `json:"title"`
                Completed bool   `json:"completed"`
            } `json:"drills"`
        } `json:"activity_types"`
    }

    // Extract JWT from request query
    jwtToken := ctx.Query("jwt")

    // Validate JWT and extract user email
    isEmail, Email := checkUserLoggedIn(jwtToken)
    if !isEmail {
        ctx.JSON(404, gin.H{
            "message": "invalid JWT",
        })
        return
    }

    // Fetch the user based on email
    var user models.User
    if err := db.DB.Where("email = ?", Email).First(&user).Error; err != nil {
        ctx.JSON(404, gin.H{
            "message": "can't fetch user",
        })
        return
    }

    userID := user.ID

    // Fetch the user's activity mappings
    var userActivities []models.UserActivityMapping
    err := db.DB.Where("user_id = ?", userID).
        Preload("ActivityMappings.DrillMappings").
        Preload("ActivityMappings").
        Preload("ActivityMappings.DrillMappings").
        Find(&userActivities).Error
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "can't fetch user activities",
        })
        return
    }

    // Prepare the progress response
    progressResponses := []ProgressResponse{}
    for _, userActivity := range userActivities {
        var activity models.Activities
        if err := db.DB.Where("id = ?", userActivity.ActivityID).First(&activity).Error; err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{
                "message": "can't fetch activity details",
            })
            return
        }

        progressResponse := ProgressResponse{
            ActivityID:   userActivity.ActivityID,
            ActivityTitle: activity.Title,
            Completed:    userActivity.Completed,
        }

        for _, activityMapping := range userActivity.ActivityMappings {
            activityTypeProgress := struct {
                ActivityTypeID string `json:"activity_type_id"`
                Title          string `json:"title"`
                Completed      bool   `json:"completed"`
                Drills         []struct {
                    DrillID   string `json:"drill_id"`
                    Title     string `json:"title"`
                    Completed bool   `json:"completed"`
                } `json:"drills"`
            }{
                ActivityTypeID: activityMapping.ID,
                Title:          activityMapping.ActivityType,
                Completed:      activityMapping.Completed,
            }

            for _, drillMapping := range activityMapping.DrillMappings {
                var drill models.Drill
                if err := db.DB.Where("id = ?", drillMapping.Drill).First(&drill).Error; err != nil {
                    ctx.JSON(http.StatusInternalServerError, gin.H{
                        "message": "can't fetch drill details",
                    })
                    return
                }

                drillProgress := struct {
                    DrillID   string `json:"drill_id"`
                    Title     string `json:"title"`
                    Completed bool   `json:"completed"`
                }{
                    DrillID:   drill.ID,
                    Title:     drill.Title,
                    Completed: drillMapping.Completed,
                }

                activityTypeProgress.Drills = append(activityTypeProgress.Drills, drillProgress)
            }

            progressResponse.ActivityTypes = append(progressResponse.ActivityTypes, activityTypeProgress)
        }

        progressResponses = append(progressResponses, progressResponse)
    }

    ctx.JSON(http.StatusOK, gin.H{
        "progress": progressResponses,
    })
}
