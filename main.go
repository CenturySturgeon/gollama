package main

import "fmt"

type LLM struct {
	model          string   // Path to the model.bin
	cuda_devices   []int    // Array of indices of the Cuda devices that will be used
	ctx_size       int      // Size of the prompt context
	temp           float32  // Temperature
	top_k          int      // Top-k sampling
	repeat_penalty float32  // Penalize repeat sequence of tokens
	ngl            int      // Number of layers to store in VRAM
	max_tokens     int      // Max number of tokens for model response
	stop           []string // Array of generation-stopping strings
	command        string   // Command to execute llama.cpp and therefore the model
}

func (llm *LLM) getLLMProps() {
	fmt.Println("Model Path: ", llm.model)
	fmt.Println("Indexes of Cuda devices to use: ", llm.cuda_devices)
	fmt.Println("Size of the prompt context: ", llm.ctx_size)
	fmt.Println("Temperature: ", llm.temp)
	fmt.Println("Top-k sampling: ", llm.top_k)
	fmt.Println("Penalize repeat sequence of tokens: ", llm.repeat_penalty)
	fmt.Println("Number of layers to store in VRAM: ", llm.ngl)
	fmt.Println("Max number of tokens for model response: ", llm.max_tokens)
	fmt.Println("List of generation-stopping strings: ", llm.stop)
}

func (llm *LLM) llmDefaults() {
	if llm.model == "" {
		llm.model = "./llama.cpp/models/ggml-vocab.bin"
	}
	if llm.cuda_devices == nil {
		llm.cuda_devices = []int{0}
	}
	if llm.ctx_size == 0 {
		llm.ctx_size = 2048
	}
	if llm.temp == 0 {
		llm.temp = 0.2
	}
	if llm.top_k == 0 {
		llm.top_k = 10000
	}
	if llm.repeat_penalty == 0 {
		llm.repeat_penalty = 1.1
	}
	if llm.max_tokens == 0 {
		llm.max_tokens = 1000
	}
}

func (llm *LLM) promptModel(prompt string) {
	llm.llmDefaults()
	llm.getLLMProps()
}

func main() {
	llm := LLM{}
	llm.promptModel("Hi")
}
