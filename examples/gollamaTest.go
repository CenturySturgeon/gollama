package main

import (
	"fmt"

	"github.com/CenturySturgeon/gollama"
)

func main() {
	llm := gollama.LLM{Model: "../llama.cpp/models/llama-2-13b-chat.ggmlv3.q4_0.bin", Llamacpp: "../llama.cpp", Ngl: 30}
	outputs, err := llm.PromptModel([]string{"Hi how are you ?"})
	if err != nil {
		fmt.Println("Error occured when executing propmt", err)
	}
	fmt.Println(outputs)
}
