package user

type UserRole int

const (
	Unknown UserRole = iota
	Master
	Normal
	Viewer
	Deputy
)

type User struct {
	Name string
	// NetAddress netip.AddrPort
	Role UserRole
}

func CreateUser(username string, role UserRole) *User {
	return &User{
		Name: username,
		// NetAddress: netip.AddrPort{},
		Role: role,
	}
}
