package v1beta1


var (
	defaultPermitWaitingTimeSeconds      int64 = 60
	defaultDeniedPGExpirationTimeSeconds int64 = 20
)

// SetDefaultsCoschedulingArgs sets the default parameters for Coscheduling plugin.
func SetDefaultsCoschedulingArgs(obj *CoschedulingArgs) {
	if obj.PermitWaitingTimeSeconds == nil {
		obj.PermitWaitingTimeSeconds = &defaultPermitWaitingTimeSeconds
	}
	if obj.DeniedPGExpirationTimeSeconds == nil {
		obj.DeniedPGExpirationTimeSeconds = &defaultDeniedPGExpirationTimeSeconds
	}
	// TODO(k/k#96427): get KubeConfigPath and KubeMaster from configuration or command args.
}