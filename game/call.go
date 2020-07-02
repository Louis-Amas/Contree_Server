package game

// Call struct
type Call struct {
	Value        int    `json:"value,omitempty"`
	Family       string `json:"family,omitempty"`
	IsCapot      bool   `json:"isCapot,omitempty"`
	IsCapotB     bool   `json:"isCapotB,omitempty"`
	IsContree    bool   `json:"isContree,omitempty"`
	IsSurContree bool   `json:"isSurContree,omitempty"`
	IsPass       bool   `json:"isPass,omitempty"`
}

// compareTo return true and error nil if call2 is valid and better than call1
// else return error
func (call1 *Call) compareTo(call2 *Call) bool {

	if call2.IsPass {
		return false
	}
	if call1.Value < call2.Value {
		return true
	}
	if !call1.IsCapot && call2.IsCapot {
		return true
	}

	if call1.IsCapot && call2.IsCapotB {
		return true
	}

	if !call1.IsContree && call2.IsContree {
		return true
	}
	if call1.IsContree && call2.IsSurContree {
		return true
	}

	return false
}
