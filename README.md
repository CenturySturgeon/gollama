# GoLlama: Llama.cpp IPC Library

=======================================================================

GoLLama is a lightweight inter-process communication library for developing LLM applications using Go and llama.cpp. It provides a simple and intuitive way to interact with the LLM runtime on top of llama.cpp using stdin/stdout.

![Diagram](https://github.com/CenturySturgeon/CenturySturgeon.github.io/blob/main/Images/GoLlama.svg)

### Usage

To use GoLlama, you can import the `gollama` package and create an instance of the `LLM` struct. You can then use the `PromptModel` method to send and receive data between your application and any LLM supported by llama.cpp. Here is an example:

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

This example demonstrates how to use Go to interact with a llama-2-13b-chat LLM instance, which is running on top of llama.cpp,  using GoLlama's `PromptModel` method. It prompts the LLM with the question "Hi, how are you?" and then reads the response back into the `outputs` variable to print it afterwards.


### Important Notes

As of August 21st 2023, llama.cpp no longer supports the GGML format. However, this doesn't affet GoLlama's performance since it runs on top of a llama.cpp running instance. This means that as long as you can run llama.cpp on your machine you can use gollama.

GGUF is a new format introduced by the llama.cpp team on August 21st 2023. It is a replacement for GGML, which is no longer supported by llama.cpp.

If you'd like to use models that use the GGML format just set llama.cpp on commit 220d9318647a8ce127dbf7c9de5400455f41e7d8 (keep in mind this probably won't work if your model is from August 21st 2023 onward) and run the setup as normal.

To use GGUF models, you can just switch to the latest commit and follow the setup instructions in llama.cpp's github repo. 

Remember: As long as you can get llama.cpp running, you can use GoLlama.

### Cloning the repo

GoLlama requires a running instance of Llama.cpp in order to communicate with any LLM. Run the following command to clone the repo alonside the llama.cpp submodule:

`git clone --recursive https://github.com/CenturySturgeon/gollama.git`