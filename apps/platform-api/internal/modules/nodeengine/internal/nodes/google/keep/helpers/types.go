package helpers

import (
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type TextContent struct {
	Text string `json:"text"`
}

type ListItem struct {
	Text    TextContent `json:"text"`
	Checked bool        `json:"checked"`
}

type ListContent struct {
	ListItems []ListItem `json:"listItems"`
}

type Section struct {
	Text *TextContent `json:"text,omitempty"`
	List *ListContent `json:"list,omitempty"`
}

func (s Section) ToMap() map[string]any {
	body := make(map[string]any)
	if s.Text != nil {
		body["text"] = map[string]any{
			"text": s.Text.Text,
		}
	}
	if s.List != nil {
		items := make([]map[string]any, 0, len(s.List.ListItems))
		for _, item := range s.List.ListItems {
			items = append(items, map[string]any{
				"text": map[string]any{
					"text": item.Text.Text,
				},
				"checked": item.Checked,
			})
		}
		body["list"] = map[string]any{
			"listItems": items,
		}
	}
	return body
}

type Note struct {
	Name       string  `json:"name"`
	Title      string  `json:"title"`
	CreateTime string  `json:"createTime"`
	UpdateTime string  `json:"updateTime"`
	TrashTime  string  `json:"trashTime,omitempty"`
	Trashed    bool    `json:"trashed"`
	Body       Section `json:"body"`
}

func (n Note) ToMap() map[string]any {
	return map[string]any{
		"name":       n.Name,
		"title":      n.Title,
		"body":       n.Body.ToMap(),
		"createTime": n.CreateTime,
		"updateTime": n.UpdateTime,
		"trashed":    n.Trashed,
	}
}

type ListNotesResponse struct {
	Notes         []Note `json:"notes"`
	NextPageToken string `json:"nextPageToken"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func NoteSchema() *gjs.Schema {
	return &gjs.Schema{
		Type: "object",
		Properties: map[string]*gjs.Schema{
			"name": {
				Type:        "string",
				Description: "Resource name of the note (e.g. notes/abc123).",
			},
			"title": {
				Type:        "string",
				Description: "Title of the note.",
			},
			"body": {
				Type:        "object",
				Description: "Content of the note, holding either text or list content.",
				Properties: map[string]*gjs.Schema{
					"text": {
						Type:        "object",
						Description: "Text content of the note.",
						Properties: map[string]*gjs.Schema{
							"text": {
								Type:        "string",
								Description: "The text of the note.",
							},
						},
					},
					"list": {
						Type:        "object",
						Description: "List content of the note.",
						Properties: map[string]*gjs.Schema{
							"listItems": {
								Type:        "array",
								Description: "The items in the list.",
								Items: &gjs.Schema{
									Type: "object",
									Properties: map[string]*gjs.Schema{
										"text": {
											Type:        "object",
											Description: "Text content of the list item.",
											Properties: map[string]*gjs.Schema{
												"text": {
													Type:        "string",
													Description: "The text of the list item.",
												},
											},
										},
										"checked": {
											Type:        "boolean",
											Description: "Whether the list item has been checked off.",
										},
									},
								},
							},
						},
					},
				},
			},
			"createTime": {
				Type:        "string",
				Description: "Time the note was created.",
			},
			"updateTime": {
				Type:        "string",
				Description: "Time the note was last updated.",
			},
			"trashed": {
				Type:        "boolean",
				Description: "Whether the note is in the trash.",
			},
		},
	}
}
