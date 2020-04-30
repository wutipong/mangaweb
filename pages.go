package main

type Page struct {
	Index int
	Name  string
}

func ListPages(file string) (pages []Page, err error) {
	var meta itemMeta
	err = meta.Read(file)
	if err != nil {
		return
	}

	if len(meta.Pages) == 0 {
		meta.GeneratePages()
		meta.Write()
	}

	fileNames := meta.Pages

	pages = make([]Page, len(fileNames))
	for i, f := range fileNames {
		pages[i] = Page{
			Name:  f,
			Index: i,
		}
	}

	return
}
