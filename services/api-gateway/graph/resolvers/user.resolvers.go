package resolvers

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"go.uber.org/zap"

	"github.com/skillofide/api-gateway/middleware"
	userv1 "github.com/skillofide/proto/user/v1"
)

// UserClients holds the gRPC clients for user resolvers.
type UserClients struct {
	UserSvc userv1.UserServiceClient
	Log     *zap.Logger
}

// GetProfile handles the getProfile GraphQL query.
func (c *UserClients) GetProfile(p graphql.ResolveParams) (interface{}, error) {
	userID := middleware.UserIDFromContext(p.Context)
	if userID == "" {
		return nil, fmt.Errorf("authentication required")
	}

	resp, err := c.UserSvc.GetProfile(p.Context, &userv1.GetProfileRequest{
		UserID: userID,
	})
	if err != nil {
		c.Log.Error("get profile resolver failed", zap.Error(err))
		return nil, fmt.Errorf("failed to get profile: %v", err)
	}

	prof := resp.Profile
	if prof == nil {
		return map[string]interface{}{
			"userId":                userID,
			"gender":                "",
			"dob":                   "",
			"whatsapp":              "",
			"phone":                 "",
			"experience":            "",
			"workExperience":        "",
			"careerGap":             "",
			"currentState":          "",
			"currentCity":           "",
			"preferredLocations":    []string{},
			"githubLink":            "",
			"linkedinLink":          "",
			"isWorkingProfessional": false,
			"resumeName":            "",
			"edu10SchoolName":       "",
			"edu10YearOfPassout":    "",
			"edu10MarksPercent":     "",
			"edu12SchoolName":       "",
			"edu12YearOfPassout":    "",
			"edu12MarksPercent":     "",
			"ugUniversityRollNo":    "",
			"ugCollegeName":         "",
			"ugCourseName":          "",
			"ugBranch":              "",
			"ugYearOfPassout":       "",
			"ugMarksPercent":        "",
			"ugCgpa":                "",
			"ugActiveBacklogs":      "",
			"pgHasCertificate":      false,
		}, nil
	}

	return map[string]interface{}{
		"userId":                prof.UserID,
		"gender":                prof.Gender,
		"dob":                   prof.Dob,
		"whatsapp":              prof.Whatsapp,
		"phone":                 prof.Phone,
		"experience":            prof.Experience,
		"workExperience":        prof.WorkExperience,
		"careerGap":             prof.CareerGap,
		"currentState":          prof.CurrentState,
		"currentCity":           prof.CurrentCity,
		"preferredLocations":    prof.PreferredLocations,
		"githubLink":            prof.GithubLink,
		"linkedinLink":          prof.LinkedinLink,
		"isWorkingProfessional": prof.IsWorkingProfessional,
		"resumeName":            prof.ResumeName,
		"edu10SchoolName":       prof.Edu10SchoolName,
		"edu10YearOfPassout":    prof.Edu10YearOfPassout,
		"edu10MarksPercent":     prof.Edu10MarksPercent,
		"edu12SchoolName":       prof.Edu12SchoolName,
		"edu12YearOfPassout":    prof.Edu12YearOfPassout,
		"edu12MarksPercent":     prof.Edu12MarksPercent,
		"ugUniversityRollNo":    prof.UGUniversityRollNo,
		"ugCollegeName":         prof.UGCollegeName,
		"ugCourseName":          prof.UGCourseName,
		"ugBranch":              prof.UGBranch,
		"ugYearOfPassout":       prof.UGYearOfPassout,
		"ugMarksPercent":        prof.UGMarksPercent,
		"ugCgpa":                prof.UGCGPA,
		"ugActiveBacklogs":      prof.UGActiveBacklogs,
		"pgHasCertificate":      prof.PGHasCertificate,
	}, nil
}

// UpsertProfile handles the upsertProfile GraphQL mutation.
func (c *UserClients) UpsertProfile(p graphql.ResolveParams) (interface{}, error) {
	userID := middleware.UserIDFromContext(p.Context)
	if userID == "" {
		return nil, fmt.Errorf("authentication required")
	}

	input, ok := p.Args["profile"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid profile input")
	}

	strVal := func(key string) string {
		v, _ := input[key].(string)
		return v
	}

	boolVal := func(key string) bool {
		v, _ := input[key].(bool)
		return v
	}

	strSliceVal := func(key string) []string {
		raw, ok := input[key].([]interface{})
		if !ok {
			return []string{}
		}
		res := make([]string, 0, len(raw))
		for _, item := range raw {
			if s, ok := item.(string); ok {
				res = append(res, s)
			}
		}
		return res
	}

	profile := &userv1.UserProfile{
		UserID:                userID,
		Gender:                strVal("gender"),
		Dob:                   strVal("dob"),
		Whatsapp:              strVal("whatsapp"),
		Phone:                 strVal("phone"),
		Experience:            strVal("experience"),
		WorkExperience:        strVal("workExperience"),
		CareerGap:             strVal("careerGap"),
		CurrentState:          strVal("currentState"),
		CurrentCity:           strVal("currentCity"),
		PreferredLocations:    strSliceVal("preferredLocations"),
		GithubLink:            strVal("githubLink"),
		LinkedinLink:          strVal("linkedinLink"),
		IsWorkingProfessional: boolVal("isWorkingProfessional"),
		ResumeName:            strVal("resumeName"),
		Edu10SchoolName:       strVal("edu10SchoolName"),
		Edu10YearOfPassout:    strVal("edu10YearOfPassout"),
		Edu10MarksPercent:     strVal("edu10MarksPercent"),
		Edu12SchoolName:       strVal("edu12SchoolName"),
		Edu12YearOfPassout:    strVal("edu12YearOfPassout"),
		Edu12MarksPercent:     strVal("edu12MarksPercent"),
		UGUniversityRollNo:    strVal("ugUniversityRollNo"),
		UGCollegeName:         strVal("ugCollegeName"),
		UGCourseName:          strVal("ugCourseName"),
		UGBranch:              strVal("ugBranch"),
		UGYearOfPassout:       strVal("ugYearOfPassout"),
		UGMarksPercent:        strVal("ugMarksPercent"),
		UGCGPA:                strVal("ugCgpa"),
		UGActiveBacklogs:      strVal("ugActiveBacklogs"),
		PGHasCertificate:      boolVal("pgHasCertificate"),
	}

	resp, err := c.UserSvc.UpsertProfile(p.Context, &userv1.UpsertProfileRequest{
		Profile: profile,
	})
	if err != nil {
		c.Log.Error("upsert profile resolver failed", zap.Error(err))
		return nil, fmt.Errorf("failed to save profile: %v", err)
	}

	return map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	}, nil
}
