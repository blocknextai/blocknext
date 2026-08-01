package organizeemails

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/helpers"
)

var (
	ErrUnsupportedCategory = apperror.Internal("unsupported category")
)

type GmailOrganizeEmailsExecutorInput struct {
	Keywords string                      `schema:"keywords"`
	Actions  []GmailOrganizeEmailsAction `schema:"actions"`
}

type GmailOrganizeEmailsAction struct {
	Type      string `schema:"type"`
	LabelName string `schema:"labelName"`
	Category  string `schema:"category"`
	Star      bool   `schema:"star"`
	Read      bool   `schema:"read"`
	Trash     bool   `schema:"trash"`
}

type GmailOrganizeEmailsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GmailOrganizeEmailsExecutorInput]
}

func NewGmailOrganizeEmailsExecutor(
	nodeID string,
	validator *jsonschema.Validator[GmailOrganizeEmailsExecutorInput],
) *GmailOrganizeEmailsExecutor {
	return &GmailOrganizeEmailsExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *GmailOrganizeEmailsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "gmail_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			emailIDs, err := e.findEmailsByKeywords(client, input.Keywords)
			if err != nil {
				return nil, err
			}

			if len(emailIDs) == 0 {
				return []map[string]any{
					{
						"status": true,
						"result": map[string]any{
							"processedEmails": 0,
							"appliedActions":  []string{},
							"scope":           getScopeInfo(input.Keywords),
						},
					},
				}, nil
			}

			processedCount, appliedActions, err := e.applyMultipleActions(client, emailIDs, input.Actions)
			if err != nil {
				return nil, err
			}

			results = append(results, map[string]any{
				"status": true,
				"result": map[string]any{
					"processedEmails": processedCount,
					"appliedActions":  appliedActions,
					"scope":           getScopeInfo(input.Keywords),
				},
			})
		}

		return results, nil
	}
}

type GmailLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GmailLabelsResponse struct {
	Labels []GmailLabel `json:"labels"`
}

type GmailErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GmailOrganizeEmailsExecutor) ensureLabelExists(client *httpclient.Client, labelName string) (string, error) {
	if labelID, found, err := e.findExistingLabel(client, labelName); err != nil {
		return "", err
	} else if found {
		return labelID, nil
	}

	return e.createNewLabel(client, labelName)
}

func (e *GmailOrganizeEmailsExecutor) findExistingLabel(client *httpclient.Client, labelName string) (string, bool, error) {
	var labelsResponse GmailLabelsResponse
	var errorResponse GmailErrorResponse

	response, err := client.Get("/users/me/labels").Do(&labelsResponse, &errorResponse)
	if err != nil {
		return "", false, err
	}

	if !response.IsSuccess() {
		return "", false, apperror.Internal(errorResponse.Error.Message)
	}

	for _, label := range labelsResponse.Labels {
		if strings.EqualFold(label.Name, labelName) {
			return label.ID, true, nil
		}
	}

	return "", false, nil
}

func (e *GmailOrganizeEmailsExecutor) createNewLabel(client *httpclient.Client, labelName string) (string, error) {
	var createResponse GmailLabel
	var errorResponse GmailErrorResponse

	response, err := client.Post("/users/me/labels").
		JSONContentType().
		Body(map[string]any{"name": labelName}).
		Do(&createResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	return createResponse.ID, nil
}

type GmailSearchResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	NextPageToken string `json:"nextPageToken"`
}

func (e *GmailOrganizeEmailsExecutor) findEmailsByKeywords(client *httpclient.Client, keywords string) ([]string, error) {
	if len(keywords) == 0 {
		return e.searchEmails(client, "")
	}

	emailIDSet := make(map[string]struct{})
	ids, err := e.searchEmails(client, keywords)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		emailIDSet[id] = struct{}{}
	}

	allEmailIDs := make([]string, 0, len(emailIDSet))
	for id := range emailIDSet {
		allEmailIDs = append(allEmailIDs, id)
	}

	return allEmailIDs, nil
}

func (e *GmailOrganizeEmailsExecutor) searchEmails(client *httpclient.Client, keywords string) ([]string, error) {
	emailIDs := make([]string, 0, 500)
	pageToken := ""

	for {
		ids, nextToken, err := e.fetchEmailPage(client, keywords, pageToken)
		if err != nil {
			return nil, err
		}

		emailIDs = append(emailIDs, ids...)

		if nextToken == "" || (keywords == "" && len(emailIDs) >= 10000) {
			break
		}
		pageToken = nextToken
	}

	return emailIDs, nil
}

func (e *GmailOrganizeEmailsExecutor) fetchEmailPage(client *httpclient.Client, keywords, pageToken string) ([]string, string, error) {
	var searchResponse GmailSearchResponse
	var errorResponse GmailErrorResponse

	request := client.Get("/users/me/messages").
		QueryParam("maxResults", "500")

	if keywords != "" {
		request = request.QueryParam("q", keywords)
	}

	if pageToken != "" {
		request = request.QueryParam("pageToken", pageToken)
	}

	response, err := request.Do(&searchResponse, &errorResponse)
	if err != nil {
		return nil, "", err
	}

	if !response.IsSuccess() {
		return nil, "", apperror.Internal(errorResponse.Error.Message)
	}

	ids := make([]string, 0, len(searchResponse.Messages))
	for _, message := range searchResponse.Messages {
		ids = append(ids, message.ID)
	}

	return ids, searchResponse.NextPageToken, nil
}

