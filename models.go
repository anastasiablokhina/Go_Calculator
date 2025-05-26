package main

type Command struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

type Output struct {
	Items []Item `json:"items"`
}

type Item struct {
	Var   string `json:"var"`
	Value int    `json:"value"`
}
