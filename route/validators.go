package route

import "regexp"

var vImageID = regexp.MustCompile(`^([A-Fa-f0-9]{12,64})$`)
var vContainerID = regexp.MustCompile(`^([A-Fa-f0-9]{12,64})$`)
var vContainerState = regexp.MustCompile(`^(|created|restarting|running|paused|exited|dead)$`) // Can also be empty
var vPassword = regexp.MustCompile(`^([a-zA-Z0-9]){5,6}$`)
var vUsername = regexp.MustCompile(`^([a-zA-Z0-9]){2,16}$`)
