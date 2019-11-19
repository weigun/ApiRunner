// pipeline
package models

type Pipeline struct {
	// Require []string `json:"require"  yaml:"require"`
	// Module Module     `json:"module"  yaml:"module"`
	Name     string     `json:"name" yaml:"name"`
	Host     string     `json:"host"  yaml:"host"`
	Def      Variables  `json:"def"  yaml:"def"`
	Steps    []ExecNode `json:"steps"  yaml:"steps"`
	Parallel bool
}

func (pl *Pipeline) GetName() string {
	return pl.Name
}

func (pl *Pipeline) Json() string {
	jsonStr, err := json.Marshal(pl)
	if err != nil {
		log.Warning(`testcase to json failed:`, err.Error())
		return `{}`
	}
	return string(jsonStr)
}

func (pl *Pipeline) GetType() int {
	return TYPE_PIPELINE
}
