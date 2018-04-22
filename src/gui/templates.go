package gui

import "html/template"

/*
  These templates are hardcoded into the application to make deployment
  a bit easier
*/

// GetPoolTemplate returns the template for rendering a pool block
func (gui *GUI) GetPoolTemplate(withChangeOption bool) (*template.Template, error) {
	var changeOption string
	if withChangeOption == true {
		changeOption = `<a href="#" id="change_pool" class="info-block dull">Change</a>`
	}
	temp, err := template.New("pool").Parse(`
    <div class="pool" data-id="{{ .ID }}">
      <h3>{{ .Name }} ` + changeOption + `</h3>
      <a href="{{ .URL }}" target="_blank" class="address">{{ .URL }}</a>
      <div class="stats">
        <table>
          <tr>
            <th>
              Hash Rate
            </th>
            <th>
              Miners
            </th>
            <th>
              Last Block Found
            </th>
          </tr>
          <tr>
            <td>
              {{ .Hashrate }}
            </td>
            <td>
              {{ .Miners }}
            </td>
            <td>
              {{ .LastBlock }}
            </td>
          </tr>
        </table>
      </div>
    </div>
`)
	if err != nil {
		return nil, err
	}
	return temp, nil
}
