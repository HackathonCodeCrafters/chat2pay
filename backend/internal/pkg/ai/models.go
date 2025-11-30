package ai

import (
	"github.com/go-skynet/go-llama.cpp"
)

type AIModel interface {
	Prompt(prompt string) (string, error)
}

type aiModel struct {
	model *llama.LLama
}

func NewAiModel(dir string) AIModel {
	model, err := llama.New(
		dir,
		llama.SetContext(1024),
		//llama.SetCPU(true),
	)
	if err != nil {
		panic(err)
	}

	return &aiModel{
		model: model,
	}
}

func (m *aiModel) Prompt(prompt string) (string, error) {
	return m.model.Predict(prompt, llama.SetTemperature(0.1))
}
