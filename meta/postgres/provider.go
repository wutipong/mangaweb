package postgres

import (
	"mangaweb/meta"

	"github.com/jmoiron/sqlx"
)

type Provider struct {
	dbx *sqlx.DB
}

func New() (p Provider, err error) {
	p.dbx, err = connectDB()

	return
}

func (p *Provider) IsItemExist(name string) bool {
	return isItemExist(p.dbx, name)
}

func (p *Provider) Write(i meta.Item) error {
	return write(i, p.dbx)
}

func (p *Provider) New(name string) (i meta.Item, err error) {
	i, err = newItem(p.dbx, name)

	return
}
func (p *Provider) Delete(i meta.Item) error {
	return deleteItem(p.dbx, i)
}
func (p *Provider) Read(name string) (i meta.Item, err error) {
	i, err = readItem(p.dbx, name)

	return
}
func (p *Provider) Open(name string) (i meta.Item, err error) {
	i, err = openItem(p.dbx, name)

	return
}

func (p *Provider) ReadAll() (items []meta.Item, err error) {
	items, err = readAllItems(p.dbx)

	return
}

func (p *Provider) Close() error {
	return p.dbx.Close()
}
