package notice

import (
	"fmt"
	"os/exec"
)

type SoundSet struct {
	Message string `yaml:"soundMessage"`
	Times   int    `yaml:"soundTimes"`
	Voice   string `yaml:"soundVoice"`
}

// MacSound only for mac
func MacSound(soundSet SoundSet) error {
	for i := 0; i < soundSet.Times; i++ {
		err := exec.Command("say", soundSet.Message, fmt.Sprintf("--voice=%s", soundSet.Voice)).Run()
		if err != nil {
			return err
		}
	}
	return nil
}