func (e *GmailOrganizeEmailsExecutor) modifyEmailLabels(client *httpclient.Client, emailIDs []string, addLabels, removeLabels []string) (int, error) {
	if len(emailIDs) == 0 || (len(addLabels) == 0 && len(removeLabels) == 0) {
		return 0, nil
	}

	body := e.buildModifyBody(addLabels, removeLabels)
	successCount := 0

	for _, emailID := range emailIDs {
		if e.modifyEmailLabel(client, emailID, body) {
			successCount++
		}
	}

	return successCount, nil
}

func (e *GmailOrganizeEmailsExecutor) buildModifyBody(addLabels, removeLabels []string) map[string]any {
	body := make(map[string]any, 2)
	if len(addLabels) > 0 {
		body["addLabelIds"] = addLabels
	}
	if len(removeLabels) > 0 {
		body["removeLabelIds"] = removeLabels
	}
	return body
}

func (e *GmailOrganizeEmailsExecutor) modifyEmailLabel(client *httpclient.Client, emailID string, body map[string]any) bool {
	var modifyResponse map[string]any
	var errorResponse GmailErrorResponse

	response, err := client.Post("/users/me/messages/"+emailID+"/modify").
		JSONContentType().
		Body(body).
		Do(&modifyResponse, &errorResponse)

	return err == nil && response.IsSuccess()
}

func (e *GmailOrganizeEmailsExecutor) applyMultipleActions(client *httpclient.Client, emailIDs []string, actions []GmailOrganizeEmailsAction) (int, []string, error) {
	addLabels := make([]string, 0, len(actions))
	removeLabels := make([]string, 0, len(actions))
	appliedActions := make([]string, 0, len(actions))

	for _, action := range actions {
		if err := e.processAction(client, action, &addLabels, &removeLabels, &appliedActions); err != nil {
			return 0, nil, err
		}
	}

	if len(addLabels) == 0 && len(removeLabels) == 0 {
		return 0, appliedActions, nil
	}

	processedCount, err := e.modifyEmailLabels(client, emailIDs, addLabels, removeLabels)
	return processedCount, appliedActions, err
}

func (e *GmailOrganizeEmailsExecutor) processAction(client *httpclient.Client, action GmailOrganizeEmailsAction, addLabels, removeLabels, appliedActions *[]string) error {
	switch action.Type {
	case "label":
		return e.processLabelAction(client, action, addLabels, appliedActions)
	case "star":
		e.processStarAction(action, addLabels, removeLabels, appliedActions)
	case "read":
		e.processReadAction(action, addLabels, removeLabels, appliedActions)
	case "category":
		return e.processCategoryAction(action, addLabels, appliedActions)
	case "trash":
		e.processTrashAction(action, addLabels, appliedActions)
	}
	return nil
}

func (e *GmailOrganizeEmailsExecutor) processLabelAction(client *httpclient.Client, action GmailOrganizeEmailsAction, addLabels, appliedActions *[]string) error {
	if action.LabelName == "" {
		return nil
	}
	labelID, err := e.ensureLabelExists(client, action.LabelName)
	if err != nil {
		return err
	}
	*addLabels = append(*addLabels, labelID)
	*appliedActions = append(*appliedActions, "label:"+action.LabelName)
	return nil
}

func (e *GmailOrganizeEmailsExecutor) processStarAction(action GmailOrganizeEmailsAction, addLabels, removeLabels, appliedActions *[]string) {
	if action.Star {
		*addLabels = append(*addLabels, "STARRED")
		*appliedActions = append(*appliedActions, "star")
	} else {
		*removeLabels = append(*removeLabels, "STARRED")
		*appliedActions = append(*appliedActions, "unstar")
	}
}

func (e *GmailOrganizeEmailsExecutor) processReadAction(action GmailOrganizeEmailsAction, addLabels, removeLabels, appliedActions *[]string) {
	if action.Read {
		*removeLabels = append(*removeLabels, "UNREAD")
		*appliedActions = append(*appliedActions, "markAsRead")
	} else {
		*addLabels = append(*addLabels, "UNREAD")
		*appliedActions = append(*appliedActions, "markAsUnread")
	}
}

func (e *GmailOrganizeEmailsExecutor) processCategoryAction(action GmailOrganizeEmailsAction, addLabels, appliedActions *[]string) error {
	if action.Category == "" {
		return nil
	}
	categoryLabel, exists := categoryLabels[action.Category]
	if !exists {
		return ErrUnsupportedCategory
	}
	*addLabels = append(*addLabels, categoryLabel)
	*appliedActions = append(*appliedActions, "category:"+action.Category)
	return nil
}

func (e *GmailOrganizeEmailsExecutor) processTrashAction(action GmailOrganizeEmailsAction, addLabels, appliedActions *[]string) {
	if action.Trash {
		*addLabels = append(*addLabels, "TRASH")
		*appliedActions = append(*appliedActions, "moved to trash")
	}
}

var categoryLabels = map[string]string{
	"primary":    "CATEGORY_PERSONAL",
	"social":     "CATEGORY_SOCIAL",
	"promotions": "CATEGORY_PROMOTIONS",
	"updates":    "CATEGORY_UPDATES",
	"forums":     "CATEGORY_FORUMS",
}

func getScopeInfo(keywords string) any {
	if len(keywords) == 0 {
		return "all_emails"
	}
	return keywords
}
