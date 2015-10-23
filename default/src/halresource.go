package defaultapp

const (
	relSelf   RelType = "self"
	relCreate RelType = "create"
	relUpdate RelType = "update"
	relDelete RelType = "delete"
	relNext   RelType = "next"
	relPrev   RelType = "prev"
	relFirst  RelType = "first"
	relLast   RelType = "last"
)

// RelType is the type of the relation
type RelType string

// HalResource is the interface that should be added to all types that should
// have links.
type HalResource interface {
	// AddLink adds a new link to the resource
	AddLink(rel RelType, href string)

	// RemoveLink removes a new link from the resource
	RemoveLink(rel RelType)
}

// DefaultHalResource is the default implementation of the HalResource interface
type DefaultHalResource struct {
	Links map[RelType]Link `datastore:"-" json:"_links,omitempty"`
}

// AddLink adds a new link to the resource
func (lr *DefaultHalResource) AddLink(rel RelType, href string) {
	if lr.Links == nil {
		lr.Links = make(map[RelType]Link)
	}

	lr.Links[rel] = Link{
		Href: href,
	}
}

// RemoveLink removes a new link from the resource
func (lr *DefaultHalResource) RemoveLink(rel RelType) {
	if lr.Links != nil {
		delete(lr.Links, rel)
	}
}

// Link href
type Link struct {
	Href string `json:"href"`
}
