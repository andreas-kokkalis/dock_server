package authspec

var AdminLoginGood = `
{
  "data": "SUCCESS"
}
`

var AdminLoginPasswordMismatch = `
{
	"errors": [
        "Password mismatch"
      ],
      "status": "Unauthorized"
    }
`

var AdminLoginUnauthorized = `
{
	"errors": [
  		"Unauthorized"
      ],
      "status": "Unauthorized"
    }
`

var AdminLogoutUnauthorized = `
{
  "errors" : [
  	"Unauthorized"
	],
  "status" : "Unauthorized"
}
`
