package usecases

type ResponseWithKeys struct {
	Msg             string
	Keys            [][]struct{ Text, Data string }
	IsMessageDelete bool
}

func (r *ResponseWithKeys) AddRow(row ...struct{ Text, Data string }) {
	r.Keys = append(r.Keys, row)
}

func (r *ResponseWithKeys) AddButton(text, data string) struct{ Text, Data string } {
	return struct{ Text, Data string }{Text: text, Data: data}
}
