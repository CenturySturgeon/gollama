# GoLlama: Llama.cpp IPC Library

=============================================

GoLLama is a lightweight IPC library for developing LLM applications using Go and llama.cpp. It provides a simple and intuitive way to interact with the LLM runtime on top of llama.cpp using stdin/stdout.

### Usage

To use Go-Llama, you can import the `gollama` package and create an instance of the `LLM` struct. You can then use the `PromptModel` method to send and receive data between your application and any LLM supported by llama.cpp. Here is an example:

```

package main

import (
"fmt"
"github.com/CenturySturgeon/gollama"
)

func  main() {
	llm := gollama.LLM{Model: "../llama.cpp/models/llama-2-13b-chat.ggmlv3.q4_0.bin", Llamacpp: "../llama.cpp", Ngl: 30}

	outputs, err := llm.PromptModel([]string{"Hi how are you ?"})

	if err != nil {
		fmt.Println("Error occured on propmt: ", err)
	}
	fmt.Println(outputs)
}

```

This example demonstrates how to interact with a llama-2-13b-chat LLM, running on top of llama.cpp, from Go using the `PromptModel` method. It prompts the LLM with the question "Hi, how are you?" and then reads the response back into the `outputs` variable to print it afterwards.

  

### Cloning the repo

GoLlama needs a running instance of Llama.cpp in order to communicate with the LLM of your choice. Run the following command to clone the repo alonside the llama.cpp submodule:
`git clone --recursive https://github.com/CenturySturgeon/gollama.git`