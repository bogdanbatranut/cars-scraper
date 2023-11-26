package url

type URLComposerImplementations struct {
	implementations map[string]IURLComposer
}

func NewURLComposerImplementations() *URLComposerImplementations {
	impl := make(map[string]IURLComposer)
	impl["autovit"] = NewAutovitURLComposer()
	impl["mobile.de"] = NewMobileDeURLComposer()
	composers := URLComposerImplementations{
		implementations: impl,
	}

	return &composers
}

func (cimpl URLComposerImplementations) GetComposerImplementation(implementationName string) IURLComposer {
	return cimpl.implementations[implementationName]
}
