package notice

type NoticerSet struct {
	BarkSet    BarkSet  `yaml:"bark"`
	SoundSet   SoundSet `yaml:"sound"`
	NoticeType int      `yaml:"noticeType"`
}

func Do(noticerSet NoticerSet) error {

	switch noticerSet.NoticeType {
	case 0:
		return nil
	case 1:
		return BarkPush(noticerSet.BarkSet)
	case 2:
		return MacSound(noticerSet.SoundSet)
	default:
		return nil
	}
}
