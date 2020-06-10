package lib

import "github.com/greymd/ojichat/generator"

func Ojichat(name string) (string, error) {
	config := generator.Config{
		TargetName:       name,
		EmojiNum:         9,
		PunctuationLevel: 3,
	}
	return generator.Start(config)
}
