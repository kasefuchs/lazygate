package kick

import (
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/util/componentutil"
)

type RawTextComponent string

func (d *RawTextComponent) Set(s string) error {
	*d = RawTextComponent(s)

	_, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, s)
	return err
}

func (d *RawTextComponent) String() string {
	return string(*d)
}

func (d *RawTextComponent) TextComponent() *component.Text {
	comp, err := componentutil.ParseTextComponent(version.MinimumVersion.Protocol, d.String())
	if err != nil {
		return &component.Text{}
	}

	return comp
}
