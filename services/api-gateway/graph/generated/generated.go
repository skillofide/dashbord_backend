// Package generated contains the programmatically built GraphQL schema for the API gateway.
// This replaces what gqlgen would normally generate — no codegen tool required.
package generated

import (
	"github.com/graphql-go/graphql"

	"github.com/skillofide/api-gateway/graph/resolvers"
)

// ─── Shared GraphQL Types ─────────────────────────────────────────────────────

var exampleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Example",
	Fields: graphql.Fields{
		"input":       &graphql.Field{Type: graphql.String},
		"output":      &graphql.Field{Type: graphql.String},
		"explanation": &graphql.Field{Type: graphql.String},
	},
})

var hintType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Hint",
	Fields: graphql.Fields{
		"order": &graphql.Field{Type: graphql.Int},
		"title": &graphql.Field{Type: graphql.String},
		"body":  &graphql.Field{Type: graphql.String},
	},
})

var starterCodesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StarterCodes",
	Fields: graphql.Fields{
		"javascript": &graphql.Field{Type: graphql.String},
		"python":     &graphql.Field{Type: graphql.String},
		"java":       &graphql.Field{Type: graphql.String},
		"cpp":        &graphql.Field{Type: graphql.String},
		"go":         &graphql.Field{Type: graphql.String},
	},
})

var problemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Problem",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.ID},
		"slug":         &graphql.Field{Type: graphql.String},
		"title":        &graphql.Field{Type: graphql.String},
		"difficulty":   &graphql.Field{Type: graphql.String},
		"topic":        &graphql.Field{Type: graphql.String},
		"xp":           &graphql.Field{Type: graphql.Int},
		"statement":    &graphql.Field{Type: graphql.String},
		"constraints":  &graphql.Field{Type: graphql.NewList(graphql.String)},
		"tags":         &graphql.Field{Type: graphql.NewList(graphql.String)},
		"examples":     &graphql.Field{Type: graphql.NewList(exampleType)},
		"hints":        &graphql.Field{Type: graphql.NewList(hintType)},
		"starterCodes": &graphql.Field{Type: starterCodesType},
		"setId":        &graphql.Field{Type: graphql.String},
		"userStatus":   &graphql.Field{Type: graphql.String},
	},
})

var practiceSetType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PracticeSet",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.ID},
		"title":         &graphql.Field{Type: graphql.String},
		"level":         &graphql.Field{Type: graphql.String},
		"levelColor":    &graphql.Field{Type: graphql.String},
		"bgColor":       &graphql.Field{Type: graphql.String},
		"totalProblems": &graphql.Field{Type: graphql.Int},
		"progress":      &graphql.Field{Type: graphql.Float},
	},
})

var listProblemsResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ListProblemsResult",
	Fields: graphql.Fields{
		"problems": &graphql.Field{Type: graphql.NewList(problemType)},
		"total":    &graphql.Field{Type: graphql.Int},
		"page":     &graphql.Field{Type: graphql.Int},
		"pageSize": &graphql.Field{Type: graphql.Int},
	},
})

var testResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "TestResult",
	Fields: graphql.Fields{
		"testCaseId":     &graphql.Field{Type: graphql.String},
		"input":          &graphql.Field{Type: graphql.String},
		"expectedOutput": &graphql.Field{Type: graphql.String},
		"actualOutput":   &graphql.Field{Type: graphql.String},
		"status":         &graphql.Field{Type: graphql.String},
		"executionMs":    &graphql.Field{Type: graphql.Int},
		"error":          &graphql.Field{Type: graphql.String},
	},
})

var submissionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Submission",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"userId":      &graphql.Field{Type: graphql.String},
		"problemId":   &graphql.Field{Type: graphql.String},
		"language":    &graphql.Field{Type: graphql.String},
		"status":      &graphql.Field{Type: graphql.String},
		"runtimeMs":   &graphql.Field{Type: graphql.Int},
		"memoryKb":    &graphql.Field{Type: graphql.Int},
		"testResults": &graphql.Field{Type: graphql.NewList(testResultType)},
		"submittedAt": &graphql.Field{Type: graphql.String},
		"completedAt": &graphql.Field{Type: graphql.String},
	},
})

var listSubmissionsResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ListSubmissionsResult",
	Fields: graphql.Fields{
		"submissions": &graphql.Field{Type: graphql.NewList(submissionType)},
		"total":       &graphql.Field{Type: graphql.Int},
		"page":        &graphql.Field{Type: graphql.Int},
		"pageSize":    &graphql.Field{Type: graphql.Int},
	},
})

var setProgressType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SetProgress",
	Fields: graphql.Fields{
		"setId":    &graphql.Field{Type: graphql.String},
		"title":    &graphql.Field{Type: graphql.String},
		"progress": &graphql.Field{Type: graphql.Float},
		"solved":   &graphql.Field{Type: graphql.Int},
		"total":    &graphql.Field{Type: graphql.Int},
	},
})

var userProgressType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserProgress",
	Fields: graphql.Fields{
		"userId":         &graphql.Field{Type: graphql.String},
		"totalSolved":    &graphql.Field{Type: graphql.Int},
		"totalAttempted": &graphql.Field{Type: graphql.Int},
		"easySolved":     &graphql.Field{Type: graphql.Int},
		"mediumSolved":   &graphql.Field{Type: graphql.Int},
		"hardSolved":     &graphql.Field{Type: graphql.Int},
		"currentStreak":  &graphql.Field{Type: graphql.Int},
		"longestStreak":  &graphql.Field{Type: graphql.Int},
		"totalXp":        &graphql.Field{Type: graphql.Int},
		"setProgress":    &graphql.Field{Type: graphql.NewList(setProgressType)},
	},
})

var runCodeResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RunCodeResult",
	Fields: graphql.Fields{
		"jobId":         &graphql.Field{Type: graphql.String},
		"overallStatus": &graphql.Field{Type: graphql.String},
		"testResults":   &graphql.Field{Type: graphql.NewList(testResultType)},
		"compileError":  &graphql.Field{Type: graphql.String},
		"runtimeMs":     &graphql.Field{Type: graphql.Int},
	},
})

var submitResultType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SubmitResult",
	Fields: graphql.Fields{
		"submissionId": &graphql.Field{Type: graphql.String},
	},
})

// ─── Schema Builder ───────────────────────────────────────────────────────────

// Clients bundles all service clients needed by the schema.
type Clients struct {
	Problems    *resolvers.ProblemClients
	Submissions *resolvers.SubmissionClients
	Progress    *resolvers.ProgressClients
}

// BuildSchema constructs the full GraphQL schema wiring resolvers to types.
func BuildSchema(clients *Clients) (graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listProblems": {
				Type: listProblemsResultType,
				Args: graphql.FieldConfigArgument{
					"setId":      {Type: graphql.String},
					"topic":      {Type: graphql.String},
					"difficulty": {Type: graphql.String},
					"page":       {Type: graphql.Int},
					"pageSize":   {Type: graphql.Int},
				},
				Resolve: clients.Problems.ListProblems,
			},
			"getProblem": {
				Type: problemType,
				Args: graphql.FieldConfigArgument{
					"id": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: clients.Problems.GetProblem,
			},
			"listPracticeSets": {
				Type:    graphql.NewList(practiceSetType),
				Resolve: clients.Problems.ListPracticeSets,
			},
			"getSubmission": {
				Type: submissionType,
				Args: graphql.FieldConfigArgument{
					"id": {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: clients.Submissions.GetSubmission,
			},
			"listSubmissions": {
				Type: listSubmissionsResultType,
				Args: graphql.FieldConfigArgument{
					"problemId": {Type: graphql.String},
					"page":      {Type: graphql.Int},
					"pageSize":  {Type: graphql.Int},
				},
				Resolve: clients.Submissions.ListSubmissions,
			},
			"getUserProgress": {
				Type:    userProgressType,
				Resolve: clients.Progress.GetUserProgress,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"submitCode": {
				Type: submitResultType,
				Args: graphql.FieldConfigArgument{
					"problemId": {Type: graphql.NewNonNull(graphql.String)},
					"language":  {Type: graphql.NewNonNull(graphql.String)},
					"code":      {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: clients.Submissions.SubmitCode,
			},
			"runCode": {
				Type: runCodeResultType,
				Args: graphql.FieldConfigArgument{
					"problemId": {Type: graphql.NewNonNull(graphql.String)},
					"language":  {Type: graphql.NewNonNull(graphql.String)},
					"code":      {Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: clients.Submissions.RunCode,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
