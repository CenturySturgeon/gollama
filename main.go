package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

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

func (llm *LLM) promptModel(prompts []string) {
	llm.llmDefaults()
	cmd, stdin, stdout, err := createPipes(buildCommand(llm))

	if err != nil {
		fmt.Println("Error creating pipes", err)
		return
	}

	// Read and collect outputs
	outputs := []string{}
	buf := make([]byte, 1024)

	for _, input := range prompts {

		// Input must contain an EOL for the LLM to correctly interpret the propmt's end
		if !strings.Contains(input, "\n") {
			input += "\n"
		}
		// Sending input
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

			if strings.Count(output, ">") >= 2 {
				break
			}
		}

		fmt.Println("Completed reading output.")
		outputs = append(outputs, strings.ReplaceAll(strings.ReplaceAll(output, "\n", ""), ">", ""))
	}

	// Print collected outputs
	fmt.Println("Outputs: ", outputs)

	// Close the communication with the LLM
	closePipes(cmd, stdin, stdout)
}

func buildCommand(llm *LLM) string {
	return fmt.Sprintf("./llama.cpp/main -m %s --color   --ctx_size %d   -n -1   -ins -b 128   --top_k %d   --temp %.1f   --repeat_penalty %.1f   --n-gpu-layers %d   -t 8", llm.model, llm.ctx_size, llm.top_k, llm.temp, llm.repeat_penalty, llm.ngl)
}

func createPipes(command string) (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	cmd := exec.Command(command)

	// Create pipes for stdin and stdout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return nil, nil, nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return nil, nil, nil, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return nil, nil, nil, err
	}

	return cmd, stdin, stdout, nil
}

func closePipes(cmd *exec.Cmd, stdin io.WriteCloser, stdout io.ReadCloser) {

	fmt.Println("Closing stdin")
	// Close stdin explicitly
	myerr := stdin.Close()

	if myerr != nil {
		fmt.Println("Error when closing the command :", myerr)
	}

	// Kill the process
	cmd.Process.Kill()
}

func main() {
	llm := LLM{model: "~/Ai/models/llama-2-13b-chat.ggmlv3.q4_0.bin", ngl: 30}
	llm.promptModel([]string{"Hi, how are you ?"})
}
