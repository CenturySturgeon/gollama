package main

import "fmt"

type LLM struct {
	model          string
	cuda_devices   []int
	ctx_size       int
	temp           float32
	top_k          int
	repeat_penalty float32
	ngl            int
	max_tokens     int
	stop           []string
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
