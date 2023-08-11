package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type LLM struct {
	model            string   // Path to the model.bin
	cuda_devices     []int    // Array of indices of the Cuda devices that will be used
	ctx_size         int      // Size of the prompt context
	temp             float32  // Temperature
	top_k            int      // Top-k sampling
	repeat_penalty   float32  // Penalize repeat sequence of tokens
	ngl              int      // Number of layers to store in VRAM
	max_tokens       int      // Max number of tokens for model response
	stop             []string // Array of generation-stopping strings
	instructionBlock string   // Instructions to format the model respone
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

func (llm *LLM) promptModel(prompts []string) {
	llm.llmDefaults()
	cmd, stdin, stdout, err := createPipes(llm)

	if err != nil {
		fmt.Println("Error creating pipes", err)
		return
	}

	// Start the llama.cpp llm communication process
	comErr := cmd.Start()
	if comErr != nil {
		fmt.Println("Error starting command:", comErr)
		return
	}

	// Array for the collection of outputs
	outputs := []string{}

	// Create a buffer for the stdout
	buf := make([]byte, 1024)

	// Create a counter for the amount of completed inputs
	counter := 0

	// Prompt all the inputs
	for i, input := range prompts {
		// When a prompt is first sent, it creates a \n> character automatically, so 'i' is incremented by 1 to reflect this
		i = i + 1

		// Add the instruction block to the input
		input = llm.instructionBlock + input

		// Input must contain an EOL for the LLM to correctly interpret the propmt's end
		if !strings.Contains(input, "\n") {
			input += "\n"
		}

		// Prompting the llm
		io.WriteString(stdin, input)

		output := ""
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading token:", err)
				}
				break
			}

			token := string(buf[:n])
			output = output + token

			if strings.Contains(token, "\n>") {
				counter += 1
				if counter > i {
					break
				}
			}
		}
		fmt.Println("Completed reading output.")
		outputs = append(outputs, strings.ReplaceAll(strings.ReplaceAll(output, "\n", ""), ">", ""))
	}
	// Print the outputs array to the screen
	fmt.Println(outputs)

	// Close the communication with the LLM
	closePipes(cmd, stdin, stdout)
}

func createPipes(llm *LLM) (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	cmd := exec.Command("./llama.cpp/main", "-m", llm.model, "--color", "--ctx_size", fmt.Sprint(llm.ctx_size), "-n", "-1", "-ins", "-b", "128", "--top_k", fmt.Sprint(llm.top_k), "--temp", fmt.Sprint(llm.temp), "--repeat_penalty", fmt.Sprint(llm.repeat_penalty), "--n-gpu-layers", fmt.Sprint(llm.ngl), "-t", "8")
	// Set the working directory if needed (for access to other directories)
	// cmd.Dir = ""

	// Create a writer for sending data to Python's stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return nil, nil, nil, err
	}

	// Create pipes for capturing Python's stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return nil, nil, nil, err
	}

	return cmd, stdin, stdout, nil
}

func closePipes(cmd *exec.Cmd, stdin io.WriteCloser, stdout io.ReadCloser) {

	fmt.Println("Closing stdin")
	// Close the stdin pipe to signal the end of input
	myerr := stdin.Close()

	if myerr != nil {
		fmt.Println("Error when closing the command :", myerr)
	}

	// Close the communication with the llm
	cmd.Process.Kill()
}

func main() {
	llm := LLM{model: "./models/llama-2-13b-chat.ggmlv3.q4_0.bin", ngl: 30}
	llm.promptModel([]string{"Hi, how are you ?", "How many pieces of bread do I need for a sandwich?"})
}
