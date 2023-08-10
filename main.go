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
}

func getLLMProps(llm LLM) {
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

func promptModel(llm LLM) {
	fmt.Println("Code goes here")
}

func main() {
	fmt.Println("Code goes here")
}
