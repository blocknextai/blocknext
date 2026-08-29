package prompts

import _ "embed"

//go:embed function_calling_system_instruction.md
var FunctionCallingSystemInstruction string

//go:embed workflow_generation_system_instruction.md
var WorkflowGenerationSystemInstruction string
