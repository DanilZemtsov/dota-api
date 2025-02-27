package model

type Hero struct {
	ID          int
	Name        string
	Attribute   string
	AttributeID int64
	Winrate     int
	Picha       string
	History     string
	Items       []Item
}
