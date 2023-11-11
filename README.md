# GoLlama: Llama.cpp IPC Library

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


## Important Notes

As of August 21st 2023, llama.cpp no longer supports the GGML format. However, this doesn't affect GoLlama's performance since it runs on top of a llama.cpp instance. This means that as long as you can run llama.cpp on your machine you can use GoLlama.

In addition, GGUF is a new format introduced by the llama.cpp team on August 21st 2023. It is a replacement for GGML.

If you'd like to use models in the GGML format, just set llama.cpp on commit 220d9318647a8ce127dbf7c9de5400455f41e7d8 (keep in mind this might not work if your model is from August 21st 2023 onward) and run the setup as normal.

To use GGUF models, you can just switch to the latest commit and follow the setup instructions in llama.cpp's github repo. 

Remember: As long as you can get llama.cpp running, you can use GoLlama.

### Cloning the repo

GoLlama requires a running instance of Llama.cpp in order to communicate with any LLM. Run the following command to clone the repo alonside the llama.cpp submodule:

`git clone --recursive https://github.com/CenturySturgeon/gollama.git`

### Downloading a model

Additionally, an LLM is required for GoLlama to work. You can use any model you want as long as it is in ggml (using llama.cpp's commit 220d9318647a8ce127dbf7c9de5400455f41e7d8 or earlier) or gguf format. Even though you can download the models and build the gguf/ggml files from source, I'd recommend you go to Hugginface and check out user TheBloke, since he has put in the work of making many LLMs available in gguf/ggml format https://huggingface.co/TheBloke.

Once you've downloaded a model, don't forget to point your GoLlama's LLM instance to the path of the model. Since LLMs are very heavy files, even when quantized, I'd recommend you store all your models in a single directory and point the GoLlama LLM instance using an absoulute path. This way you don't have to download the same model multiple times during the development of multiple Gollama applications.

```
llm := gollama.LLM{Model: "ABSOLUTE_PATH_TO_YOUR_MODEL", ... }
```

## Docs

All implementations of the described methods are covered on the examples folder.

### PromptModel

PromptModel method orderly prompts the LLM with the provided prompts in the array, engaging in a sort of conversation. It returns an array with the respones of the LLM, each response matching with the index of its prompt.

```
llm.PromptModel([]string{"Your Prompt Here", "Another Prompt Here"})
```

### BufferPromptModel

Method that returns the model's response in real time as it is being built. A prompt passed as a string and a channel are required, the response tokens will be sent to the channel as they're being written.

```
llm.BufferPromptModel("Your Prompt Here", outputChannel)
```