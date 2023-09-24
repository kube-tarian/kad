package agent

type appStatus string

const (
	appIntallingStatus      appStatus = "Installing"
	appIntalledStatus       appStatus = "Installed"
	appUpgradingStatus      appStatus = "Upgrading"
	appUpgradedStatus       appStatus = "Upgraded"
	appIntallFailedStatus   appStatus = "Installion Failed"
	appUpgradeFaileddStatus appStatus = "Upgrade Failed"
)

type appDeployAction string

const (
	appInstallAction appDeployAction = "install"
	appUpgradeAction appDeployAction = "upgrade"
)
