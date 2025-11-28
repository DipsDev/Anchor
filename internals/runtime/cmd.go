package runtime

func ApplyEnvironmentCmd(applyConfig EnvironmentStatusOptions) error {
	return modifyEnvironment(applyEnvironment, applyConfig)
}

func StopEnvironmentCmd(stopConfig EnvironmentStatusOptions) error {
	return modifyEnvironment(stopEnvironment, stopConfig)
}
