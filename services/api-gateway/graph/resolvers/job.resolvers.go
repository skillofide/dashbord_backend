package resolvers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/graphql-go/graphql"
)

type JobClients struct {
	// Add dependencies here if needed (e.g., config, logger, DB)
}

func NewJobClients() *JobClients {
	return &JobClients{}
}

// SearchJobs calls the Careerjet API to fetch jobs
func (c *JobClients) SearchJobs(p graphql.ResolveParams) (interface{}, error) {
	keywords, _ := p.Args["keywords"].(string)
	location, _ := p.Args["location"].(string)
	page, ok := p.Args["page"].(int)
	if !ok || page < 1 {
		page = 1
	}

	// NOTE: Careerjet requires an Affiliate ID to function properly.
	// We are using a dummy ID here. In production, register on careerjet.com
	affid := "test_affiliate_id" 
	userIP := "192.168.1.1" // In a real app, extract from p.Context via HTTP headers
	userAgent := "Skillofied/1.0" // In a real app, extract from p.Context

	baseURL := "http://public.api.careerjet.net/search"
	
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Careerjet API URL: %v", err)
	}

	q := u.Query()
	q.Set("affid", affid)
	q.Set("user_ip", userIP)
	q.Set("user_agent", userAgent)
	
	if keywords != "" {
		q.Set("keywords", keywords)
	}
	if location != "" {
		q.Set("location", location)
	}
	q.Set("page", fmt.Sprintf("%d", page))

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to Careerjet: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("careerjet API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Jobs []struct {
			URL         string `json:"url"`
			Title       string `json:"title"`
			Company     string `json:"company"`
			Locations   string `json:"locations"`
			Description string `json:"description"`
			Salary      string `json:"salary"`
			Date        string `json:"date"`
			Site        string `json:"site"`
		} `json:"jobs"`
		Pages int `json:"pages"`
		Total int `json:"hits"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Format response to match the GraphQL output types
	var jobs []map[string]interface{}
	for _, j := range result.Jobs {
		jobs = append(jobs, map[string]interface{}{
			"url":         j.URL,
			"title":       j.Title,
			"company":     j.Company,
			"locations":   j.Locations,
			"description": j.Description,
			"salary":      j.Salary,
			"date":        j.Date,
			"site":        j.Site,
		})
	}

	return map[string]interface{}{
		"jobs":  jobs,
		"total": result.Total,
		"pages": result.Pages,
	}, nil
}
