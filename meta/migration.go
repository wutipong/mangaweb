package meta

type MigrateFunction func(m Item) (out Item, err error)

var migrateFunctions = make(map[int]MigrateFunction)

func init() {
	migrateFunctions[0] = migrateV0
}

func migrateV0(m Item) (out Item, err error) {
	out = m
	return
}

func Migrate(m Item) (out Item, err error) {
	return doMigrate(m, migrateFunctions, CurrentItemVersion)
}

func doMigrate(m Item, functions map[int]MigrateFunction, targetVersion int) (out Item, err error) {
	temp := m
	for v := m.Version; v < targetVersion; v++ {
		f := functions[v]
		temp, err = f(temp)
		if err != nil {
			return
		}
		temp.Version = v
	}

	temp.Version = targetVersion
	out = temp
	return
}
