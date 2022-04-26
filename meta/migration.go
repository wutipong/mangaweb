package meta

type MigrateFunction func(m Meta) (out Meta, err error)

var migrateFunctions = make(map[int]MigrateFunction)

func init() {
	migrateFunctions[0] = migrateV0
}

func migrateV0(m Meta) (out Meta, err error) {
	m.PopulateTags()
	out = m
	return
}

func Migrate(m Meta) (out Meta, err error) {
	return doMigrate(m, migrateFunctions, CurrentVersion)
}

func doMigrate(m Meta, functions map[int]MigrateFunction, targetVersion int) (out Meta, err error) {
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
