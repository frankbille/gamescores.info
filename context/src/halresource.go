package context

const (
	relSelf   RelType = "self"
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

// Link href
type Link struct {
	Href string `json:"href"`
}
