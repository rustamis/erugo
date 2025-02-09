package models

import "time"

type Setting struct {
	Id            string    `json:"id"`
	Value         string    `json:"value"`
	PreviousValue string    `json:"previous_value"`
	SettingGroup  string    `json:"setting_group"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SettingRequest struct {
	Id           string `json:"id"`
	Value        string `json:"value"`
	SettingGroup string `json:"setting_group"`
}

type SettingResponse struct {
	Id            string    `json:"id"`
	Value         string    `json:"value"`
	PreviousValue string    `json:"previous_value"`
	SettingGroup  string    `json:"setting_group"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *SettingRequest) ToSetting() Setting {
	return Setting{
		Id:           s.Id,
		Value:        s.Value,
		SettingGroup: s.SettingGroup,
	}
}

func (s *Setting) ToResponse() SettingResponse {
	return SettingResponse{
		Id:            s.Id,
		Value:         s.Value,
		PreviousValue: s.PreviousValue,
		SettingGroup:  s.SettingGroup,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
