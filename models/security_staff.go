package models

// SecurityStaff represents airport security personnel
type SecurityStaff struct {
	ID          string
	Name        string
	Position    string
	IsActive    bool
	AssignedLot string
}

func NewSecurityStaff(id, name, position string) *SecurityStaff {
	return &SecurityStaff{
		ID:          id,
		Name:        name,
		Position:    position,
		IsActive:    true,
		AssignedLot: "",
	}
}

func (ss *SecurityStaff) AssignToLot(lotID string) {
	ss.AssignedLot = lotID
}

func (ss *SecurityStaff) UnassignFromLot() {
	ss.AssignedLot = ""
}

func (ss *SecurityStaff) SetActive(status bool) {
	ss.IsActive = status
}

func (ss *SecurityStaff) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"ID":          ss.ID,
		"Name":        ss.Name,
		"Position":    ss.Position,
		"IsActive":    ss.IsActive,
		"AssignedLot": ss.AssignedLot,
	}
}
