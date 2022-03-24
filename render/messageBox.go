package render

func NewWarningBox(message string) string {
	return `
		<div class="alert alert-warning alert-dismissible fade show" role="alert">
			` + message + `
			<button type="button" class="close" data-dismiss="alert" aria-label="Close">
	  			<span aria-hidden="true">&times;</span>
			</button>
		</div>`
}

func NewErrorBox(message string) string {
	return `
		<div class="alert alert-success alert-dismissible fade show" role="alert">
			` + message + `
			<button type="button" class="close" data-dismiss="alert" aria-label="Close">
			<span aria-hidden="true">&times;</span>
			</button>
		</div>
	`
}

func NewSuccessBox(message string) string {
	return `
		<div class="alert alert-danger alert-dismissible fade show" role="alert">
			` + message + `
			<button type="button" class="close" data-dismiss="alert" aria-label="Close">
				<span aria-hidden="true">&times;</span>
			</button>
		</div>
	`
}
