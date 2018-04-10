package model

import (
	"fmt"
	"github.com/ulyssessouza/sterna/migration"
)

type Migration struct {
	Green migration.Config
	Blue  migration.Config
}

func (m *Migration) Check() {
	// TODO Check if BLUE exists
}

func (m *Migration) execute(line *string) uint8 {
	return 1
}

func (m *Migration) Execute() {
	var status uint8 = 0
	for {
		var line string
		fmt.Printf("migrator (%d)$\n", status)
		fmt.Scanln(&line)
		status += m.execute(&line)
	}
}

func (m *Migration) CopyDb() {
	// TODO Run an advertisement procedure that tells the client systems to pass in read only mode. Ex: Setting a flag on consul/etcd
	// TODO Flush data on BLUE
	// TODO Lock BLUE tuning it into read only
	// TODO Dump the content of all databases to a persistent FS. Ex: S3 or Hadoop
	// TODO Create service/pods for GREEN based on the configuration copied from BLUE on k8s. (I think it's easier than determining if the configuration of an existing one is compatible)
	// TODO Fill GREEN with the data dumped from BLUE
}

func (m *Migration) MigrateBlue() {
	// TODO Check if configured update scripts apply to m database version. Abort if not
	// TODO Run the migration scripts against GREEN
	// TODO Update DB version
}

func (m *Migration) UpdateLiveConfiguration() {
	// TODO Update consul/etcd with GREEN coordinates as the new current database
}

func (m *Migration) SpawnBlueApps() {
	// At m point the new application should spawn and reach the new configuration
}

func (m *Migration) LoadIntoBlue() {
	// TODO Configure load balancer to target BLUE on new sessions
	// In case of problems, you can still come back to GREEN by rolling back the LOAD step
}

func (m *Migration) ScaleBalance() {
	// TODO Scale down BLUE
	// TODO Scale up GREEN
	// TODO Repeat until 100% GREEN
}
