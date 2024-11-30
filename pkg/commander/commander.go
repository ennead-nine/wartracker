package commander

type Commander struct {
	Id       string `json:"id" yaml:"id"`
	NoteName string `json:"note-name" yaml:"noteName"`
}

type CommanderData struct {
	Date           string `json:"date" yaml:"date"`
	Name           string `json:"name" yaml:"name"`
	Kills          int64  `json:"kills" yaml:"kills"`
	HQPower        int64  `json:"hq-power" yaml:"HqPower"`
	TotalHeroPower int64  `json:"total-hero-power" yaml:"totalHeroPower"`
	Commander
}
