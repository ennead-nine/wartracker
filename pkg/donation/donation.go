package donation

type Donation struct {
	Id          string `json:"id" yaml:"id" db:"id"`
	Amount      int    `json:"amount" yaml:"amount" db:"amount"`
	Date        string `json:"date" yaml:"date" db:"date"`
	CommanderId string `json:"commander-id" yaml:"commanderId" db:"commander_id"`
}

type DonationMap map[string]Donation
