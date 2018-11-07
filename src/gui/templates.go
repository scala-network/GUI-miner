package gui

import "html/template"

/*
  These templates are hardcoded into the application to make deployment
  a bit easier
*/

// GetPoolTemplate returns the template for rendering a pool block
func (gui *GUI) GetPoolTemplate(withChangeOption bool) (*template.Template, error) {
	// var changeOption string
	// if withChangeOption == true {
		// changeOption = `<a href="#" id="change_pool" class="info-block dull">Change</a>`
	// }
	temp, err := template.New("pool").Parse(`
	<div class="table-body" data-id="{{ .ID }}">
		<div class="table-col text-left">{{ .URL }} <span>({{ .APIType }})</span></div>
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
