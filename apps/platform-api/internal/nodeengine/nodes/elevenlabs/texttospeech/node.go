package texttospeech

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type ElevenlabsTextToSpeechNode struct {
	nodes.Node
}

func NewElevenlabsTextToSpeechNode(nodeID string) *ElevenlabsTextToSpeechNode {
	return &ElevenlabsTextToSpeechNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "ElevenLabs Text To Speech",
			Description: "Convert text to speech using ElevenLabs.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Audio"},
			SubCategories: []string{"ElevenLabs"},
			Tags: []string{
				"ai",
				"audio",
				"voice",
				"speech",
				"tts",
				"generation",
			},
			SupportedCredentials: []string{
				"elevenlabs_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Text to convert to speech.",
					},
					"voiceId": {
						Type:        "string",
						Title:       "Voice ID",
						Description: "ElevenLabs voice identifier to use for synthesis.",
						Enum: []any{
							"21m00Tcm4TlvDq8ikWAM",
							"29vD33N1CtxCmqQRPOHJ",
							"2EiwWnXFnvU5JabPnv8n",
							"5Q0t7uMcjvnagumLfvZi",
							"9BWtsMINqrJLrRacOk9x",
							"AZnzlk1XvdvUeBnXmlld",
							"CYw3kZ02Hs0563khs1Fj",
							"CwhRBWXzGAHq8TQ4Fs17",
							"D38z5RcWu1voky8WS1ja",
							"EXAVITQu4vr4xnSDxMaL",
							"ErXwobaYiN019PkySvjV",
							"FGY2WhTYpPnrIDTdsKH5",
							"GBv7mTt0atIp3Br8iCZE",
							"IKne3meq5aSn9XLyUdCD",
							"JBFqnCBsd6RMkjVDRZzb",
							"LcfcDJNUP1GQjkzn1xUU",
							"MF3mGyEYCl7XYWbV9V6O",
							"N2lVS1w4EtoT3dr4eOWO",
							"ODq5zmih8GrVes37Dizd",
							"SAz9YHcvj6GT2YYXdXww",
							"SOYHLrjzK2X1ezoPC6cr",
							"TX3LPaxmHKxFdv7VOQHJ",
							"ThT5KcBeYPX3keUQqHPh",
							"TxGEqnHWrfWFTfGW9XjX",
							"VR6AewLTigWG4xSOukaG",
							"XB0fDUnXU5powFXDhCwa",
							"Xb7hH8MSUJpSbSDYk0k2",
							"XrExE9yKIg1WjnnlVkGX",
							"Yko7PKHZNXotIFUBG7I9",
							"ZQe5CZNOzWyzPSCn5a3c",
							"Zlb1dXrM653N07WRdFW3",
							"bIHbv24MWmeRgasZH58o",
							"bVMeCyTHy58xNoL34h3p",
							"cgSgspJ2msm6clMCkdW9",
							"cjVigY5qzO86Huf0OWal",
							"flq6f7yk4E4fJM5XTYuZ",
							"g5CIjZEefAph4nQFvHAz",
							"iP95p4xoKVk53GoZ742B",
							"jBpfuIE2acCO8z3wKNLl",
							"jsCqWAovK2LkecY7zXl4",
							"knrPHWnBmmDHMoiMeP3l",
							"nPczCjzI2devNBz1zQrb",
							"oWAxZDx7w5VEj9dCyTzz",
							"onwK4e9ZLuTAKqWW03F9",
							"pFZP5JQG7iQjIQuC4Bku",
							"pMsXgVXv3BLzUgSXRplE",
							"pNInz6obpgDQGcFmaJgB",
							"piTKgcLEGmPE4e6mEKli",
							"pqHfZKP75CvOlQylNhV4",
							"t0jbNlBVZ17f02VDIeMI",
							"wViXBPUzp2ZZixB1xQuM",
							"yoZ06aMxZJJ28mfd3POQ",
							"z9fAnlkpzviPz146aGWa",
							"zcAOhNBS3c14rBihAFp1",
							"zrHiDhphv9ZnVXBqCLjz",
						},
						Default: json.RawMessage(`"JBFqnCBsd6RMkjVDRZzb"`),
					},
					"modelId": {
						Type:        "string",
						Title:       "Model ID",
						Description: "ElevenLabs synthesis model identifier.",
						Enum: []any{
							"eleven_v3",
							"eleven_ttv_v3",
							"eleven_multilingual_v2",
							"eleven_flash_v2_5",
							"eleven_flash_v2",
							"eleven_english_sts_v2",
							"eleven_multilingual_sts_v2",
						},
						Default: json.RawMessage(`"eleven_flash_v2_5"`),
					},
					"outputFormat": {
						Type:        "string",
						Title:       "Output Format",
						Description: "Audio output format identifier.",
						Enum:        []any{"mp3_44100_128"},
						Default:     json.RawMessage(`"mp3_44100_128"`),
					},
				},
				Required: []string{
					"text",
					"voiceId",
					"modelId",
					"outputFormat",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"audio": {
							Type:        "string",
							Description: "URL of the generated audio file.",
							Format:      "uri",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:    true,
				Destructive: new(false),
			},
		},
	}
}
