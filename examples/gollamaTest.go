package main

import (
	"fmt"

	"github.com/CenturySturgeon/gollama"
)

// Types the given string to the terminal on the same line.
func typeToTerminal(text string) {
	for _, char := range text {
		fmt.Printf("%c", char)
	}
}

func main() {
	// Create a LLM instance
	llm := gollama.LLM{Model: "../llama.cpp/models/llama-2-13b-chat.ggmlv3.q4_0.bin", Llamacpp: "../llama.cpp", Ngl: 30}

	// =============================================== Simple Prompting ===============================================

	// Prompt the LLM and collects the response
	outputs, err := llm.PromptModel([]string{"Hi how are you ?"})
	if err != nil {
		fmt.Println("Error occured when executing propmt", err)
	}
	// Print the response
	fmt.Println(outputs)

	// =============================================== Runtime Prompting ===============================================

	// Create a channel to receive the LLM's response in runtime
	outputChan := make(chan string)
	// Create a go routine to ensure the the main program's execution isn't blocked, allowing it to receive and display characters as they are generated.
	go llm.BufferPromptModel("How are you doing?", outputChan)

	// Create a variable to store the whole response
	outputText := ""

	// Print empty line to differentiate from previous output
	fmt.Println()

	// Iterate over the channel's elements in a continuous loop, rather than stopping at the end of the channel
	for char := range outputChan {
		outputText += string(char)
		typeToTerminal(string(char))
	}

	// Print the whole response
	fmt.Printf("\nLLM output: %s\n", outputText)
}
