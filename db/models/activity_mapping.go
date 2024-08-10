package models

import "gorm.io/gorm"



// type DrillCompletion struct {
//     gorm.Model
//     ID              string `json:"id" gorm:"primaryKey"`
//     ActivityMapID   string `json:"activity_map_id"`
//     DrillID           string `json:"drill"`
// }

// type ActivityAndActivityTypeMapping struct {
//     gorm.Model
//     ID                string                      `json:"id" gorm:"primaryKey"`
//     UserActivityMapID string                      `json:"user_activity_map_id"`
//     ActivityType      string                      `json:"activity_type"`
//     DrillMappings     []DrillCompletion           `json:"drill_mappings" gorm:"foreignKey:ActivityMapID"`
// }


// type UserActivityMapping struct {
//     gorm.Model
//     ID               string                       `json:"id" gorm:"primaryKey"`
//     UserID           string                       `json:"user_id"`
//     ActivityID       string                       `json:"activity_id"`
//     ActivityMappings []ActivityAndActivityTypeMapping `json:"activity_mappings" gorm:"foreignKey:UserActivityMapID"`
// }

type DrillCompletion struct {
    gorm.Model
    ID              string `json:"id" gorm:"primaryKey"`
    ActivityMapID   string `json:"activity_map_id"`
    DrillID         string `json:"drill"`
}

type ActivityAndActivityTypeMapping struct {
    gorm.Model
    ID                string                      `json:"id" gorm:"primaryKey"`
    UserActivityMapID string                      `json:"user_activity_map_id"`
    ActivityType      string                      `json:"activity_type"`
    DrillMappings     []DrillCompletion           `json:"drill_mappings" gorm:"foreignKey:ActivityMapID"`
}

type UserActivityMapping struct {
    gorm.Model
    ID               string                       `json:"id" gorm:"primaryKey"`
    UserID           string                       `json:"user_id"`
    ActivityID       string                       `json:"activity_id"`
    ActivityMappings []ActivityAndActivityTypeMapping `json:"activity_mappings" gorm:"foreignKey:UserActivityMapID"`
}



