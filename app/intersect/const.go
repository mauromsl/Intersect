package intersect

const issuePrefixRE string = "^#[0-9]+: "

//we can't escape the backtick character in a string literal...
const TRELLO_COMMENT_TMPL = `
` + "`COMMENT_ID: %d`" + `
` + "`USER: %s`" + `
` + "`CREATED: %s`" + `

%s
`
