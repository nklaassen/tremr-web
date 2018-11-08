package datastore

type Tremor struct {
	Tid       int `json:"tid"`
	Resting   int `json:"resting"`
	Postural  int `json:"postural"`
	Completed bool `json:"completed"`
	Date      string `json:"date"`
}

type TremorRepo interface {
	Add(*Tremor) error
	GetAll() ([]Tremor, error)
}
