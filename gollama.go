package gollama

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type LLM struct {
	Model            string   // Path to the model.bin
	Llamacpp         string   // Path to the llama.cpp folder
	CudaDevices      []int    // Array of indices of the Cuda devices that will be used
	CtxSize          int      // Size of the prompt context
	Temp             float32  // Temperature
	TopK             int      // Top-k sampling
	RepeatPenalty    float32  // Penalize repeat sequence of tokens
	Ngl              int      // Number of layers to store in VRAM
	MaxTokens        int      // Max number of tokens for model response
	Stop             []string // Array of generation-stopping strings
	InstructionBlock string   // Instructions to format the model response
}

func (llm *LLM) GetLLMProps() {
	fmt.Println("Model Path: ", llm.Model)
	fmt.Println("Llama.cpp Path: ", llm.Llamacpp)
	fmt.Println("Indexes of Cuda devices to use: ", llm.CudaDevices)
	fmt.Println("Size of the prompt context: ", llm.CtxSize)
	fmt.Println("Temperature: ", llm.Temp)
	fmt.Println("Top-k sampling: ", llm.TopK)
	fmt.Println("Penalize repeat sequence of tokens: ", llm.RepeatPenalty)
	fmt.Println("Number of layers to store in VRAM: ", llm.Ngl)
	fmt.Println("Max number of tokens for model response: ", llm.MaxTokens)
	fmt.Println("List of generation-stopping strings: ", llm.Stop)
}

func (llm *LLM) llmDefaults() {
	if llm.Model == "" {
		llm.Model = "./llama.cpp/models/ggml-vocab.bin"
	}
	if llm.Llamacpp == "" {
		llm.Llamacpp = "./llama.cpp"
	}
	if llm.CudaDevices == nil {
		llm.CudaDevices = []int{0}
	}
	if llm.CtxSize == 0 {
		llm.CtxSize = 2048
	}
	if llm.Temp == 0 {
		llm.Temp = 0.2
	}
	if llm.TopK == 0 {
		llm.TopK = 10000
	}
	if llm.RepeatPenalty == 0 {
		llm.RepeatPenalty = 1.1
	}
	if llm.MaxTokens == 0 {
		llm.MaxTokens = 1000
	}
}

func createPipes(llm *LLM) (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	mainPath := llm.Llamacpp + "/main"
	cmd := exec.Command(mainPath, "-m", llm.Model, "--color", "--ctx_size", fmt.Sprint(llm.CtxSize), "-n", "-1", "-ins", "-b", "128", "--top_k", fmt.Sprint(llm.TopK), "--temp", fmt.Sprint(llm.Temp), "--repeat_penalty", fmt.Sprint(llm.RepeatPenalty), "--n-gpu-layers", fmt.Sprint(llm.Ngl), "-t", "8")
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
	// Close the stdin pipe to signal the end of input
	myerr := stdin.Close()

	if myerr != nil {
		fmt.Println("Error when closing the command:", myerr)
	}

	// Close the communication with the llm
	cmd.Process.Kill()
}

// PromptModel method orderly prompts the LLM with the provided prompts in the array, engaging in a sort of conversation.
// It returns an array with the respones of the LLM, each response matching with the index of its prompt.
func (llm *LLM) PromptModel(prompts []string) ([]string, error) {
	llm.llmDefaults()
	cmd, stdin, stdout, err := createPipes(llm)

	if err != nil {
		fmt.Println("Error creating pipes:", err)
		return []string{"Error creating pipes."}, err
	}

	// Start the llama.cpp llm communication process
	comErr := cmd.Start()
	if comErr != nil {
		fmt.Println("Error starting command:", comErr)
		return []string{"Error starting command."}, comErr
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
		input = llm.InstructionBlock + input

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
		outputs = append(outputs, strings.ReplaceAll(strings.ReplaceAll(output, "\n", ""), ">", ""))
	}

	// Close the communication with the LLM
	closePipes(cmd, stdin, stdout)

	// Return the LLM responses
	return outputs, nil
}
