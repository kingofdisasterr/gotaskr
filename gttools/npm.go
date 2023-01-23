package gttools

import (
	"os/exec"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/goext"
)

// NpmTool provides access to the helper methods for npm.
type NpmTool struct {
}

func CreateNpmTool() *NpmTool {
	return &NpmTool{}
}

// NpmCommonSettings are common settings used for all commands.
type NpmCommonSettings struct {
	WorkingDirectory string
	OutputToConsole  bool
	CustomArguments  []string
}

// Customize adds a custom argument to the settings object.
func (s *NpmCommonSettings) Customize(setting string) *NpmCommonSettings {
	s.CustomArguments = append(s.CustomArguments, setting)
	return s
}

// NpmInitSettings are the settings used for Init.
type NpmInitSettings struct {
	NpmCommonSettings
	PackageSpec string
}

func (tool *NpmTool) Init(outputToConsole bool, settings *NpmInitSettings) error {
	if settings == nil {
		settings = &NpmInitSettings{}
	}
	args := []string{
		"init",
		settings.PackageSpec,
		"-y",
	}
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

// NpmRunSettings are the settings used for Run.
type NpmRunSettings struct {
	NpmCommonSettings
	Script string
}

func (tool *NpmTool) Run(settings *NpmRunSettings) error {
	if settings == nil {
		settings = &NpmRunSettings{}
	}
	args := []string{
		"run",
		settings.Script,
	}
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", args...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

func (tool *NpmTool) RunScript(outputToConsole bool, script string) error {
	return tool.Run(&NpmRunSettings{
		NpmCommonSettings: NpmCommonSettings{
			OutputToConsole: outputToConsole,
		},
		Script: script,
	})
}

// NpmCleanInstallSettings are the settings used for CleanInstall.
type NpmCleanInstallSettings struct {
	NpmCommonSettings
	CacheDir      string
	NoAudit       bool
	PreferOffline bool
}

func (tool *NpmTool) CleanInstall(outputToConsole bool, settings *NpmCleanInstallSettings) error {
	if settings == nil {
		settings = &NpmCleanInstallSettings{}
	}
	args := []string{
		"ci",
	}
	args = goext.AddIf(args, settings.NoAudit, "--no-audit")
	args = goext.AddIf(args, settings.PreferOffline, "--prefer-offline")
	args = tool.addCache(args, settings.CacheDir)
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

// NpmInstallSettings are the settings used for Install.
type NpmInstallSettings struct {
	NpmCommonSettings
	CacheDir      string
	PackageSpec   string
	NoAudit       bool
	PreferOffline bool
	SaveProd      bool
	SaveDev       bool
	SaveOptional  bool
	SaveExact     bool
}

func (tool *NpmTool) Install(outputToConsole bool, settings *NpmInstallSettings) error {
	if settings == nil {
		settings = &NpmInstallSettings{}
	}
	args := []string{
		"install",
		settings.PackageSpec,
	}
	args = goext.AddIf(args, settings.SaveProd, "--save-prod")
	args = goext.AddIf(args, settings.SaveDev, "--save-dev")
	args = goext.AddIf(args, settings.SaveOptional, "--save-optional")
	args = goext.AddIf(args, settings.SaveExact, "--save-exact")
	args = goext.AddIf(args, settings.NoAudit, "--no-audit")
	args = goext.AddIf(args, settings.PreferOffline, "--prefer-offline")
	args = tool.addCache(args, settings.CacheDir)
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	return execr.RunCommand(settings.OutputToConsole, cmd)
}

// NpmBinSettings are the settings used for Bin.
type NpmBinSettings struct {
	NpmCommonSettings
	Global bool
}

func (tool *NpmTool) Bin(settings *NpmBinSettings) (string, error) {
	if settings == nil {
		settings = &NpmBinSettings{}
	}
	args := []string{
		"bin",
	}
	args = goext.AddIf(args, settings.Global, "--global")
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory
	stdout, _, err := execr.RunCommandGetOutput(settings.OutputToConsole, cmd)
	return stdout, err
}

// NpmPublishSettings are the settings used for Publish.
type NpmPublishSettings struct {
	NpmCommonSettings
	PackageSpec string
	Tag         string
}

func (tool *NpmTool) Publish(settings *NpmPublishSettings) (string, error) {
	if settings == nil {
		settings = &NpmPublishSettings{}
	}
	args := []string{
		"publish",
		settings.PackageSpec,
	}
	args = goext.AddIf(args, len(settings.Tag) > 0, "--tag", settings.Tag)
	args = append(args, settings.CustomArguments...)

	cmd := exec.Command("npm", goext.RemoveEmpty(args)...)
	cmd.Dir = settings.WorkingDirectory

	stdout, _, err := execr.RunCommandGetOutput(settings.OutputToConsole, cmd)

	return stdout, err
}

////////////////////////////////////////////////////////////
// Internal Methods
////////////////////////////////////////////////////////////

func (tool *NpmTool) addCache(args []string, cacheDir string) []string {
	return goext.AddIf(args, len(cacheDir) > 0, "--cache", cacheDir)
}
