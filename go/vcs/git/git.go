package git

import (
	vcs ".."
	"log"
)

type Git struct {
}

func init() {
	vcs.Register("git", &Git{})
}

func (self *Git) Install(bundle *vcs.Bundle) error {
	log.Printf("Cloning %s to %s\n", bundle.Url, bundle.Path)
	return vcs.Run("git", "clone", "--recursive", bundle.Url, bundle.Path)
}

func (self *Git) Update(bundle *vcs.Bundle) error {
	log.Printf("Pulling in %s", bundle.Path)
	return vcs.Run("git", "pull", "--ff", "--ff-only")
}
