package gui

import "html/template"

/*
  These templates are hardcoded into the application to make deployment
  a bit easier
*/

// GetPoolTemplate returns the template for rendering a pool block
func (gui *GUI) GetPoolTemplate() (*template.Template, error) {
	var poolType string
	if gui.config.CoinType == "bloc" {
		poolType = `<span>({{ .APIType }})</span>`
	}
	temp, err := template.New("pool").Parse(`
	<div class="table-body" data-id="{{ .ID }}">
		<div class="table-col text-left">{{ .URL }} ` + poolType + `</span></div>
		<div class="table-col opacity">{{ .Hashrate }}</div>
		<div class="table-col opacity">{{ .Miners }}</div>
		<div class="table-col opacity">{{ .Fee }}</div>
		<div class="table-col opacity">{{ .Payout }}</div>
		<div class="table-col text-right opacity">{{ .LastBlock }}</div>
		<div class="clearfix"></div>
	</div>
`)
	if err != nil {
		return nil, err
	}
	return temp, nil
}
