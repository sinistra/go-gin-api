package frontierlink

type ActivityNote struct {
	Notes        string `xml:"notes,omitempty" json:"notes,omitempty"`
	CreationDate string `xml:"creationDate,omitempty" json:"creationDate,omitempty"`
}

type ActivityNoteList struct {
	ActivityNote []*ActivityNote `xml:"activityNote,omitempty" json:"activityNote,omitempty"`
}

type ProductOrderActivity struct {
	ActivityType     string            `xml:"activityType,omitempty" json:"activityType,omitempty"`
	ActualStartDate  string            `xml:"actualStartDate,omitempty" json:"actualStartDate,omitempty"`
	ActualEndDate    string            `xml:"actualEndDate,omitempty" json:"actualEndDate,omitempty"`
	ActivityNoteList *ActivityNoteList `xml:"activityNoteList,omitempty" json:"activityNoteList,omitempty"`
}

type ServiceOrderDetail struct {
	ServiceOrderID   string `xml:"serviceOrderID,omitempty" json:"serviceOrderID,omitempty"`
	ServiceOrderType string `xml:"serviceOrderType,omitempty" json:"serviceOrderType,omitempty"`
}

type ServiceOrderActivity struct {
	ActivityType     string            `xml:"activityType,omitempty" json:"activityType,omitempty"`
	ActualStartDate  string            `xml:"actualStartDate,omitempty" json:"actualStartDate,omitempty"`
	ActualEndDate    string            `xml:"actualEndDate,omitempty" json:"actualEndDate,omitempty"`
	ActivityNoteList *ActivityNoteList `xml:"activityNoteList,omitempty" json:"activityNoteList,omitempty"`
}

type ServiceOrder struct {
	ServiceOrderDetail   *ServiceOrderDetail     `xml:"serviceOrderDetail,omitempty" json:"serviceOrderDetail,omitempty"`
	ServiceOrderActivity []*ServiceOrderActivity `xml:"serviceOrderActivity,omitempty" json:"serviceOrderActivity,omitempty"`
}

type ServiceOrderList struct {
	ServiceOrder []*ServiceOrder `xml:"serviceOrder,omitempty" json:"serviceOrder,omitempty"`
}

type ListActivities struct {
	ProductOrderID       string                  `xml:"productOrderID,omitempty" json:"productOrderID,omitempty"`
	ProductOrderActivity []*ProductOrderActivity `xml:"productOrderActivity,omitempty" json:"productOrderActivity,omitempty"`
	ServiceOrderList     *ServiceOrderList       `xml:"serviceOrderList,omitempty" json:"serviceOrderList,omitempty"`
}
