package state


var (
	State GlobalState
)

type GlobalState interface {
	CheckLogin(username, password string) error
}

func init() {
	// TODO
	State = &FileState{
		users:map[string]User{
			"hello":User{Username:"hello", Password:"world"},
		},
	}
}
