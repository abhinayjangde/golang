package person

type Person struct {
	name  string
	email string
}

func (p Person) GetName() string {
	return p.name
}

func (p Person) GetEmail() string {
	return p.email
}
func (p *Person) SetName(name string) {
	p.name = name
}
func NewPerson(name string, email string) *Person {
	return &Person{
		name,
		email,
	}
}
