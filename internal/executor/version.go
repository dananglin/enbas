package executor

func (v *VersionExecutor) Execute() error {
	v.printer.PrintVersion(v.full)

	return nil
}
