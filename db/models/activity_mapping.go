package models

type UserActivityMapping struct {
    ID               string                       `json:"id" gorm:"primaryKey"`
    UserID           string                       `json:"user_id"`
    ActivityID       string                       `json:"activity_id"`
    ActivityMappings []ActivityAndActivityTypeMapping `json:"activity_mappings" gorm:"foreignKey:UserActivityMapID"`
    Completed        bool                         `json:"completed"`
}

type ActivityAndActivityTypeMapping struct {
    ID                string                      `json:"id" gorm:"primaryKey"`
    UserActivityMapID string                      `json:"user_activity_map_id"`
    ActivityType      string                      `json:"activity_type"`
    DrillMappings     []DrillCompletion           `json:"drill_mappings" gorm:"foreignKey:ActivityMapID"`
    Completed         bool                        `json:"completed"`
}

type DrillCompletion struct {
    ID              string `json:"id" gorm:"primaryKey"`
    ActivityMapID   string `json:"activity_map_id"`
    Drill           string `json:"drill"`
    Completed       bool   `json:"completed"`
}
